package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func initRedis() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}
	fmt.Printf("Attempting to connect to Redis at %s\n", redisURL)

	redisClient = redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
	} else {
		fmt.Println("Successfully connected to Redis")
	}
}

func getFromCache(key string) (Latest, error) {
	ctx := context.Background()
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return Latest{}, err
	}

	var latest Latest
	err = json.Unmarshal([]byte(val), &latest)
	if err != nil {
		return Latest{}, err
	}

	return latest, nil
}

func setToCache(key string, value Latest, expiration time.Duration) error {
	ctx := context.Background()
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return redisClient.Set(ctx, key, jsonValue, expiration).Err()
}
