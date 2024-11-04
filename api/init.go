package handler

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func init() {
	os.Setenv("REDIS_URL", "redis://default:O1Z0NHQMsRSJpanSKiQo6Qb7d6KLrUno@redis-16009.c100.us-east-1-4.ec2.redns.redis-cloud.com:16009")

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		panic("REDIS_URL environment variable not set")
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}
	opt.PoolSize = 10
	opt.MinIdleConns = 3

	// go func() {
	// 	for {
	// 		if row, err := callCaibai(); err == nil {
	// 			insertToRedis(row)
	// 		}
	// 		time.Sleep(10 * time.Minute)
	// 	}
	// }()

	redisClient = redis.NewClient(opt)
}
