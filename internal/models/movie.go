package models

import (
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	TMDB_ID     int    `gorm:"not null"`
	Title       string `gorm:"not null"`
	PosterPath  string
	ReleaseDate string
}

type MovieResponse struct {
	TMDB_ID     int    `json:"tmdb_id"`
	Title       string `json:"title"`
	PosterPath  string `json:"poster_path"`
	ReleaseDate string `json:"release_date"`
}