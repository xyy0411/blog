package config

import (
	"github.com/go-redis/redis"
	"github.com/xyy0411/blog/global"
)

var (
	RedisClient *redis.Client
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		global.Logger.Error("连接redis失败，error:", err)
		return
	}
	global.Logger.Info("连接redis成功")
}
