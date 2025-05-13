# Zione Backend API

A production-grade RESTful API built with Go for the personal website zionechain.cfd.

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
- **Framework**: Standard HTTP package
- **Authentication**: JWT

## Setup Instructions

### Prerequisites

- Go 1.20+
- Make (optional)

### Local Development

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
```

## API Endpoints

| Method | Route                    | Description              | Access |
| ------ | ------------------------ | ------------------------ | ------ |
| GET    | /                        | API status               | Public |
| GET    | /health                  | Health check             | Public |
| GET    | /api                     | API welcome              | Public |

## Deployment

The application can be deployed to any Go-compatible hosting platform. 