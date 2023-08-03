package main

import (
	"context"
	"time"
)

func readCache(ctx context.Context, key string) []string {
	return rdb.SMembers(ctx, key).Val()
}

func writeCache(ctx context.Context, key string, vals []string) error {
	tx := rdb.TxPipeline()
	tx.SAdd(ctx, key, vals)
	tx.Expire(ctx, key, time.Second*15)
	_, err := tx.Exec(ctx)

	return err
}
