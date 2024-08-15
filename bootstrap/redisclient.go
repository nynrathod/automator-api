package bootstrap

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

// Initialize Redis client
var redisClient *redis.Client

func InitRedis() {
	var err error
	opt, _ := redis.ParseURL("rediss://default:AcvqAAIjcDFhYzk0MWRhOGY0MGU0NTdiODMwODg3ZGE0NTZhMTA4ZnAxMA@stirred-satyr-52202.upstash.io:6379")
	redisClient = redis.NewClient(opt)

	//redisClient = redis.NewClient(&redis.Options{
	//	Addr: "redis-16336.c264.ap-south-1-1.ec2.cloud.redislabs.com:16336",
	//	//Username: "",
	//	Password: "XJecZsRZaacUxVVp9RQxyeXQGKvtgMN7",
	//	DB:       0,
	//})

	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")
}
