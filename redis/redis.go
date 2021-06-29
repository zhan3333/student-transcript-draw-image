package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

var RDB *redis.Client

func InitRedis() error {
	index, err := strconv.ParseInt(os.Getenv("REDIS_INDEX"), 10, 64)
	if err != nil {
		return fmt.Errorf("REDIS_INDEX 读取失败: %w", err)
	}
	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS"),
		Password: "",         // no password set
		DB:       int(index), // use default DB
	})
	if err = RDB.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("ping redis err: %w", err)
	}
	return nil
}
