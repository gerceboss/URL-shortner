package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

//redis as a database

var Ctx=context.Background()

func CreateClient(dbNo int) *redis.Client{
	rdb:=redis.NewClient(&redis.Options{
		Addr:os.Getenv("DB_ADDRESS"),
		DB:dbNo,
		Password:os.Getenv("DB_PASSWORD"),
	})

	return rdb
}