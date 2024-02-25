package models

import (
	"gorm.io/gorm"
)

// Database model for storing url data
type URL struct {
	gorm.Model
	ShortURL    string `gorm:"unique"`
	LongURL     string
}
