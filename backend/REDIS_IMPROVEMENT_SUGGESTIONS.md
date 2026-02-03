# Redis 改进建议方案

## 当前问题分析

### 1. **Redis 使用场景单一**
- ✗ 仅用于简单的缓存（articles、likes）
- ✗ 没有充分利用 Redis 的多样化数据结构
- ✗ 缺少业务关键场景的优化

### 2. **Redis 配置缺陷**
- ✗ 连接池配置不完整（PoolSize, MinIdleConns 未配置）
- ✗ 无连接重试机制
- ✗ 无心跳检测（Ping）
- ✗ 错误处理不够完善
- ✗ 连接字符串格式错误（`"localhost: 6379"` 有多余空格）

### 3. **缓存策略问题**
- ✗ 文章缓存失效机制简单（仅 Delete）
- ✗ 点赞计数没有持久化到数据库
- ✗ 缺少缓存预热机制
- ✗ 缺少缓存穿透/击穿/雪崩的防护

### 4. **业务功能缺失**
- ✗ 登录会话管理不完整（没有 Session/Token 缓存）
- ✗ 缺少用户黑名单功能
- ✗ 汇率数据未充分利用 Redis 加速
- ✗ 缺少限流/计数器功能
- ✗ 缺少分布式锁防止并发问题

---

## 改进方案（优先级排序）

### **优先级 1：基础配置完善**

#### 1.1 改进 Redis 配置（config/redis.go）
```go
package config

import (
	"exchangeapp/global"
	"log"
	"time"

	"github.com/go-redis/redis"
)

func InitRedis() {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",  // 修复格式
		DB:           0,
		Password:     "",
		PoolSize:     10,               // 连接池大小
		MinIdleConns: 5,                // 最小空闲连接数
		MaxRetries:   3,                // 重试次数
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// 带重试的连接测试
	var err error
	for i := 0; i < 3; i++ {
		_, err = RedisClient.Ping().Result()
		if err == nil {
			log.Println("Redis connected successfully")
			global.RedisDb = RedisClient
			return
		}
		log.Printf("Redis connection attempt %d failed: %v\n", i+1, err)
		time.Sleep(time.Second * time.Duration(i+1))
	}
	log.Fatalf("Failed to connect to Redis after retries: %v", err)
}
```

### **优先级 2：缓存改进**

#### 2.1 实现缓存工具库（utils/cache.go - 新建）
```go
package utils

import (
	"encoding/json"
	"exchangeapp/global"
	"time"

	"github.com/go-redis/redis"
)

// 缓存穿透防护：使用空值缓存
func GetCache(key string) (string, error) {
	val, err := global.RedisDb.Get(key).Result()
	if err == redis.Nil {
		return "", nil // 返回空字符串而非错误
	}
	return val, err
}

// 通用缓存设置（支持 JSON 序列化）
func SetCache(key string, data interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return global.RedisDb.Set(key, string(jsonData), expiration).Err()
}

// 获取并反序列化缓存
func GetCacheJSON(key string, dest interface{}) error {
	val, err := global.RedisDb.Get(key).Result()
	if err == redis.Nil {
		return nil // 缓存不存在
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// 缓存穿透防护：设置空值缓存
func SetEmptyCache(key string) error {
	return global.RedisDb.Set(key, "", 30*time.Second).Err()
}

// 检查是否为空值缓存
func IsEmptyCache(val string) bool {
	return val == ""
}

// 批量删除缓存（用于失效）
func DeleteCache(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return global.RedisDb.Del(keys...).Err()
}
```

#### 2.2 改进文章控制器（controllers/article_controller.go）
- 添加缓存穿透防护（空值缓存）
- 改进缓存失效策略
- 使用统一的缓存工具函数

#### 2.3 改进汇率缓存策略（controllers/exchange_rate_controller.go）
```go
// 新增函数
const (
	exchangeRateCacheKey = "exchange_rates"
	exchangeRateCacheTTL = 1 * time.Hour // 汇率数据1小时更新一次
)

func GetExchangeRates(ctx *gin.Context) {
	// 先查缓存
	var exchangeRates []models.ExchangeRate
	if err := GetCacheJSON(exchangeRateCacheKey, &exchangeRates); err == nil && len(exchangeRates) > 0 {
		ctx.JSON(http.StatusOK, exchangeRates)
		return
	}

	// 缓存不存在，查数据库
	if err := global.Db.Find(&exchangeRates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置缓存
	_ = SetCache(exchangeRateCacheKey, exchangeRates, exchangeRateCacheTTL)
	ctx.JSON(http.StatusOK, exchangeRates)
}
```

---

### **优先级 3：会话和认证管理**

#### 3.1 创建会话管理工具（utils/session.go - 新建）
```go
package utils

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"fmt"
	"time"
)

const (
	sessionPrefix = "session:"
	sessionTTL    = 24 * time.Hour // 会话有效期24小时
)

// 存储用户会话
func CreateSession(userID uint, username string) (string, error) {
	sessionKey := fmt.Sprintf("%s%d", sessionPrefix, userID)
	sessionData := map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"login_at": time.Now().Unix(),
	}
	return sessionKey, SetCache(sessionKey, sessionData, sessionTTL)
}

// 验证会话
func GetSession(userID uint) (map[string]interface{}, error) {
	sessionKey := fmt.Sprintf("%s%d", sessionPrefix, userID)
	val, err := global.RedisDb.Get(sessionKey).Result()
	if err != nil {
		return nil, err
	}
	// 解析 JSON 并返回
	var session map[string]interface{}
	err = json.Unmarshal([]byte(val), &session)
	return session, err
}

// 销毁会话（登出）
func DestroySession(userID uint) error {
	sessionKey := fmt.Sprintf("%s%d", sessionPrefix, userID)
	return global.RedisDb.Del(sessionKey).Err()
}

// 刷新会话过期时间
func RefreshSession(userID uint) error {
	sessionKey := fmt.Sprintf("%s%d", sessionPrefix, userID)
	return global.RedisDb.Expire(sessionKey, sessionTTL).Err()
}
```

#### 3.2 创建黑名单管理（utils/blacklist.go - 新建）
```go
package utils

import (
	"exchangeapp/global"
	"fmt"
	"time"
)

const blacklistPrefix = "token_blacklist:"

// 将 Token 加入黑名单（登出时）
func AddToBlacklist(token string, expiration time.Duration) error {
	key := fmt.Sprintf("%s%s", blacklistPrefix, token)
	return global.RedisDb.Set(key, "1", expiration).Err()
}

// 检查 Token 是否在黑名单
func IsTokenBlacklisted(token string) bool {
	key := fmt.Sprintf("%s%s", blacklistPrefix, token)
	val, err := global.RedisDb.Get(key).Result()
	return val != "" && err == nil
}
```

---

### **优先级 4：限流和防护**

#### 4.1 创建限流器（utils/ratelimit.go - 新建）
```go
package utils

import (
	"exchangeapp/global"
	"fmt"
	"strconv"
	"time"
)

// 滑动窗口限流
func CheckRateLimit(userID string, action string, limit int, windowSize time.Duration) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s:%s", userID, action)
	
	current := global.RedisDb.Incr(key).Val()
	if current == 1 {
		// 第一次访问，设置过期时间
		global.RedisDb.Expire(key, windowSize)
	}
	
	return current <= int64(limit), nil
}

// 获取剩余请求次数
func GetRemainingQuota(userID string, action string, limit int) (int, error) {
	key := fmt.Sprintf("rate_limit:%s:%s", userID, action)
	current, err := global.RedisDb.Get(key).Int64()
	if err != nil {
		return limit, nil
	}
	remaining := limit - int(current)
	if remaining < 0 {
		return 0, nil
	}
	return remaining, nil
}
```

#### 4.2 创建分布式锁（utils/lock.go - 新建）
```go
package utils

import (
	"exchangeapp/global"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// 分布式锁
func AcquireLock(resource string, ttl time.Duration) (string, error) {
	key := fmt.Sprintf("lock:%s", resource)
	lockID := uuid.New().String()
	
	// SET NX EX 原子操作：只在 key 不存在时设置
	ok, err := global.RedisDb.SetNX(key, lockID, ttl).Result()
	if !ok {
		return "", fmt.Errorf("lock already held")
	}
	return lockID, err
}

// 释放锁（Lua 脚本防止误释放）
func ReleaseLock(resource string, lockID string) error {
	key := fmt.Sprintf("lock:%s", resource)
	
	luaScript := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	
	result, err := global.RedisDb.Eval(luaScript, []string{key}, lockID).Result()
	if result == int64(1) {
		return nil
	}
	return fmt.Errorf("failed to release lock: %v", err)
}
```

---

### **优先级 5：业务优化**

#### 5.1 点赞计数持久化
```go
// 每收集到一定数量点赞，或定时将 Redis 数据回源到数据库
func SyncLikesToDB(articleID uint) error {
	key := fmt.Sprintf("article:%d:likes", articleID)
	likes, err := global.RedisDb.Get(key).Int64()
	if err != nil {
		return err
	}
	
	// 更新数据库
	return global.Db.Model(&models.Article{}).
		Where("id = ?", articleID).
		Update("likes", likes).Error
}
```

#### 5.2 热点数据预热
```go
// 应用启动时预热缓存
func WarmupCache() {
	// 预热所有文章
	var articles []models.Article
	global.Db.Find(&articles)
	SetCache(cacheKey, articles, 10*time.Minute)
	
	// 预热所有汇率
	var rates []models.ExchangeRate
	global.Db.Find(&rates)
	SetCache(exchangeRateCacheKey, rates, 1*time.Hour)
}
```

---

## Redis 经典使用场景对标

| 场景 | 当前使用 | 建议方案 |
|------|--------|--------|
| **缓存** | ✓ 简单 | ✓ 优化穿透/击穿/雪崩防护 |
| **会话管理** | ✗ | **✓ 实现 Session 存储** |
| **黑名单** | ✗ | **✓ Token 黑名单管理** |
| **计数器** | ✓ 点赞 | **✓ 限流、访问计数** |
| **发布订阅** | ✗ | 可选（实时通知） |
| **消息队列** | ✗ | 可选（异步任务） |
| **分布式锁** | ✗ | **✓ 防止并发问题** |
| **HyperLogLog** | ✗ | 可选（UV 统计） |
| **热点数据预热** | ✗ | **✓ 应用启动预热** |
| **缓存预热/更新** | 简单 | **✓ 定时任务更新** |

---

## 实施步骤

1. **第一阶段（基础）**
   - [ ] 修复 Redis 配置
   - [ ] 创建缓存工具库
   - [ ] 改进现有缓存逻辑

2. **第二阶段（认证/会话）**
   - [ ] 实现会话管理
   - [ ] 实现 Token 黑名单
   - [ ] 改进认证中间件

3. **第三阶段（防护）**
   - [ ] 实现限流功能
   - [ ] 实现分布式锁
   - [ ] 添加到关键业务逻辑

4. **第四阶段（优化）**
   - [ ] 点赞计数持久化
   - [ ] 热点数据预热
   - [ ] 定时缓存更新任务

---

## 预期面试亮点

✓ Redis 配置最佳实践（连接池、超时、重试）
✓ 缓存穿透/击穿/雪崩三大问题解决方案
✓ 会话管理和身份认证安全性
✓ 分布式锁防止并发问题
✓ 限流保护系统
✓ Lua 脚本原子性操作
✓ 完整的错误处理和日志记录
✓ 代码结构清晰，易于维护和扩展
