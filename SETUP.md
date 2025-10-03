# Quick Setup Guide

## For Development

1. **Copy environment template:**
   ```bash
   cp env.example .env
   ```

2. **Generate JWT secret:**
   ```bash
   openssl rand -base64 32
   ```
   Add the generated string to your `.env` file as `JWT_SECRET`

3. **Configure database in `.env`:**
   ```env
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=your_database_name
   ```

4. **Run migrations:**
   ```bash
   migrate -path ./migrations -database "postgresql://username:password@localhost:5432/database_name?sslmode=disable" up
   ```

5. **Start the server:**
   ```bash
   go run cmd/main.go
   ```

## For Production Deployment

1. **Set deployment environment variables:**
   ```bash
   export DEPLOY_SERVER="user@your-server-ip"
   export DEPLOY_KEY_PATH="$HOME/.ssh/your_key"
   ```

2. **Create production environment file:**
   ```bash
   cp env.example .env.prod
   # Edit .env.prod with production values
   ```

3. **Deploy:**
   ```bash
   ./deploy/deploy.sh
   ```

## Environment Variable Reference

### Required Variables
- `JWT_SECRET` - JWT signing key (generate with `openssl rand -base64 32`)
- `DB_USER` - Database username (if not using `DATABASE_URL`)
- `DB_PASSWORD` - Database password (if not using `DATABASE_URL`)
- `DB_NAME` - Database name (if not using `DATABASE_URL`)

### Optional Variables
- `ENVIRONMENT` - App environment (default: `development`)
- `PORT` - Server port (default: `8080`)
- `GIN_MODE` - Gin framework mode (auto-detected based on ENVIRONMENT)
- `TOKEN_LIFETIME` - JWT token lifetime (default: `24h`)
- `ALLOWED_ORIGINS` - CORS allowed origins (comma-separated, default: allow all)
- `STATIC_PATH` - Static files directory (default: `./uploads`)
- `STATIC_URL` - Static files URL prefix (default: `/static`)
- `UPLOADS_PATH` - File uploads directory (default: `./uploads`)
- `API_PREFIX` - API routes prefix (default: `/api`)
- `DATABASE_URL` - Full database connection string (alternative to individual DB_* vars)
- `DB_HOST` - Database host (default: `localhost`)
- `DB_PORT` - Database port (default: `5432`)
- `DB_SSLMODE` - SSL mode (default: `disable`)
- `DB_MAX_OPEN_CONNS` - Max open database connections (default: `25`)
- `DB_MAX_IDLE_CONNS` - Max idle database connections (default: `5`)

## Security Checklist

- [ ] Strong JWT secret generated
- [ ] Database credentials secured
- [ ] `.env` files not committed to git
- [ ] Production environment configured
- [ ] SSL enabled for production database
- [ ] Server access properly secured
