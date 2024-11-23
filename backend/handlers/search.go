package handlers

import (
	"hero-hunter/config"
	"hero-hunter/models"
	"hero-hunter/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func SearchHandler(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}

	// Check Redis cache
	cacheData, found := services.GetCacheData(query)
	if found {
		log.Println("Cache hit. Returning cached data.")
		c.JSON(http.StatusOK, models.SearchResponse{
			Superheroes: cacheData.Superheroes,
			Movies:      cacheData.Movies,
			TotalPages:  5,
		})
		return
	}
	log.Println("Cache miss. Checking MongoDB.")

	// Cache miss. Check MongoDB
	collection := config.Client.Database("superheroDB").Collection("search_results")
	var mongoData models.CacheData
	err := collection.FindOne(config.MongoCtx, bson.M{"query": query}).Decode(&mongoData)
	if err == nil {
		// MongoDB hit
		log.Println("MongoDB hit. Returning data from MongoDB.")
		// Cache to Redis for next time
		services.SetCacheData(query, &mongoData) // Pass pointer here
		c.JSON(http.StatusOK, models.SearchResponse{
			Superheroes: mongoData.Superheroes,
			Movies:      mongoData.Movies,
			TotalPages:  5,
		})
		return
	}
	log.Println("MongoDB miss. Making external API calls.")

	// Fetch data from external APIs
	superheroes := services.FetchSuperheroes(query)
	movies := services.FetchMovies(query)

	// Only cache and store the data if there is valid data for both
	if len(superheroes) > 0 || len(movies) > 0 {
		// Create a pointer to CacheData (use pointer to avoid copying large struct)
		cacheData := &models.CacheData{
			Query:       query,
			Superheroes: superheroes,
			Movies:      movies,
		}

		// Cache the fetched data in Redis and MongoDB
		services.SetCacheData(query, cacheData) // Pass pointer here

		// Save to MongoDB if we have valid data
		_, err := collection.InsertOne(config.MongoCtx, cacheData)
		if err != nil {
			log.Printf("Error inserting data into MongoDB: %v", err)
		} else {
			log.Println("Data saved to MongoDB.")
		}

		// Return the response with valid data
		c.JSON(http.StatusOK, models.SearchResponse{
			Superheroes: superheroes,
			Movies:      movies,
			TotalPages:  5,
		})
	} else {
		// If no valid data was returned, return an empty response
		log.Println("No valid data found for the query.")
		c.JSON(http.StatusOK, models.SearchResponse{
			Superheroes: nil,
			Movies:      nil,
			TotalPages:  0,
		})
	}
}
