package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"url-shortening-service/models"
	"url-shortening-service/utils"
	"gorm.io/gorm"
)

// CreateShortURL - Generates and stores a new short URL
func CreateShortURL(c *gin.Context, db *gorm.DB) {
	// Extract request json body to get LongURL value
	var input struct {
		LongURL string `json:"longURL"`
		Expiry   string `json:"expiry"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Generates short unique url using rand and time generator
	shortURL := utils.GenerateShortURL()

	// Parse Expiry field to get expiration time
	expiryDuration, err := time.ParseDuration(input.Expiry)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry format"})
		return
	}
	expiryTime := time.Now().Add(expiryDuration)

	url := models.URL{
		ShortURL: shortURL,
		LongURL:  input.LongURL,
		Expiry:   &expiryTime,
	}

	// Store url data to the url database model
	if result := db.Create(&url); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create short URL"})
		return
	}

	utils.Logger.Println("Created short URL:", shortURL)
	utils.Logger.Println("Created long URL:", url)

	c.JSON(http.StatusOK, gin.H{"shortURL": shortURL})
}

// RedirectURL - Redirects short URL to original URL
func RedirectURL(c *gin.Context, db *gorm.DB) {
	// Extract the short URL identifier from the request parameters.
	shortURL := c.Param("shortURL")
	var url models.URL

	// Query the database for the URL model for short URL.
	if result := db.Where("short_url = ?", shortURL).First(&url); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	} else {
		 // Check if the URL has expired
		 if url.Expiry != nil && time.Now().After(*url.Expiry) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "This URL has expired"})
            return
        }
	}

	 // Check if the LongURL field is empty.
	if url.LongURL == "" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Long URL is empty"})
        return
    }

	url.AccessCount++
	db.Save(&url)

	utils.Logger.Println("Redirected short URL:", shortURL)
	utils.Logger.Println("Redirected longURL URL:", url.LongURL)

	c.Redirect(http.StatusMovedPermanently, url.LongURL)
}

// DeleteURL - Deletes a short URL
func DeleteURL(c *gin.Context, db *gorm.DB) {
	// Extract the short URL identifier from the request parameters.
	shortURL := c.Param("shortURL")

	// Delete short url from the database
	if result := db.Where("short_url = ?", shortURL).Delete(&models.URL{}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete URL"})
		return
	}

	utils.Logger.Println("Deleted short URL:", shortURL)
	c.JSON(http.StatusOK, gin.H{"message": "URL deleted"})
}

// URLMetrics fetches and returns metrics for a given short URL.
func URLMetrics(c *gin.Context, db *gorm.DB) {
	// Extract the short URL identifier from the request parameters.
	shortURL := c.Param("shortURL")
	var url models.URL

	// Query the database for the URL model for short URL.
	if result := db.Where("short_url = ?", shortURL).First(&url); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	// Return the metrics. You can extend this to include more detailed metrics.
	c.JSON(http.StatusOK, gin.H{
		"shortURL":    url.ShortURL,
		"longURL":     url.LongURL,
		"accessCount": url.AccessCount,
		"created_at":  url.CreatedAt,
		"expiry":      url.Expiry,
	})
}