package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB //define global database
	RedisDb *redis.Client
)