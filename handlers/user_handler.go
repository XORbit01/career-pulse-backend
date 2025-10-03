package handlers

import (
	"database/sql"
	"net/http"

	"github.com/XORbit01/jobseeker-backend/models"
	"github.com/XORbit01/jobseeker-backend/repos"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo *repos.UserRepository
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		userRepo: repos.NewUserRepository(db),
	}
}

func RegisterUserRoutes(router *gin.RouterGroup, db *sql.DB) {
	handler := NewUserHandler(db)

	router.GET("/me", handler.GetCurrentUser)
	router.DELETE("/me", handler.DeleteCurrentUser)
}

//	@Summary		Get current user
//	@Description	Returns the currently authenticated user's `profile`
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.SuccessResponse{data=models.User}
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Router			/users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   &models.ErrorInfo{Code: "UNAUTHORIZED"},
		})
		return
	}

	user, err := h.userRepo.GetByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "User not found",
			Error:   &models.ErrorInfo{Code: "USER_NOT_FOUND"},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

//	@Summary		Delete current user
//	@Description	Deletes the currently authenticated user account
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.SuccessResponse{data=nil}
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users/me [delete]
func (h *UserHandler) DeleteCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   &models.ErrorInfo{Code: "UNAUTHORIZED"},
		})
		return
	}

	if err := h.userRepo.Delete(userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to delete user",
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}
