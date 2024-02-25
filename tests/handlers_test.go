package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	
	"url-shortening-service/handlers"
	"url-shortening-service/models"
	"url-shortening-service/utils"
)

var logBuffer bytes.Buffer

// SetupTestRouter initializes a Gin router with registered handlers for testing.
// It returns the router and a reference to the in-memory database.
func setupTestRouter() (*gin.Engine, *gorm.DB) {
	
	// Redirect logger output to the buffer
	utils.Logger.SetOutput(&logBuffer)
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&models.URL{})

	r.POST("/create", func(c *gin.Context) { handlers.CreateShortURL(c, db) })
	r.GET("/:shortURL", func(c *gin.Context) { handlers.RedirectURL(c, db) })
	r.DELETE("/:shortURL", func(c *gin.Context) { handlers.DeleteURL(c, db) })
	r.GET("/metrics/:shortURL", func(c *gin.Context) { handlers.URLMetrics(c, db) })

	return r, db
}

func tearDownTest(db *gorm.DB) {
	logBuffer.Reset()
	// Reset logger output to default if necessary
	utils.Logger.SetOutput(os.Stderr)
	db.Exec("DELETE FROM urls") 
}

type testCase struct {
    description    string
    method         string
    url            string
    body           string
    expectedStatus int
    expectedBody   string // Optional: for more specific response body checks
}

func TestCreateShortURL(t *testing.T) {
    r, db := setupTestRouter()
    defer tearDownTest(db)

    testCases := []testCase{
        {
            description:    "Valid input",
            method:         "POST",
            url:            "/create",
            body:           `{"longURL": "https://example.com", "expiry": "1m"}`,
            expectedStatus: http.StatusOK,
        },
        {
            description:    "Invalid input - Empty URL",
            method:         "POST",
            url:            "/create",
            body:           `{"longURL": ""}`,
            expectedStatus: http.StatusBadRequest,
        },
		{
            description:    "Invalid url",
            method:         "POST",
            url:            "/create",
            body:           `{"longURL": "invalid"}`,
            expectedStatus: http.StatusBadRequest,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.description, func(t *testing.T) {
            req, _ := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            r.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)
        })
    }
}

func TestRedirectURL(t *testing.T) {
	r, db := setupTestRouter()
	defer tearDownTest(db)

	// Create a URL entry directly in the database
	db.Create(&models.URL{ShortURL: "test123", LongURL: "https://example.com"})
	
    testCases := []testCase{
        {
            description:    "Redirect existing URL",
            method:         "GET",
            url:            "/test123",
            expectedStatus: http.StatusMovedPermanently,
        },
        {
            description:    "URL does not exist",
            method:         "GET",
            url:            "/nonexistent",
            expectedStatus: http.StatusNotFound,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.description, func(t *testing.T) {
            req, _ := http.NewRequest(tc.method, tc.url, nil)

            w := httptest.NewRecorder()
            r.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)
        })
    }
}

func TestDeleteURL(t *testing.T) {
	r, db := setupTestRouter()
	defer tearDownTest(db)

	// Create a URL entry directly in the database
	db.Create(&models.URL{ShortURL: "test123", LongURL: "https://example.com"})
	
    testCases := []testCase{
        {
            description:    "Delete existing URL",
            method:         "DELETE",
            url:            "/test123",
            expectedStatus: http.StatusOK,
        },
    }
    for _, tc := range testCases {
        t.Run(tc.description, func(t *testing.T) {
            req, _ := http.NewRequest(tc.method, tc.url, nil)

            w := httptest.NewRecorder()
            r.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)
        })
    }
}

func TestURLMetrics(t *testing.T) {
	r, db := setupTestRouter()
	defer tearDownTest(db)

	// Create a URL entry directly in the database
	db.Create(&models.URL{ShortURL: "metrics123", LongURL: "https://example.com", AccessCount: 5})
	testCases := []testCase{
        {
            description:    "Metrics for existing URL",
            method:         "GET",
            url:            "/metrics/metrics123",
            expectedStatus: http.StatusOK,
        },
        {
            description:    "Metrics for non-existent URL",
            method:         "GET",
            url:            "/metrics/nonexistent",
            expectedStatus: http.StatusNotFound,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.description, func(t *testing.T) {
            req, _ := http.NewRequest(tc.method, tc.url, nil)

            w := httptest.NewRecorder()
            r.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)
            if tc.expectedStatus == http.StatusOK {
                var response map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.NoError(t, err)
                assert.Equal(t, float64(5), response["accessCount"]) // json.Unmarshal decodes numbers as floats
            }
        })
    }
}