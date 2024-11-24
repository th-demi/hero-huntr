package handlers

import (
	"hero-hunter/config"
	"hero-hunter/models"
	"hero-hunter/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

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

	// Fetch data from external APIs if not found in cache
	superheroes := services.FetchSuperheroes(query)
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
		log.Println("Calling paginateCombined function...")
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

// paginateCombined handles pagination for combined superhero and movie results
func paginateCombined(superheroes []models.Superhero, movies []models.Movie, page, limit int) *models.CacheData {
	// Create a combined slice of interfaces to handle both types
	combined := make([]interface{}, 0, len(superheroes)+len(movies))

	// Add all items to the combined slice
	for _, hero := range superheroes {
		combined = append(combined, hero)
	}
	for _, movie := range movies {
		combined = append(combined, movie)
	}

	// Calculate start and end indices for the requested page
	start := (page - 1) * limit
	end := start + limit
	if end > len(combined) {
		end = len(combined)
	}
	if start >= len(combined) {
		return &models.CacheData{
			Superheroes: []models.Superhero{},
			Movies:      []models.Movie{},
		}
	}

	// Get the slice for the current page
	pageSlice := combined[start:end]

	// Separate back into superheroes and movies
	var paginatedHeroes []models.Superhero
	var paginatedMovies []models.Movie

	for _, item := range pageSlice {
		switch v := item.(type) {
		case models.Superhero:
			paginatedHeroes = append(paginatedHeroes, v)
		case models.Movie:
			paginatedMovies = append(paginatedMovies, v)
		}
	}

	return &models.CacheData{
		Superheroes: paginatedHeroes,
		Movies:      paginatedMovies,
	}
}

func calculateTotalPages(totalItems, limit int) int {
	return (totalItems + limit - 1) / limit
}
