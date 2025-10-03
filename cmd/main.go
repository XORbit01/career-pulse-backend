//	@title			Job Seeker API
//	@version		1.0
//	@description	## API documentation for Job Seeker platform
//	@description	This API follows a unified response **envelope pattern**. All responses will have `success`, `message`, and `data` or `error` keys. Example:
//	@description
//	@description	Success:
//	@description	```json
//	@description	{ "success": true, "message": "Fetched", "data": { ... } }
//	@description	```
//	@description
//	@description				Error:
//	@description				```json
//	@description				{ "success": false, "message": "Something went wrong", "error": { "code": "ERROR_CODE" } }
//	@description				```
//	@BasePath					/api
//	@schemes					http
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and your JWT token.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/XORbit01/jobseeker-backend/config"
	"github.com/XORbit01/jobseeker-backend/db"
	_ "github.com/XORbit01/jobseeker-backend/docs"
	"github.com/XORbit01/jobseeker-backend/handlers"
	"github.com/XORbit01/jobseeker-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerFiles "github.com/swaggo/files"
)

func main() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: %s file not found, using system environment variables\n", envFile)
	} else {
		log.Printf("✅ Loaded environment from %s\n", envFile)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	database, err := db.Connect(cfg.DB, cfg.MaxOpenConns, cfg.MaxIdleConns)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Set Gin mode based on configuration
	gin.SetMode(cfg.GinMode)

	router := gin.Default()
	router.RedirectFixedPath = false
	// static files
	router.Static(cfg.StaticURL, cfg.StaticPath)
	router.Use(middleware.CORS(cfg.AllowedOrigins))
	router.Use(middleware.ConfigMiddleware(cfg.StaticURL, cfg.UploadsPath))
	apiGroup := router.Group(cfg.APIPrefix)

	authGroup := apiGroup.Group("/auth") // /api/auth
	handlers.RegisterAuthRoutes(authGroup, database)

	protectedGroup := apiGroup.Group("/")
	protectedGroup.Use(middleware.AuthMiddleware())
	// image upload
	protectedGroup.POST("/upload", handlers.UploadFile)

	// User routes
	userGroup := protectedGroup.Group("/users")
	handlers.RegisterUserRoutes(userGroup, database)

	// Employer profile route
	employerGroup := protectedGroup.Group("/employers")
	handlers.RegisterEmployerRoutes(employerGroup, database)

	// job seekers
	jobSeekerGroup := protectedGroup.Group("/job-seekers")
	handlers.RegisterJobSeekerRoutes(jobSeekerGroup, database)

	// jobs public
	publicJobGroup := apiGroup.Group("/jobs")
	handlers.RegisterJobRoutes(publicJobGroup, database)

	// jobs private
	privateJobGroup := protectedGroup.Group("/jobs")
	handlers.RegisterJobRoutesPrivate(privateJobGroup, database)

	applicationGroup := protectedGroup.Group("/applications")
	handlers.RegisterApplicationRoutes(applicationGroup, database)
	// profile public
	publicProfileGroup := apiGroup.Group("/profile")
	handlers.RegisterPublicProfileRoutes(publicProfileGroup, database)

	// chat
	handlers.RegisterChatRoutes(protectedGroup, database)

	// swagger files
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
