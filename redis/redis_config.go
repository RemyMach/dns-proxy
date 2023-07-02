package redis

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		log.Fatal("Redis client is not initialized")
	}
	return redisClient
}

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,	
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
}

func GetStartedLambdaEvent(channel string) *redis.PubSub {
	pubsub := redisClient.Subscribe(channel)
	_, err := pubsub.Receive()
	if err != nil {
		log.Fatal(err)
	}
	return pubsub
}