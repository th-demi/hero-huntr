package handlers

import (
	"encoding/json"
	"hero-hunter/models"
	"math"
	"os"
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

// calculateTotalPages calculates the total number of pages for pagination
func calculateTotalPages(totalItems, limit int) int {
	return (totalItems + limit - 1) / limit
}
