package models

import (
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	TMBD_ID string `gorm:"unique;not null"`
	Title    string `gorm:"unique;not null"`
	PosterPath string `gorm:"not null"`
}