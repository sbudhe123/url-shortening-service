package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"url-shortening-service/models"
	"url-shortening-service/utils"
	"gorm.io/gorm"
)

// CreateShortURL - Generates and stores a new short URL
func CreateShortURL(c *gin.Context, db *gorm.DB) {
	// Extract request json body to get LongURL value
	var input struct {
		LongURL string `json:"longURL"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Generates short unique url using rand and time generator
	shortURL := utils.GenerateShortURL()

	// Store url data to the url database model
	url := models.URL{ShortURL: shortURL, LongURL: input.LongURL}
	if result := db.Create(&url); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create short URL"})
		return
	}

	utils.Logger.Println("Created short URL:", shortURL)
	utils.Logger.Println("Created long URL:", input.LongURL)

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
	}

	 // Check if the LongURL field is empty.
	if url.LongURL == "" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Long URL is empty"})
        return
    }

	utils.Logger.Println("Redirected short URL:", shortURL)
	utils.Logger.Println("Redirected longURL URL:", url.LongURL)

	c.Redirect(http.StatusMovedPermanently, url.LongURL)
}