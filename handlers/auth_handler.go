package handlers

import (
	"database/sql"
	"net/http"

	"github.com/XORbit01/jobseeker-backend/middleware"
	"github.com/XORbit01/jobseeker-backend/models"
	"github.com/XORbit01/jobseeker-backend/repos"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo      *repos.UserRepository
	tokenLifetime string
}

func NewAuthHandler(db *sql.DB, tokenLifetime string) *AuthHandler {
	return &AuthHandler{
		userRepo:      repos.NewUserRepository(db),
		tokenLifetime: tokenLifetime,
	}
}

func RegisterAuthRoutes(router *gin.RouterGroup, db *sql.DB) {
	tokenLifetime := "24h"
	handler := NewAuthHandler(db, tokenLifetime)

	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
}

//	@Summary		Register a new user
//	@Description	Create a new user account and return a JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.UserInput	true	"User input"
//	@Success		201		{object}	models.SuccessResponse{data=models.TokenResponse}
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input models.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid input",
			Error:   &models.ErrorInfo{Code: "INVALID_PAYLOAD", Details: err.Error()},
		})
		return
	}

	if _, err := h.userRepo.GetByEmail(input.Email); err == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Email already in use",
			Error:   &models.ErrorInfo{Code: "EMAIL_EXISTS"},
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Error hashing password",
			Error:   &models.ErrorInfo{Code: "HASH_ERROR"},
		})
		return
	}

	userID, err := h.userRepo.Create(input, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Details: err.Error()},
		})
		return
	}

	token, err := middleware.GenerateToken(userID, input.Role, h.tokenLifetime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Error generating token",
			Error:   &models.ErrorInfo{Code: "JWT_ERROR"},
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Message: "User registered successfully",
		Data: models.TokenResponse{
			Token: token,
			Role:  input.Role,
		},
	})
}

//	@Summary		Login an existing user
//	@Description	Authenticate user and return JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.LoginInput	true	"Login credentials"
//	@Success		200		{object}	models.SuccessResponse{data=models.TokenResponse}
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid input",
			Error:   &models.ErrorInfo{Code: "INVALID_PAYLOAD", Details: err.Error()},
		})
		return
	}

	user, err := h.userRepo.GetByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Invalid email or password",
			Error:   &models.ErrorInfo{Code: "AUTH_FAILED"},
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Invalid email or password",
			Error:   &models.ErrorInfo{Code: "AUTH_FAILED"},
		})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Role, h.tokenLifetime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Error generating token",
			Error:   &models.ErrorInfo{Code: "JWT_ERROR"},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Login successful",
		Data: models.TokenResponse{
			Token: token,
			Role:  user.Role,
		},
	})
}
