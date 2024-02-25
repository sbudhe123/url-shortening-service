package main

import (
	"github.com/gin-gonic/gin"
	"url-shortening-service/utils"
	"url-shortening-service/handlers"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"url-shortening-service/models"
)

func main() {
	// Initialize the logger to log information, warnings, and errors.
	utils.InitLogger()

	// Create a connection to the SQLite database. The database file is named "urls.db".
	db, err := gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})
	if err != nil {
		utils.Logger.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate will automatically create (if they don't exist) URL table if it does not exist depeneding on the models.URL fields
	if err := db.AutoMigrate(&models.URL{}); err != nil {
		utils.Logger.Fatalf("Could not migrate the database: %v", err)
	}

	// Initialize the Gin router
	r := gin.Default()

	// Define routes and respective handlers for serving request
	r.POST("/create", func(c *gin.Context) { handlers.CreateShortURL(c, db) })
	r.GET("/:shortURL", func(c *gin.Context) { handlers.RedirectURL(c, db) })

	// staret the server
	r.Run(":8080")
}