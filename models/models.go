package models

import (
	"gorm.io/gorm"
	"time"
)

// Database model for storing url data
type URL struct {
	gorm.Model
	ShortURL    string `gorm:"unique"`
	LongURL     string
	AccessCount int64 // Access count for URL
	Expiry      *time.Time // Expirty for URL
}
