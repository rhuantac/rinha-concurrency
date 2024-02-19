package config

import (
	"os"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

func SetupRedis() *redsync.Redsync{
	client := goredislib.NewClient(&goredislib.Options{
		Addr: os.Getenv("REDIS_HOST"),
	})
	pool := goredis.NewPool(client)
	return redsync.New(pool)
}
