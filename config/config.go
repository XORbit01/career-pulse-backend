package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type DBConfig struct {
	DSN      string // check if connection string is set
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Config struct {
	Environment   string
	Port          string
	JWTSecret     string
	DB            DBConfig
	TokenLifetime string
	// Server configuration
	GinMode string
	// CORS configuration
	AllowedOrigins []string
	// Static files configuration
	StaticPath  string
	StaticURL   string
	UploadsPath string
	// API configuration
	APIPrefix string
	// Database connection pool
	MaxOpenConns int
	MaxIdleConns int
}

func Load() (*Config, error) {

	dsn := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	tokenLifetime := os.Getenv("TOKEN_LIFETIME")
	if tokenLifetime == "" {
		tokenLifetime = "24h"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}
	if dsn != "" {
		// Server configuration
		ginMode := os.Getenv("GIN_MODE")
		if ginMode == "" {
			if env == "production" {
				ginMode = "release"
			} else {
				ginMode = "debug"
			}
		}

		// CORS configuration
		allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
		var allowedOrigins []string
		if allowedOriginsStr != "" {
			allowedOrigins = strings.Split(allowedOriginsStr, ",")
			for i, origin := range allowedOrigins {
				allowedOrigins[i] = strings.TrimSpace(origin)
			}
		}

		// Static files configuration
		staticPath := os.Getenv("STATIC_PATH")
		if staticPath == "" {
			staticPath = "./uploads"
		}

		staticURL := os.Getenv("STATIC_URL")
		if staticURL == "" {
			staticURL = "/static"
		}

		uploadsPath := os.Getenv("UPLOADS_PATH")
		if uploadsPath == "" {
			uploadsPath = "./uploads"
		}

		// API configuration
		apiPrefix := os.Getenv("API_PREFIX")
		if apiPrefix == "" {
			apiPrefix = "/api"
		}

		// Database connection pool
		maxOpenConns := 25
		if maxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpenConnsStr != "" {
			if parsed, err := strconv.Atoi(maxOpenConnsStr); err == nil {
				maxOpenConns = parsed
			}
		}

		maxIdleConns := 5
		if maxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdleConnsStr != "" {
			if parsed, err := strconv.Atoi(maxIdleConnsStr); err == nil {
				maxIdleConns = parsed
			}
		}

		return &Config{
			Environment:    env,
			Port:           port,
			JWTSecret:      jwtSecret,
			TokenLifetime:  tokenLifetime,
			GinMode:        ginMode,
			AllowedOrigins: allowedOrigins,
			StaticPath:     staticPath,
			StaticURL:      staticURL,
			UploadsPath:    uploadsPath,
			APIPrefix:      apiPrefix,
			MaxOpenConns:   maxOpenConns,
			MaxIdleConns:   maxIdleConns,
			DB: DBConfig{
				DSN: dsn,
			},
		}, nil
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, errors.New("DB_USER environment variable is required")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, errors.New("DB_PASSWORD environment variable is required")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, errors.New("DB_NAME environment variable is required")
	}

	dbSSLMode := os.Getenv("DB_SSLMODE")
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	// Server configuration
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		if env == "production" {
			ginMode = "release"
		} else {
			ginMode = "debug"
		}
	}

	// CORS configuration
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsStr != "" {
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
		for i, origin := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(origin)
		}
	}

	// Static files configuration
	staticPath := os.Getenv("STATIC_PATH")
	if staticPath == "" {
		staticPath = "./uploads"
	}

	staticURL := os.Getenv("STATIC_URL")
	if staticURL == "" {
		staticURL = "/static"
	}

	uploadsPath := os.Getenv("UPLOADS_PATH")
	if uploadsPath == "" {
		uploadsPath = "./uploads"
	}

	// API configuration
	apiPrefix := os.Getenv("API_PREFIX")
	if apiPrefix == "" {
		apiPrefix = "/api"
	}

	// Database connection pool
	maxOpenConns := 25
	if maxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpenConnsStr != "" {
		if parsed, err := strconv.Atoi(maxOpenConnsStr); err == nil {
			maxOpenConns = parsed
		}
	}

	maxIdleConns := 5
	if maxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdleConnsStr != "" {
		if parsed, err := strconv.Atoi(maxIdleConnsStr); err == nil {
			maxIdleConns = parsed
		}
	}

	return &Config{
		Environment:    env,
		Port:           port,
		JWTSecret:      jwtSecret,
		TokenLifetime:  tokenLifetime,
		GinMode:        ginMode,
		AllowedOrigins: allowedOrigins,
		StaticPath:     staticPath,
		StaticURL:      staticURL,
		UploadsPath:    uploadsPath,
		APIPrefix:      apiPrefix,
		MaxOpenConns:   maxOpenConns,
		MaxIdleConns:   maxIdleConns,
		DB: DBConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			SSLMode:  dbSSLMode,
		},
	}, nil
}
