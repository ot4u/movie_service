package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type TMDBMovie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	PosterPath  string `json:"poster_path"`
	ReleaseDate string `json:"release_date"`
}

type TMDBSearchResponse struct {
	Results []TMDBMovie `json:"results"`
}

func SearchMovies(query string) ([]TMDBMovie, error) {
	apiKey := os.Getenv("TMDB_API_KEY")
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s", apiKey, query)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TMDBSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Results, nil
}