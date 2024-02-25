package utils

import (
    "math/rand"
    "time"
)

// GenerateShortURL creates a short URL identifier.
func GenerateShortURL() string {
    rand.Seed(time.Now().UnixNano())
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    b := make([]rune, 8) // Generating an 8-character string.
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}