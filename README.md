# Zione Backend API

A production-grade RESTful API built with Go, Gin, MySQL, and GORM for the personal website zionechain.cfd.

## Features

- User authentication with JWT and role-based access control
- Project showcase with categories, tags, and media attachments
- Blog system with categories, media support, and SEO-optimized structure
- Admin panel functionality via API endpoints

## Architecture

The application follows standard Go project layout and clean architecture principles:

```
.
├── api/            # API documentation (Swagger)
├── cmd/            # Application entry points
├── configs/        # Configuration files
├── internal/       # Application core code (private)
│   ├── controllers/  # HTTP handlers
│   ├── models/       # Database models
│   ├── middleware/   # HTTP middleware
│   ├── database/     # Database connection and migrations
│   ├── services/     # Business logic
│   └── utils/        # Utility functions
├── pkg/            # Public libraries
├── tests/          # Integration and end-to-end tests
└── deployments/    # Deployment configurations
```

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin
- **Database**: MySQL
- **ORM**: GORM
- **Authentication**: JWT
- **Documentation**: Swagger via Swaggo

## Setup Instructions

### Prerequisites

- Go 1.20+
- Docker and Docker Compose (optional)
- Make (optional)

### Local Development

#### Option 1: Using Docker

1. Clone or download the project
   ```bash
   # Download the project and navigate to the directory
   cd zione-backend
   ```

2. Create `.env` file from example
   ```bash
   cp .env.example .env
   # Edit the .env file with your preferred settings
   ```

3. Run the application with Docker Compose
   ```bash
   docker-compose up --build
   ```

4. Access the API at http://localhost:8080/api
   
5. Access Swagger documentation at http://localhost:8080/swagger/index.html

#### Option 2: Without Docker

1. Clone or download the project
   ```bash
   # Download the project and navigate to the directory
   cd zione-backend
   ```

2. Create `.env` file (optional, application uses default values if not provided)
   ```bash
   # Example environment variables
   APP_ENV=development
   APP_PORT=3000
   ```

3. Run the application directly
   ```bash
   go run cmd/api/main.go
   ```

4. Access the API at http://localhost:3000/
   
5. Access API welcome endpoint at http://localhost:3000/api

### Running Tests

```bash
# Run unit tests
go test ./internal/...

# Run integration tests (requires Docker)
docker-compose -f docker-compose.test.yml up --build
```

## API Endpoints

| Method | Route                    | Description              | Access |
| ------ | ------------------------ | ------------------------ | ------ |
| POST   | /api/auth/login          | Login via phone/password | Public |
| POST   | /api/auth/register       | Register new user        | Public |
| GET    | /api/projects            | Get list of projects     | Public |
| POST   | /api/projects            | Create project           | Admin  |
| GET    | /api/blog                | Get blog posts           | Public |
| POST   | /api/blog                | Create blog post         | Admin  |
| GET    | /api/categories/projects | Get project categories   | Public |
| GET    | /api/categories/blog     | Get blog categories      | Public |

For a complete list of endpoints, refer to the Swagger documentation.

## Deployment

The application can be deployed using Docker Compose or any container orchestration platform like Kubernetes.

```bash
# Production deployment
docker-compose -f docker-compose.prod.yml up -d
``` 