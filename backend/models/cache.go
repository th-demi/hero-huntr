// models/cache.go

package models

// Superhero struct with name, image, power, and alignment
type Superhero struct {
	Name      string `json:"name"`
	Image     string `json:"image"`
	Power     string `json:"power"`     // Added Power field
	Alignment string `json:"alignment"` // Added Alignment field
}

type Movie struct {
	Title  string `json:"Title"`
	Poster string `json:"Poster"`
	Year   string `json:"Year"` // Added Year field
}

type CacheData struct {
	Query       string      `json:"query"`
	Superheroes []Superhero `json:"superheroes"`
	Movies      []Movie     `json:"movies"`
}
