package handlers

import (
	"encoding/json"
	"hero-hunter/config"
	"hero-hunter/models"
	"hero-hunter/services"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// LevenshteinDistance calculates the Levenshtein distance between two strings
func LevenshteinDistance(a, b string) int {
	lena, lenb := len(a), len(b)
	dp := make([][]int, lena+1)

	for i := range dp {
		dp[i] = make([]int, lenb+1)
	}

	for i := 0; i <= lena; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= lenb; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= lena; i++ {
		for j := 1; j <= lenb; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			dp[i][j] = int(math.Min(
				float64(dp[i-1][j-1]+cost),
				math.Min(
					float64(dp[i-1][j]+1),
					float64(dp[i][j-1]+1),
				),
			))
		}
	}

	return dp[lena][lenb]
}

// Load superhero names from hero_names.json
func loadSuperheroNames() ([]string, error) {
	file, err := os.ReadFile("./handlers/hero_names.json")
	if err != nil {
		return nil, err
	}

	var names []string
	err = json.Unmarshal(file, &names)
	if err != nil {
		return nil, err
	}
	return names, nil
}

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
