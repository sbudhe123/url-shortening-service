package utils

import (
	"log"
	"os"
)

// LogFile defines the path to the log file.
const LogFile = "./url-shortening-service.log"

var Logger *log.Logger

func InitLogger() {
	file, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	Logger = log.New(file, "URL_SHORTENER: ", log.Ldate|log.Ltime|log.Lshortfile)
}