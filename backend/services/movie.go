// services/movie.go

package services

import (
	"hero-hunter/models"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

func FetchMovies(query string) []models.Movie {
	client := resty.New()

	// Get the OMDB API key from environment variable
	apiKey := os.Getenv("OMDB_API_KEY")
	if apiKey == "" {
		log.Fatal("OMDB_API_KEY is not set in the .env file")
	}

	// Construct the URLs for both movie and series searches
	urlMovie := "http://www.omdbapi.com/?apikey=" + apiKey + "&s=" + query + "&type=movie"
	urlSeries := "http://www.omdbapi.com/?apikey=" + apiKey + "&s=" + query + "&type=series"

	// Create structs to hold the API responses
	var resultMovie struct {
		Response string `json:"Response"`
		Search   []struct {
			Title  string `json:"Title"`
			Poster string `json:"Poster"` // Poster URL
			Year   string `json:"Year"`   // Year of release
		} `json:"Search"`
	}
	var resultSeries struct {
		Response string `json:"Response"`
		Search   []struct {
			Title  string `json:"Title"`
			Poster string `json:"Poster"` // Poster URL
			Year   string `json:"Year"`   // Year of release
		} `json:"Search"`
	}

	// Send the GET requests for both movies and series
	_, errMovie := client.R().SetResult(&resultMovie).Get(urlMovie)
	_, errSeries := client.R().SetResult(&resultSeries).Get(urlSeries)

	// Handle any errors from the API request
	if errMovie != nil {
		log.Printf("Error fetching movie results: %v", errMovie)
	}

	if errSeries != nil {
		log.Printf("Error fetching series results: %v", errSeries)
	}

	// Create a slice of movies and series with only Title, Poster URL, and Year
	var combinedResults []models.Movie

	if resultMovie.Response == "True" {
		for _, movie := range resultMovie.Search {
			combinedResults = append(combinedResults, models.Movie{
				Title:  movie.Title,
				Poster: movie.Poster, // Use the Poster URL from the API response
				Year:   movie.Year,   // Use the Year from the API response
			})
		}
	} else {
		log.Printf("No movies found for query: %s", query)
	}

	if resultSeries.Response == "True" {
		for _, series := range resultSeries.Search {
			combinedResults = append(combinedResults, models.Movie{
				Title:  series.Title,
				Poster: series.Poster, // Use the Poster URL from the API response
				Year:   series.Year,   // Use the Year from the API response
			})
		}
	} else {
		log.Printf("No series found for query: %s", query)
	}

	// Return the combined list of movies and series
	return combinedResults
}
