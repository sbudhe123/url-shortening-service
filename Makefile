.PHONY: build run set-env

# Binary output name
BINARY_NAME=url-shortening-service

build:
	@echo "Building the application..."
	@go build -o $(BINARY_NAME) main.go

run: build
	@echo "Running the application..."
	./$(BINARY_NAME)

# Run unit tests.
test:
	@echo "Running tests..."
	@go test ./... -v

# Generate and display test coverage in the CLI.
coverage:
	@echo "Generating test coverage report..."
	@go test ./... -coverprofile=coverage.out 
	@echo "Displaying test coverage..."
	@go tool cover -func=coverage.out
	@echo "Total coverage on CLI"
	@go tool cover -func=coverage.out | grep total | awk '{print $$3}'
