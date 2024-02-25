# URL Shortener Service

This URL Shortener Service is a simple, lightweight application written in Go, using the Gin web framework and GORM for ORM. It allows users to create shorter aliases for long URLs, similar to services like bit.ly or TinyURL.

## Project Structure

The project is organized into several key directories:

- `/models` - Contains the data models for the application.
- `/handlers` - Houses the HTTP handlers that process requests for creating, redirecting, and deleting short URLs.
- `/utils` - Includes utility functions
- `/tests` - Contains unit tests for utility functions and handlers.

Additionally, the project's entry point is located at the root:

- `main.go` - Initializes the application.
