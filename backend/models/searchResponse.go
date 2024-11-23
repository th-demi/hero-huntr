package models

type SearchResponse struct {
	Superheroes []Superhero `json:"superheroes"`
	Movies      []Movie     `json:"movies"`
	TotalPages  int         `json:"totalPages"`
}
