package models

import "gorm.io/gorm"

// UserLike — таблица "пользователь лайкнул фильм"
type UserLike struct {
	gorm.Model
	UserID  uint
	MovieID uint
	
	User  User  `gorm:"foreignKey:UserID"`
	Movie Movie `gorm:"foreignKey:MovieID"`
}