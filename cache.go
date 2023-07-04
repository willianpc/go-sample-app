package main

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

func client() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
}

func readCache(key string) []string {
	rdb := client()

	return rdb.SMembers(context.Background(), key).Val()
}

func writeCache(key string, vals []string) error {
	ctx := context.Background()
	rdb := client()

	tx := rdb.TxPipeline()
	tx.SAdd(ctx, key, vals)
	tx.Expire(ctx, key, time.Second*15)
	_, err := tx.Exec(ctx)

	return err
}
