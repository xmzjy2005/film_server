package db

import (
	"context"
	"film_server/config"
	"github.com/redis/go-redis/v9"
	"time"
)

var Rdb *redis.Client
var Cxt = context.Background()

func InitRedisConn() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:        config.RedisAddr,
		Password:    config.RedisPassword,
		DB:          config.RedisDBNo,
		PoolSize:    10,               // 最大连接数
		DialTimeout: time.Second * 10, // 超时时间
	})
	//这个没有再文档里找到？
	_, err := Rdb.Ping(Cxt).Result()
	if err != nil {
		panic(err)
	}
	return nil

}
