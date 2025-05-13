# Zione Backend API

A production-grade RESTful API built with Go for the personal website zionechain.cfd.

## Features

- User authentication with JWT and role-based access control
- Project showcase with categories, tags, and media attachments
- Blog system with categories, media support, and SEO-optimized structure
- Admin panel functionality via API endpoints
- Resume/Personal information management with detailed sections

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

## Getting Started with Resume API

Follow these steps to set up and test the Resume API:

### Prerequisites

- Go 1.20 or later
- MySQL 5.7 or later
- Git

### Setup 

1. Clone the repository:
   ```bash
   git clone https://github.com/arashdm2020/banckend-zione.git
   cd banckend-zione
   ```

2. Set up your environment variables by creating a `.env` file:
   ```bash
   # Application settings
   APP_ENV=development
   APP_PORT=3000
   APP_HOST=0.0.0.0
   APP_NAME=zione-backend
   APP_URL=http://localhost:3000
   APP_SECRET=your_secret_key_change_in_production
   
   # Database settings
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=zionec_db
   DB_USER=zionec_user
   DB_PASSWORD=your_db_password
   DB_CHARSET=utf8mb4
   DB_MAX_IDLE_CONNS=10
   DB_MAX_OPEN_CONNS=100
   DB_CONN_MAX_LIFETIME=3600
   
   # JWT settings
   JWT_SECRET=your_jwt_secret_change_in_production
   JWT_ACCESS_TOKEN_EXPIRY=15m
   JWT_REFRESH_TOKEN_EXPIRY=168h
   
   # CORS settings
   CORS_ALLOWED_ORIGINS=http://localhost:3000,https://zionechain.cfd
   CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
   CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization
   
   # Logging settings
   LOG_LEVEL=info
   LOG_FORMAT=json
   
   # TLS settings
   TLS_ENABLED=false
   TLS_CERT_FILE=./certs/server.crt
   TLS_KEY_FILE=./certs/server.key
   ```
   > Note: The `.env` file is excluded from Git. You should always create it locally and never commit it to version control.

3. Create the database:
   ```sql
   CREATE DATABASE zionec_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   CREATE USER 'zionec_user'@'localhost' IDENTIFIED BY 'your_db_password';
   GRANT ALL PRIVILEGES ON zionec_db.* TO 'zionec_user'@'localhost';
   FLUSH PRIVILEGES;
   ```

4. Install dependencies and run the application:
   ```bash
   go mod download
   ./start.sh
   ```
   Or run directly without using the script:
   ```bash
   go run cmd/api/main.go
   ```

The server will start and automatically create all necessary database tables, including the resume-related tables.

### Testing Resume API Endpoints

Here are some sample cURL commands to test the Resume API endpoints:

#### Get All Resume Sections

```bash
curl -X GET http://localhost:3000/api/resume/complete
```

#### Create Personal Information

```bash
curl -X POST http://localhost:3000/api/resume/personal \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "job_title": "Full Stack Developer",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "address": "123 Main St, City, Country",
    "website": "https://johndoe.com",
    "linkedin": "https://linkedin.com/in/johndoe",
    "github": "https://github.com/johndoe",
    "summary": "Experienced developer with 5+ years in web development"
  }'
```

#### Add a Skill

```bash
curl -X POST http://localhost:3000/api/resume/skills \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Go",
    "proficiency": 85,
    "category": "Programming Languages"
  }'
```

#### Add Work Experience

```bash
curl -X POST http://localhost:3000/api/resume/experience \
  -H "Content-Type: application/json" \
  -d '{
    "job_title": "Backend Developer",
    "company": "Tech Solutions Inc.",
    "location": "Remote",
    "start_date": "2021-01-01T00:00:00Z",
    "current_job": true,
    "description": "Developing RESTful APIs using Go and microservices architecture"
  }'
```

#### Add Education

```bash
curl -X POST http://localhost:3000/api/resume/education \
  -H "Content-Type: application/json" \
  -d '{
    "institution": "Tech University",
    "degree": "Bachelor of Science",
    "field": "Computer Science",
    "location": "San Francisco, CA",
    "start_date": "2016-09-01T00:00:00Z",
    "end_date": "2020-05-30T00:00:00Z"
  }'
```

These commands will help you quickly test the functionality of the Resume API. For protected endpoints, you'll need to include an authorization token, which you can obtain by using the login endpoint.

## Application Workflow

### System Architecture

The application follows a layered architecture approach with clear separation of concerns:

1. **Presentation Layer (API)**: The HTTP handlers in `cmd/api/main.go` handle the incoming HTTP requests, process them, and return appropriate HTTP responses.

2. **Controller Layer**: The controllers in `internal/controllers/` define the business logic processing for each API endpoint, handling validation, error handling, and response formatting.

3. **Model Layer**: The models in `internal/models/` define the data structures and database schema. These models represent the business entities in the system.

4. **Database Layer**: The database connections and migrations in `internal/database/` manage data persistence and retrieval.

This layered approach ensures high maintainability, testability, and separation of concerns.

### API Request Lifecycle

1. **Request Routing**: When a client sends an HTTP request, it first goes through the router in `cmd/api/main.go`, which determines which handler should process the request based on the URL path.

2. **Authentication & Authorization**: For protected endpoints, the request is authenticated using JWT tokens to verify the user's identity and permissions.

3. **Request Processing**: The corresponding controller method processes the request, performs any required business logic, interacts with the database, and prepares the response.

4. **Response Generation**: The API formats the data as JSON and sends it back to the client with appropriate HTTP status codes and headers.

### Data Flow

1. **Client Request**: A client (web app, mobile app, etc.) sends an HTTP request to the API.
2. **Route Matching**: The router matches the URL path to the appropriate handler.
3. **Controller Processing**: The controller processes the request, validates inputs, and interacts with models.
4. **Database Operations**: The models interact with the database to perform CRUD operations.
5. **Response Generation**: The controller formats the results and sends a response back to the client.

## API Structure

The API follows RESTful conventions and is organized around resources:

### Authentication System

The authentication system uses JWT (JSON Web Tokens) for secure user authentication:

- **Registration**: Users can create an account with username/password
- **Login**: Users can authenticate and receive a JWT token
- **Authorization**: Protected endpoints verify the JWT token to ensure the user has appropriate permissions

### Project Management

The Projects API allows creating and retrieving projects with:

- **Categories**: Projects can be organized by categories
- **Tags**: Projects can be tagged for better searchability
- **Media Attachments**: Projects can have images and other media files attached

### Blog System

The Blog API provides endpoints for managing blog posts:

- **Posts**: Full CRUD operations for blog posts
- **Categories**: Blog posts can be organized by categories
- **Media Support**: Blog posts can include images and other media
- **SEO**: Blog posts contain metadata for SEO optimization

### Resume System

The Resume API manages personal/professional information with multiple sections:

- **Personal Information**: Basic contact and biographical details
- **Skills**: Technical and soft skills with proficiency levels
- **Work Experience**: Employment history with descriptions
- **Education**: Academic background and qualifications
- **Certificates**: Professional certifications
- **Languages**: Language proficiency levels
- **Publications**: Academic or professional publications
- **Complete Resume**: Comprehensive view of all resume sections

## API Endpoints

### General Endpoints

| Method | Route                     | Description              | Access |
| ------ | ------------------------- | ------------------------ | ------ |
| GET    | /                         | API status               | Public |
| GET    | /health                   | Health check             | Public |
| GET    | /api                      | API welcome              | Public |

### Authentication Endpoints

| Method | Route                     | Description              | Access |
| ------ | ------------------------- | ------------------------ | ------ |
| POST   | /api/auth/login           | Login via phone/password | Public |
| POST   | /api/auth/register        | Register new user        | Public |

### Project Endpoints

| Method | Route                     | Description              | Access |
| ------ | ------------------------- | ------------------------ | ------ |
| GET    | /api/projects             | Get list of projects     | Public |
| POST   | /api/projects             | Create project           | Admin  |

### Blog Endpoints

| Method | Route                     | Description              | Access |
| ------ | ------------------------- | ------------------------ | ------ |
| GET    | /api/blog                 | Get blog posts           | Public |
| POST   | /api/blog                 | Create blog post         | Admin  |
| GET    | /api/categories/projects  | Get project categories   | Public |
| GET    | /api/categories/blog      | Get blog categories      | Public |

### Resume Endpoints

| Method | Route                         | Description                  | Access |
| ------ | ----------------------------- | ---------------------------- | ------ |
| GET    | /api/resume/personal          | Get personal information     | Public |
| POST   | /api/resume/personal          | Create personal information  | Admin  |
| PUT    | /api/resume/personal/:id      | Update personal information  | Admin  |
| DELETE | /api/resume/personal/:id      | Delete personal information  | Admin  |
| GET    | /api/resume/skills            | Get skills                   | Public |
| POST   | /api/resume/skills            | Create skill                 | Admin  |
| PUT    | /api/resume/skills/:id        | Update skill                 | Admin  |
| DELETE | /api/resume/skills/:id        | Delete skill                 | Admin  |
| GET    | /api/resume/experience        | Get work experience          | Public |
| POST   | /api/resume/experience        | Create work experience       | Admin  |
| PUT    | /api/resume/experience/:id    | Update work experience       | Admin  |
| DELETE | /api/resume/experience/:id    | Delete work experience       | Admin  |
| GET    | /api/resume/education         | Get education details        | Public |
| POST   | /api/resume/education         | Create education detail      | Admin  |
| PUT    | /api/resume/education/:id     | Update education detail      | Admin  |
| DELETE | /api/resume/education/:id     | Delete education detail      | Admin  |
| GET    | /api/resume/certificates      | Get certificates             | Public |
| POST   | /api/resume/certificates      | Create certificate           | Admin  |
| PUT    | /api/resume/certificates/:id  | Update certificate           | Admin  |
| DELETE | /api/resume/certificates/:id  | Delete certificate           | Admin  |
| GET    | /api/resume/languages         | Get languages                | Public |
| POST   | /api/resume/languages         | Create language              | Admin  |
| PUT    | /api/resume/languages/:id     | Update language              | Admin  |
| DELETE | /api/resume/languages/:id     | Delete language              | Admin  |
| GET    | /api/resume/publications      | Get publications             | Public |
| POST   | /api/resume/publications      | Create publication           | Admin  |
| PUT    | /api/resume/publications/:id  | Update publication           | Admin  |
| DELETE | /api/resume/publications/:id  | Delete publication           | Admin  |
| GET    | /api/resume/complete          | Get complete resume          | Public |

## Resume Data Models

The Resume API uses the following data models:

### Personal Information

```json
{
  "id": 1,
  "full_name": "John Doe",
  "job_title": "Full Stack Developer",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "address": "123 Main St, City, Country",
  "website": "https://johndoe.com",
  "linkedin": "https://linkedin.com/in/johndoe",
  "github": "https://github.com/johndoe",
  "twitter": "https://twitter.com/johndoe",
  "summary": "Experienced developer with 5+ years in web development",
  "profile_image": "https://example.com/profile.jpg",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Skill

```json
{
  "id": 1,
  "name": "JavaScript",
  "proficiency": 90,
  "category": "Programming Languages",
  "icon_url": "https://example.com/js-icon.png",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Experience

```json
{
  "id": 1,
  "job_title": "Senior Developer",
  "company": "Tech Company",
  "location": "San Francisco, CA",
  "start_date": "2020-01-01T00:00:00Z",
  "end_date": "2022-12-31T00:00:00Z",
  "current_job": false,
  "description": "Led development of web applications using React and Node.js",
  "achievements": "Improved system performance by 40%, Mentored junior developers",
  "website": "https://techcompany.com",
  "logo_url": "https://example.com/company-logo.png",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Education

```json
{
  "id": 1,
  "institution": "University of Technology",
  "degree": "Bachelor of Science",
  "field": "Computer Science",
  "location": "Boston, MA",
  "start_date": "2016-09-01T00:00:00Z",
  "end_date": "2020-06-30T00:00:00Z",
  "current": false,
  "gpa": "3.8/4.0",
  "description": "Focus on software engineering and data structures",
  "logo_url": "https://example.com/university-logo.png",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Certificate

```json
{
  "id": 1,
  "name": "AWS Certified Developer",
  "issuer": "Amazon Web Services",
  "issue_date": "2022-01-15T00:00:00Z",
  "expiry_date": "2025-01-15T00:00:00Z",
  "no_expiry": false,
  "credential_id": "AWS-DEV-123456",
  "credential_url": "https://aws.amazon.com/verification/12345",
  "description": "Certification for cloud development on AWS platform",
  "logo_url": "https://example.com/aws-logo.png",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Language

```json
{
  "id": 1,
  "name": "English",
  "proficiency": "Native",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Publication

```json
{
  "id": 1,
  "title": "Modern Web Development Techniques",
  "publisher": "Tech Journal",
  "authors": "John Doe, Jane Smith",
  "publish_date": "2022-03-10T00:00:00Z",
  "url": "https://techjournal.com/articles/123",
  "doi": "10.1234/tj.2022.123456",
  "description": "Research on efficient web development methodologies",
  "image_url": "https://example.com/publication-cover.jpg",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Complete Resume

The complete resume endpoint returns an object containing all resume sections:

```json
{
  "personal_info": [/* PersonalInfo objects */],
  "skills": [/* Skill objects */],
  "experience": [/* Experience objects */],
  "education": [/* Education objects */],
  "projects": [/* Project objects */],
  "certificates": [/* Certificate objects */],
  "languages": [/* Language objects */],
  "publications": [/* Publication objects */]
}
```

## API Response Format

All API responses follow a consistent JSON format:

### Success Response

```json
{
  "message": "Success message",
  "data": {
    // Response data
  }
}
```

### Error Response

```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```

## HTTP Status Codes

The API uses standard HTTP status codes:

- **200 OK**: The request succeeded
- **201 Created**: Resource created successfully
- **400 Bad Request**: Invalid request parameters
- **401 Unauthorized**: Authentication required
- **403 Forbidden**: Insufficient permissions
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Server-side error

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

3. Create a `.env` file with your database configuration
   ```bash
   nano .env
   ```
   
   Add the following content (adjust values for your environment):
   ```
   # Application settings
   APP_ENV=production
   APP_PORT=3000
   APP_HOST=0.0.0.0
   APP_NAME=zione-backend
   APP_URL=https://yourapi.example.com
   APP_SECRET=your_secure_secret_key
   
   # Database settings
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=your_production_db
   DB_USER=your_production_user
   DB_PASSWORD=your_secure_password
   DB_CHARSET=utf8mb4
   DB_MAX_IDLE_CONNS=10
   DB_MAX_OPEN_CONNS=100
   DB_CONN_MAX_LIFETIME=3600
   
   # Other settings...
   # (Copy from the Setup section above and adjust as needed)
   ```

4. Make all scripts executable
   ```bash
   chmod +x *.sh
   ```

5. Build and run the application
   ```bash
   # Build the application first
   ./build.sh
   
   # Start the application (will use the binary if available, otherwise will use go run)
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

## Request Logging System

The API includes a comprehensive request logging system that monitors and records all incoming HTTP requests. This helps with debugging, monitoring, and security auditing.

### Logging Features

- **Request Details**: HTTP method, path, query parameters, and body content
- **Response Information**: Status code and success/error indication
- **Performance Metrics**: Request processing time (latency)
- **Client Data**: IP address and user agent
- **Error Tracking**: Any errors that occurred during request processing

### Log Storage

Logs are stored in two locations:
1. **Console Output**: All requests are logged to the console in real-time
2. **Daily Log Files**: Logs are saved to files in the `logs/` directory, with one file per day (format: `YYYY-MM-DD.log`)

### Log Format

Each log entry uses the following format:
```
[REQUEST] YYYY/MM/DD - HH:MM:SS | Status | METHOD PATH | StatusCode | Latency | ClientIP | Parameters | User-Agent | Errors
```

Example:
```
[REQUEST] 2023/05/15 - 14:32:45 | Success | GET /api/resume/skills | 200 | 15.2ms | 192.168.1.5 | ?category=frontend | User-Agent: Mozilla/5.0... | 
```

For POST and PUT requests, the request body is included in the log (limited to 1KB to prevent log file bloat).

### Accessing Logs

To view the latest logs:
```bash
tail -f logs/$(date +%Y-%m-%d).log
```

To search logs for specific requests:
```bash
grep "GET /api/resume" logs/2023-05-15.log
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -m 'Add feature'`
4. Push to the branch: `git push origin feature-name`
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 