package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func ConnectRedis() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	db := 0

	addr := fmt.Sprintf("%s:%s", host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("ไม่สามารถเชื่อม Redis: %v", err)
	}
	log.Println("เชื่อมต่อกับ Redis สำเร็จแล้ว")
	return rdb
}

func SetToken(rdb *redis.Client, key string, value string, duration time.Duration) error {
	return rdb.Set(Ctx, key, value, duration).Err()
}

func GetToken(rdb *redis.Client, key string) (string, error) {
	return rdb.Get(Ctx, key).Result()
}

func IsTokenValid(rdb *redis.Client, key string, token string) (bool, error) {
	ctx := context.Background()
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if val != token {
		return false, nil
	}
	return true, nil
}
