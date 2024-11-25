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

	// Get the optional filters: powerMin, powerMax, and alignment
	powerMinStr := c.DefaultQuery("powerMin", "")
	powerMaxStr := c.DefaultQuery("powerMax", "")
	alignment := c.DefaultQuery("alignment", "")

	// Convert powerMin and powerMax to integers if provided
	var powerMin, powerMax int
	if powerMinStr != "" {
		powerMin, err = strconv.Atoi(powerMinStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid powerMin value"})
			return
		}
	}
	if powerMaxStr != "" {
		powerMax, err = strconv.Atoi(powerMaxStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid powerMax value"})
			return
		}
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
		// Apply filters to the cached superheroes
		filteredSuperheroes := filterSuperheroes(cacheData.Superheroes, powerMin, powerMax, alignment)
		paginatedData := paginateCombined(filteredSuperheroes, cacheData.Movies, page, limit)
		totalItems := len(filteredSuperheroes) + len(cacheData.Movies)
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
		// Apply filters to the MongoDB superheroes
		filteredSuperheroes := filterSuperheroes(mongoData.Superheroes, powerMin, powerMax, alignment)
		paginatedData := paginateCombined(filteredSuperheroes, mongoData.Movies, page, limit)
		totalItems := len(filteredSuperheroes) + len(mongoData.Movies)
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

		// Apply filters to the superheroes
		filteredSuperheroes := filterSuperheroes(superheroes, powerMin, powerMax, alignment)

		// Save to cache
		cacheData := &models.CacheData{
			Query:       closestMatch,
			Superheroes: filteredSuperheroes,
			Movies:      movies,
		}
		services.SetCacheData(closestMatch, cacheData)

		// Save to MongoDB
		_, err = collection.InsertOne(config.MongoCtx, cacheData)
		if err != nil {
			log.Printf("Error inserting data into MongoDB: %v", err)
		}

		// Paginate the combined results
		paginatedData := paginateCombined(filteredSuperheroes, movies, page, limit)
		totalItems := len(filteredSuperheroes) + len(movies)
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
		// Apply filters to the superheroes
		filteredSuperheroes := filterSuperheroes(superheroes, powerMin, powerMax, alignment)

		cacheData := &models.CacheData{
			Query:       query,
			Superheroes: filteredSuperheroes,
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
		paginatedData := paginateCombined(filteredSuperheroes, movies, page, limit)
		totalItems := len(filteredSuperheroes) + len(movies)
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

// filterSuperheroes filters the list of superheroes based on powerMin, powerMax, and alignment
// filterSuperheroes filters the list of superheroes based on powerMin, powerMax, and alignment
func filterSuperheroes(superheroes []models.Superhero, powerMin, powerMax int, alignment string) []models.Superhero {
	var filtered []models.Superhero
	for _, hero := range superheroes {
		// Convert hero.Power to an integer (if it can be converted)
		power, err := strconv.Atoi(hero.Power)
		if err != nil {
			// If conversion fails, skip this superhero or handle the error accordingly
			log.Printf("Invalid power value for superhero '%s': %s", hero.Name, hero.Power)
			continue
		}

		// Apply the filters
		if (powerMin == 0 || power >= powerMin) && (powerMax == 0 || power <= powerMax) {
			if alignment == "" || hero.Alignment == alignment {
				filtered = append(filtered, hero)
			}
		}
	}
	return filtered
}
