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
- **Database**: MySQL

## Setup Instructions

### Prerequisites

- Go 1.20+
- MySQL (for local database usage)
- Make (optional)

## Local Development

### Option 1: Running on Linux/MacOS

1. Clone or download the project
   ```bash
   # Download the project and navigate to the directory
   git clone https://github.com/arashdm2020/banckend-zione.git
   cd zione-backend
   ```

2. Create `.env` file (optional, application uses default values if not provided)
   ```bash
   # Example environment variables
   APP_ENV=development
   APP_PORT=3000
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=zione_db
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   ```

3. Run the application directly
   ```bash
   go run cmd/api/main.go
   ```

4. Access the API at http://localhost:3000/
   
5. Access API welcome endpoint at http://localhost:3000/api

### Option 2: Running on Windows

1. Clone the repository
   ```powershell
   git clone https://github.com/arashdm2020/banckend-zione.git zione-backend
   cd zione-backend
   ```

2. Create a batch file for running the application
   ```batch
   @echo off
   :: Set environment variables
   set APP_ENV=development
   set APP_PORT=3000
   set APP_HOST=localhost
   set DB_HOST=localhost
   set DB_PORT=3306
   set DB_NAME=zione_db
   set DB_USER=root
   set DB_PASSWORD=your_password
   set DB_CHARSET=utf8mb4

   :: Run the application
   go run cmd/api/main.go
   ```
   Save this as `run.bat` in the project root

3. Install dependencies
   ```powershell
   go mod tidy
   ```

4. Run the application
   ```powershell
   .\run.bat
   ```

5. Access your API at http://localhost:3000/

## Database Setup

### MySQL Setup (Local Development)

1. Install MySQL
   ```bash
   # Ubuntu/Debian
   sudo apt install mysql-server
   
   # MacOS with Homebrew
   brew install mysql
   ```

2. Create a database and user
   ```sql
   CREATE DATABASE zione_db;
   CREATE USER 'zione_user'@'localhost' IDENTIFIED BY 'your_password';
   GRANT ALL PRIVILEGES ON zione_db.* TO 'zione_user'@'localhost';
   FLUSH PRIVILEGES;
   ```

3. Update your `.env` or environment variables with your database credentials

## Testing

### Running Unit Tests

```bash
# Run unit tests
go test ./internal/...
```

### API Testing

You can test the API using various tools:

1. **curl** from command line
   ```bash
   # Test the root endpoint
   curl http://localhost:3000/
   
   # Test the API welcome endpoint
   curl http://localhost:3000/api
   ```

2. **Postman**: Download from [postman.com](https://www.postman.com/downloads/)

3. **Web Browser**: For GET requests, simply open the endpoints in your browser

## API Endpoints

| Method | Route                     | Description              | Access |
| ------ | ------------------------- | ------------------------ | ------ |
| GET    | /                         | API status               | Public |
| GET    | /health                   | Health check             | Public |
| GET    | /api                      | API welcome              | Public |
| POST   | /api/auth/login           | Login via phone/password | Public |
| POST   | /api/auth/register        | Register new user        | Public |
| GET    | /api/projects             | Get list of projects     | Public |
| POST   | /api/projects             | Create project           | Admin  |
| GET    | /api/blog                 | Get blog posts           | Public |
| POST   | /api/blog                 | Create blog post         | Admin  |
| GET    | /api/categories/projects  | Get project categories   | Public |
| GET    | /api/categories/blog      | Get blog categories      | Public |

## Deployment

### Shared Hosting Deployment

1. SSH into your server
   ```bash
   ssh user@your-server -p port
   ```

2. Clone the repository in your www directory
   ```bash
   cd ~/www
   mkdir zione-backend
   cd zione-backend
   git clone https://github.com/arashdm2020/banckend-zione.git .
   ```

3. Set up the database configuration
   ```bash
   nano start.sh
   ```
   Edit the database configuration with your hosting provider's database details

4. Make all scripts executable
   ```bash
   chmod +x *.sh
   ```

5. Run the application
   ```bash
   # Build and run the compiled binary (recommended for production)
   ./build.sh
   ./run.sh
   
   # Or directly run the Go code
   ./start.sh
   ```

6. Stop the application when needed
   ```bash
   ./stop.sh
   ```

7. Access your API at
   ```
   http://your-domain.com:port/
   ```

### Troubleshooting Deployment

1. If you encounter line ending issues on Linux servers when uploading files from Windows:
   ```bash
   sed -i 's/\r$//' *.sh
   ```

2. If the port is already in use:
   ```bash
   # Check which process is using the port
   lsof -i :port_number
   
   # Kill that process
   kill -9 process_id
   ```

3. For viewing logs:
   ```bash
   tail -f app.log
   ```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -m 'Add feature'`
4. Push to the branch: `git push origin feature-name`
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 