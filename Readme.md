# URL Shortener Service

## Introduction

The URL Management Service is a Go-based RESTful API designed for simplifying URL management. This service supports URL shortening, redirection, deletion, and metrics tracking, leveraging the Gin web framework and GORM.

## Project Structure

The project is organized into several key directories:

- `/models` - Contains the data models for the application.
- `/handlers` - Houses the HTTP handlers that process requests for creating, redirecting, and deleting short URLs.
- `/utils` - Includes utility functions
- `/tests` - Contains unit tests for utility functions and handlers.

Additionally, the project's entry point is located at the root:

- `main.go` - Initializes the application.

## Features

- **Shorten URLs**: Convert long URLs into short, user-friendly links.
- **URL Redirection**: Redirect users from short links to the original URLs.
- **Delete URLs**: Remove short URLs from the database.
- **Metrics**: View access counts and other metrics for short URLs.

## Getting Started

### Prerequisites

- Go (version 1.15 or newer)
- SQLite3 (for local development)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/sbudhe123/url-shortening-service
cd url-shortening-service
```

2. **Install Dependencies**

    ```bash
    go get github.com/gin-gonic/gin
    go get gorm.io/gorm
    go get gorm.io/driver/sqlite
    go mod tidy
    ```

### Use make to build, run, test and coverage

Build the application

```bash
make build
```

To start the service

```bash
make run
```

To run tests

```bash
make test
```

For coverage

```bash
make coverage
```

### Usage
After starting the service, you can interact with the API using curl or Postman. Here are some examples:

```bash
curl -X POST http://localhost:8080/create -H 'Content-Type: application/json' -d '{"longURL": "https://example.com"}'
```

### Redirect from a short URL:

Navigate to http://localhost:8080/{shortURL} in your browser or using curl.

### Delete a short url:

```bash
curl -X DELETE http://localhost:8080/delete/{shortURL}
```

### To view metrics

```bash
curl http://localhost:8080/metrics/{shortURL}
```

## Limitations of Current Implementation

1. **Single Instance Database**: The current SQLite3 database may become a bottleneck under high load conditions due to its file-based nature and write-lock limitations.

2. **Lack of Horizontal Scaling**: The application is designed to run as a monolith, making it challenging to scale horizontally across multiple machines or instances.

3. **Synchronous Processing**: All operations, including URL redirection and metrics updates, are processed synchronously, which could lead to increased response times under heavy load.

4. **No Built-In Caching**: Without caching, every request hits the database directly, increasing the load on the database and potentially leading to slower response times as user numbers grow.

5. **Limited Fault Tolerance**: The current setup lacks mechanisms for automatic failover or recovery, making it susceptible to single points of failure.

## Scalability

### Strategies to Scale to Millions of Users

1. **Microservices Architecture**: Refactor the application into microservices to improve scalability and maintainability. Each service can be scaled independently based on demand.

2. **Load Balancing**: Deploy a load balancer in front of the application servers to distribute traffic evenly across multiple instances, improving response times and system availability.

3. **Database Optimization**:
    - Use database sharding to distribute the data across multiple databases, reducing the load on any single database server.
    - Implement read replicas to handle read-heavy loads efficiently.

4. **Caching**:
    - Introduce caching layers, such as Redis, for frequently accessed data like URL mappings, to reduce database load and improve response times.
    - Use edge caching or Content Delivery Networks (CDNs) to cache responses closer to the users, reducing latency.

5. **Asynchronous Processing**: Utilize message queues for asynchronous processing of tasks that do not require immediate completion, such as metrics collection.

6. **Auto-Scaling**: Employ auto-scaling for your application and database layers based on load, ensuring that resources are allocated efficiently and can handle spikes in traffic.

7. **Rate Limiting**: Implement rate limiting to prevent abuse and ensure that the service remains available to all users.

8. **Monitoring and Alerting**: Establish comprehensive monitoring and alerting to detect performance bottlenecks and system failures early, allowing for proactive management of scalability issues.
