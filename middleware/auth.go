package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/XORbit01/jobseeker-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// TokenClaims represents the JWT token claims
type TokenClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware protects routes by requiring a valid JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		// 1. Try standard Authorization header
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Fallback to ?token query parameter (for WebSockets)
		if tokenStr == "" {
			tokenStr = c.Query("token")
		}

		// 3. If no token found
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Missing authentication token",
				Error:   &models.ErrorInfo{Code: "MISSING_TOKEN"},
			})
			c.Abort()
			return
		}

		// 4. Validate token
		claims, err := validateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Invalid or expired token",
				Error:   &models.ErrorInfo{Code: "INVALID_TOKEN", Details: err.Error()},
			})
			c.Abort()
			return
		}

		// 5. Save user ID and role in context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID int, role string, tokenLifetime string) (string, error) {
	duration, err := time.ParseDuration(tokenLifetime)
	if err != nil {
		duration = 24 * time.Hour
	}

	claims := TokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET environment variable is required")
	}

	return token.SignedString([]byte(jwtSecret))
}

// validateToken validates a JWT token and returns its claims
func validateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			return nil, errors.New("JWT_SECRET environment variable is required")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	return claims, nil
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "User role not found",
				Error:   &models.ErrorInfo{Code: "ROLE_MISSING"},
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		if !slices.Contains(allowedRoles, role) {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Success: false,
				Message: "Access denied: insufficient permissions",
				Error:   &models.ErrorInfo{Code: "FORBIDDEN_ROLE"},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
