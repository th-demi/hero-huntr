package services

import (
	"context"
	"encoding/json"
	"hero-hunter/config"
	"hero-hunter/models"
	"log"
	"time"

	"github.com/go-redis/redis/v8" // Import this for redis.Nil
)

var ctx = context.Background()

// GetCacheData fetches the cached data from Redis
func GetCacheData(query string) (*models.CacheData, bool) {
	cacheKey := "search:" + query
	cachedResult, err := config.Rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// Cache miss
		return nil, false
	} else if err != nil {
		log.Printf("Error fetching from cache: %v", err)
		return nil, false
	}

	var cacheData models.CacheData
	err = json.Unmarshal([]byte(cachedResult), &cacheData)
	if err != nil {
		log.Printf("Error unmarshalling cache data: %v", err)
		return nil, false
	}
	return &cacheData, true
}

// SetCacheData saves the cache data to Redis
func SetCacheData(query string, data *models.CacheData) {
	cacheKey := "search:" + query
	cacheDataJSON, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling cache data: %v", err)
		return
	}
	config.Rdb.Set(ctx, cacheKey, string(cacheDataJSON), 10*time.Minute)
}
func GetCacheClosestMatch(query string) (string, bool) {
	// Redis key for closest match
	redisKey := "closest_match:" + query

	// Fetch from Redis
	val, err := config.Rdb.Get(context.Background(), redisKey).Result()
	if err == redis.Nil {
		// If the key doesn't exist in Redis
		return "", false
	} else if err != nil {
		log.Printf("Error fetching from Redis: %v", err)
		return "", false
	}

	// Return the closest match
	return val, true
}

// SetCacheClosestMatch stores the closest match for the query in Redis
func SetCacheClosestMatch(query, closestMatch string) error {
	// Redis key for closest match
	redisKey := "closest_match:" + query

	// Store the closest match in Redis
	err := config.Rdb.Set(context.Background(), redisKey, closestMatch, 24*time.Hour).Err()
	if err != nil {
		log.Printf("Error storing closest match in Redis: %v", err)
		return err
	}

	// Successfully stored
	return nil
}
