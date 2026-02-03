package dao

import (
	"exchangeapp/global"
	"strconv"

	"github.com/go-redis/redis"
)

type LikeDAO interface {
	LikeArticle(articleID string) error
	GetArticleLikes(articleID string) (int64, error)
	UnlikeArticle(articleID string) error
}

type likeDAO struct {
	redisDb *redis.Client
}

func NewLikeDAO() LikeDAO {
	return &likeDAO{
		redisDb: global.RedisDb,
	}
}

func (ld *likeDAO) LikeArticle(articleID string) error {
	likeKey := "article:" + articleID + ":likes"
	return ld.redisDb.Incr(likeKey).Err()
}

func (ld *likeDAO) GetArticleLikes(articleID string) (int64, error) {
	likeKey := "article:" + articleID + ":likes"
	val, err := ld.redisDb.Get(likeKey).Result()

	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	likes, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return likes, nil
}

func (ld *likeDAO) UnlikeArticle(articleID string) error {
	likeKey := "article:" + articleID + ":likes"
	// 使用Lua脚本防止点赞数为负
	script := redis.NewScript(`
		local key = KEYS[1]
		local current = tonumber(redis.call('GET', key) or 0)
		if current > 0 then
			return redis.call('DECR', key)
		else
			return 0
		end
	`)

	_, err := script.Run(ld.redisDb, []string{likeKey}).Result()
	return err
}
