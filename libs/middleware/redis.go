package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const CtxRedisConn = "redis-conn"

func Redis() gin.HandlerFunc {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		c.Set(CtxRedisConn, client)
	}
}
