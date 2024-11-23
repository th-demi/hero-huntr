package config

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

// Initialize Redis
func InitRedis() {
	// Get Redis configuration from environment variables
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD") // Can be empty if no password is required

	if redisHost == "" || redisPort == "" {
		log.Fatal("REDIS_HOST or REDIS_PORT environment variable is not set")
	}

	// Initialize Redis client
	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword, // Default is no password
		DB:       0,             // Default DB is 0
	})

	// Test Redis connection
	_, err := Rdb.Ping(Rdb.Context()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
}
