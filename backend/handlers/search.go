package handlers

import (
	"hero-hunter/config"
	"hero-hunter/models"
	"hero-hunter/services"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// SearchHandler handles the search functionality
func SearchHandler(c *gin.Context) {
	// Get the query, page, and limit from the query parameters
	query := c.DefaultQuery("query", "")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "12")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 12
	}

	// Check Redis cache for the query
	closestMatchQuery, found := services.GetCacheClosestMatch(query)
	if found {
		// If a closest match exists in Redis, update the query to the closest match
		log.Printf("Found closest match in Redis for query '%s': %s", query, closestMatchQuery)
		query = closestMatchQuery
	}

	// Check Redis cache for the query
	cacheData, found := services.GetCacheData(query)
	if found {
		log.Println("Cache hit. Returning cached data.")
		paginatedData := paginateCombined(cacheData.Superheroes, cacheData.Movies, page, limit)
		totalItems := len(cacheData.Superheroes) + len(cacheData.Movies)
		c.JSON(http.StatusOK, models.SearchResponse{
			Superheroes: paginatedData.Superheroes,
			Movies:      paginatedData.Movies,
			TotalPages:  calculateTotalPages(totalItems, limit),
		})
		return
	}

	// Check MongoDB cache
	collection := config.Client.Database("superheroDB").Collection("search_results")
	var mongoData models.CacheData
	err = collection.FindOne(config.MongoCtx, bson.M{"query": query}).Decode(&mongoData)
	if err == nil {
		log.Println("MongoDB hit. Returning data from MongoDB.")
		services.SetCacheData(query, &mongoData)
		paginatedData := paginateCombined(mongoData.Superheroes, mongoData.Movies, page, limit)
		totalItems := len(mongoData.Superheroes) + len(mongoData.Movies)
		c.JSON(http.StatusOK, models.SearchResponse{
			Superheroes: paginatedData.Superheroes,
			Movies:      paginatedData.Movies,
			TotalPages:  calculateTotalPages(totalItems, limit),
		})
		return
	}

	// If no superheroes found for the query, calculate Levenshtein distance
	superheroes := services.FetchSuperheroes(query)
	if len(superheroes) == 0 {
		// Load superhero names from JSON file
		superheroNames, err := loadSuperheroNames()
		if err != nil {
			log.Printf("Error loading superhero names: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading superhero names."})
			return
		}

		// Find the closest match using Levenshtein distance
		closestMatch := query
		minDistance := math.MaxInt

		for _, name := range superheroNames {
			distance := LevenshteinDistance(query, name)
			if distance < minDistance {
				minDistance = distance
				closestMatch = name
			}
		}

		log.Printf("Closest match for query '%s': %s", query, closestMatch)

		// Store the closest match in Redis with the actual query
		services.SetCacheClosestMatch(query, closestMatch)

		// Fetch superheroes and movies for the closest match
		superheroes = services.FetchSuperheroes(closestMatch)
		movies := services.FetchMovies(closestMatch)

		// Save to cache
		cacheData := &models.CacheData{
			Query:       closestMatch,
			Superheroes: superheroes,
			Movies:      movies,
		}
		services.SetCacheData(closestMatch, cacheData)

		// Save to MongoDB
		_, err = collection.InsertOne(config.MongoCtx, cacheData)
		if err != nil {
			log.Printf("Error inserting data into MongoDB: %v", err)
		}

		// Paginate the combined results
		paginatedData := paginateCombined(superheroes, movies, page, limit)
		totalItems := len(superheroes) + len(movies)
		c.JSON(http.StatusOK, models.SearchResponse{
			Superheroes: paginatedData.Superheroes,
			Movies:      paginatedData.Movies,
			TotalPages:  calculateTotalPages(totalItems, limit),
		})
		return
	}

	// If superheroes were found directly, proceed with paginating and returning the results
	movies := services.FetchMovies(query)
	if len(superheroes) > 0 || len(movies) > 0 {
		cacheData := &models.CacheData{
			Query:       query,
			Superheroes: superheroes,
			Movies:      movies,
		}

		// Save to Redis cache
		services.SetCacheData(query, cacheData)

		// Save to MongoDB
		_, err := collection.InsertOne(config.MongoCtx, cacheData)
		if err != nil {
			log.Printf("Error inserting data into MongoDB: %v", err)
		}

		// Paginate the combined results
		paginatedData := paginateCombined(superheroes, movies, page, limit)
		totalItems := len(superheroes) + len(movies)
		c.JSON(http.StatusOK, models.SearchResponse{
			Superheroes: paginatedData.Superheroes,
			Movies:      paginatedData.Movies,
			TotalPages:  calculateTotalPages(totalItems, limit),
		})
	} else {
		c.JSON(http.StatusOK, models.SearchResponse{
			Superheroes: nil,
			Movies:      nil,
			TotalPages:  0,
		})
	}
}
