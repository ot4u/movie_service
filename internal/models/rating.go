package models

import "gorm.io/gorm"

type Rating struct {
	gorm.Model
	UserID  uint
	MovieID uint
	Score   int // от 1 до 10

	User  User  `gorm:"foreignKey:UserID"`
	Movie Movie `gorm:"foreignKey:MovieID"`
}

type RatedMovieResponse struct {
	TMDB_ID     int    `json:"tmdb_id"`
	Title       string `json:"title"`
	PosterPath  string `json:"poster_path"`
	ReleaseDate string `json:"release_date"`
	Score       int    `json:"score"`
}