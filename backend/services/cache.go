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
