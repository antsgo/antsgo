package redis

import (
	"fmt"

	"github.com/antsgo/antsgo/conf"
	"github.com/go-redis/redis"
)

func NewRedis(c conf.ConfigRedis) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.Db,
	})

	if _, err = client.Ping().Result(); err != nil {
		fmt.Println("redis连接失败:", err)
		return
	}
	fmt.Println("Redis已连接...")
	return
}
