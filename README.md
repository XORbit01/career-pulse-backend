# Career Pulse Backend

A Go-based REST API for the Career Pulse job seeking platform, featuring user authentication, job management, and real-time chat functionality.

## Features

- **User Authentication**: JWT-based authentication with role-based access control
- **Job Management**: CRUD operations for job postings and applications
- **Real-time Chat**: WebSocket-based messaging between employers and job seekers
- **File Upload**: Support for profile images and document uploads
- **API Documentation**: Auto-generated Swagger documentation

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## Setup

### 1. Clone the Repository

```bash
git clone <your-repo-url>
cd career-pulse-backend
```

### 2. Environment Configuration

Copy the environment template and configure your settings:

```bash
cp env.example .env
```

Edit `.env` with your actual values:

```env
# Application Environment
ENVIRONMENT=development

# Server Port
PORT=8080

# JWT Configuration (REQUIRED)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Token Lifetime
TOKEN_LIFETIME=24h

# Database Configuration
# Option 1: Use full DATABASE_URL (recommended)
DATABASE_URL=postgresql://username:password@localhost:5432/database_name?sslmode=disable

# Option 2: Use individual parameters
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_database_user
DB_PASSWORD=your_database_password
DB_NAME=your_database_name
DB_SSLMODE=disable
```

**Important**: Generate a secure JWT secret:
```bash
openssl rand -base64 32
```

### 3. Database Setup

Create a PostgreSQL database and run migrations:

```bash
# Install migration tool (if not already installed)
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path ./migrations -database "postgresql://username:password@localhost:5432/database_name?sslmode=disable" up
```

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Application

```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8080`

## API Documentation

Once the server is running, visit:
- Swagger UI: `http://localhost:8080/docs/`
- API Base: `http://localhost:8080/api`

## Deployment

### Using the Deployment Script

The project includes a deployment script that requires environment variables:

```bash
# Set deployment configuration
export DEPLOY_SERVER="user@your-server-ip"
export DEPLOY_KEY_PATH="$HOME/.ssh/your_key"
export REMOTE_DIR="/path/to/remote/directory"  # optional
export BINARY_NAME="jobseeker"                # optional

# Run deployment
./deploy/deploy.sh
```

### Manual Deployment

1. Build the binary:
```bash
GOOS=linux GOARCH=amd64 go build -o jobseeker ./cmd
```

2. Copy to your server:
```bash
scp jobseeker user@server:/path/to/app/
scp .env.prod user@server:/path/to/app/.env
```

3. Set up systemd service (see `deploy/jobseeker.service`)

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `ENVIRONMENT` | Application environment | No | `development` |
| `PORT` | Server port | No | `8080` |
| `GIN_MODE` | Gin framework mode | No | Auto-detected |
| `JWT_SECRET` | JWT signing secret | **Yes** | - |
| `TOKEN_LIFETIME` | JWT token lifetime | No | `24h` |
| `ALLOWED_ORIGINS` | CORS allowed origins (comma-separated) | No | Allow all |
| `STATIC_PATH` | Static files directory path | No | `./uploads` |
| `STATIC_URL` | Static files URL prefix | No | `/static` |
| `UPLOADS_PATH` | File uploads directory | No | `./uploads` |
| `API_PREFIX` | API routes prefix | No | `/api` |
| `DATABASE_URL` | Full database connection string | No | - |
| `DB_HOST` | Database host | No* | `localhost` |
| `DB_PORT` | Database port | No* | `5432` |
| `DB_USER` | Database username | No* | - |
| `DB_PASSWORD` | Database password | No* | - |
| `DB_NAME` | Database name | No* | - |
| `DB_SSLMODE` | SSL mode | No | `disable` |
| `DB_MAX_OPEN_CONNS` | Max open database connections | No | `25` |
| `DB_MAX_IDLE_CONNS` | Max idle database connections | No | `5` |
| `ENV_FILE` | Environment file path | No | `.env` |

*Required if `DATABASE_URL` is not provided

## Project Structure

```
career-pulse-backend/
├── cmd/                    # Application entry point
├── config/                 # Configuration management
├── db/                     # Database connection
├── handlers/               # HTTP request handlers
├── middleware/             # Custom middleware
├── models/                 # Data models
├── repos/                  # Repository layer
├── migrations/             # Database migrations
├── docs/                   # Auto-generated API docs
├── deploy/                 # Deployment scripts
└── uploads/                # File upload directory
```

## Security Considerations

- Never commit `.env` files or other sensitive configuration
- Use strong, unique JWT secrets in production
- Enable SSL for database connections in production
- Regularly rotate JWT secrets
- Use environment-specific configurations

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
