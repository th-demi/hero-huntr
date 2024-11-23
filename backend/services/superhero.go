package services

import (
	"hero-hunter/models"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

func FetchSuperheroes(query string) []models.Superhero {
	client := resty.New()

	// Get the access token from environment variable
	accessToken := os.Getenv("SUPERHERO_API_ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatal("SUPERHERO_API_ACCESS_TOKEN is not set in the .env file")
	}

	// Construct the URL to search for superheroes by name with the access token
	url := "https://superheroapi.com/api/" + accessToken + "/search/" + query

	// Create a struct to hold the API response
	var result struct {
		Response string `json:"response"`
		Results  []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Powerstats struct {
				Power string `json:"power"` // Power as a string
			} `json:"powerstats"`
			Biography struct {
				Alignment string `json:"alignment"` // Alignment as a string
			} `json:"biography"`
			Image struct {
				URL string `json:"url"` // Nested struct to extract image URL
			} `json:"image"`
		} `json:"results"`
	}

	// Send the GET request and store the result in 'result'
	_, err := client.R().
		SetResult(&result). // Automatically unmarshals the JSON response into 'result'
		Get(url)

	if err != nil {
		log.Printf("Error fetching superheroes: %v", err)
		return nil
	}

	// Create a slice of Superheroes with name, image URL, power, and alignment
	var superheroes []models.Superhero
	if result.Response == "success" {
		// Loop through the results and add the superhero name, image URL, power, and alignment
		for _, hero := range result.Results {
			superheroes = append(superheroes, models.Superhero{
				Name:      hero.Name,
				Image:     hero.Image.URL,           // Extract the image URL from the nested struct
				Power:     hero.Powerstats.Power,    // Extract the 'power' from nested struct
				Alignment: hero.Biography.Alignment, // Extract the 'alignment' from nested struct
			})
		}
	}

	log.Printf("Fetched %d superheroes", len(superheroes))
	// Return the list of superheroes
	return superheroes
}
