package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/XORbit01/jobseeker-backend/middleware"
	"github.com/XORbit01/jobseeker-backend/models"
	"github.com/XORbit01/jobseeker-backend/repos"
	"github.com/gin-gonic/gin"
)

// EmployerHandler handles employer profile routes
type EmployerHandler struct {
	employerRepo *repos.EmployerRepository
}

// NewEmployerHandler creates a new EmployerHandler
func NewEmployerHandler(db *sql.DB) *EmployerHandler {
	return &EmployerHandler{
		employerRepo: repos.NewEmployerRepository(db),
	}
}

// RegisterEmployerRoutes registers employer profile endpoints
func RegisterEmployerRoutes(router *gin.RouterGroup, db *sql.DB) {
	handler := NewEmployerHandler(db)
	// public one
	employerGroup := router.Group("/")
	router.GET("/:id", handler.GetPublicEmployerProfile)
	employerGroup.Use(middleware.RoleMiddleware("employer"))
	{
		employerGroup.POST("/profile", handler.CreateProfile)
		employerGroup.GET("/profile", handler.GetProfile)
		employerGroup.PUT("/profile", handler.UpdateProfile)
		employerGroup.DELETE("/profile", handler.DeleteProfile)
	}
}

// CreateProfile godoc
//	@Summary	Create employer profile
//	@Tags		Employers
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		input	body		models.EmployerProfileInput	true	"Profile input"
//	@Success	201		{object}	models.SuccessResponse{data=models.EmployerProfile}
//	@Failure	400		{object}	models.ErrorResponse
//	@Failure	401		{object}	models.ErrorResponse
//	@Failure	500		{object}	models.ErrorResponse
//	@Router		/employers/profile [post]
func (h *EmployerHandler) CreateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   &models.ErrorInfo{Code: "UNAUTHORIZED"},
		})
		return
	}

	var input models.EmployerProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: err.Error(),
			Error:   &models.ErrorInfo{Code: "VALIDATION_ERROR"},
		})
		return
	}

	if _, err := h.employerRepo.GetByUserID(userID.(int)); err == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Profile already exists",
			Error:   &models.ErrorInfo{Code: "PROFILE_EXISTS"},
		})
		return
	}

	profileID, err := h.employerRepo.Create(userID.(int), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to create profile",
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Details: err.Error()},
		})
		return
	}

	profile, err := h.employerRepo.GetByID(profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch created profile",
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Message: "Profile created successfully",
		Data:    profile,
	})
}

// GetProfile godoc
//	@Summary	Get employer profile
//	@Tags		Employers
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{object}	models.SuccessResponse{data=models.EmployerProfile}
//	@Failure	401	{object}	models.ErrorResponse
//	@Failure	404	{object}	models.ErrorResponse
//	@Router		/employers/profile [get]
func (h *EmployerHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   &models.ErrorInfo{Code: "UNAUTHORIZED"},
		})
		return
	}

	profile, err := h.employerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "Profile not found",
			Error:   &models.ErrorInfo{Code: "NOT_FOUND"},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Profile fetched successfully",
		Data:    profile,
	})
}

// UpdateProfile godoc
//	@Summary	Update employer profile
//	@Tags		Employers
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		input	body		models.EmployerProfileInput	true	"Updated profile"
//	@Success	200		{object}	models.SuccessResponse{data=models.EmployerProfile}
//	@Failure	400		{object}	models.ErrorResponse
//	@Failure	401		{object}	models.ErrorResponse
//	@Failure	404		{object}	models.ErrorResponse
//	@Failure	500		{object}	models.ErrorResponse
//	@Router		/employers/profile [put]
func (h *EmployerHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   &models.ErrorInfo{Code: "UNAUTHORIZED"},
		})
		return
	}

	var input models.EmployerProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: err.Error(),
			Error:   &models.ErrorInfo{Code: "VALIDATION_ERROR"},
		})
		return
	}

	profile, err := h.employerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "Profile not found",
			Error:   &models.ErrorInfo{Code: "NOT_FOUND"},
		})
		return
	}

	if err := h.employerRepo.Update(profile.ID, input); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to update profile",
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Details: err.Error()},
		})
		return
	}

	updatedProfile, err := h.employerRepo.GetByID(profile.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch updated profile",
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    updatedProfile,
	})
}

// DeleteProfile godoc
//	@Summary	Delete employer profile
//	@Tags		Employers
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{object}	models.SuccessResponse{data=nil}
//	@Failure	401	{object}	models.ErrorResponse
//	@Failure	404	{object}	models.ErrorResponse
//	@Failure	500	{object}	models.ErrorResponse
//	@Router		/employers/profile [delete]
func (h *EmployerHandler) DeleteProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   &models.ErrorInfo{Code: "UNAUTHORIZED"},
		})
		return
	}

	profile, err := h.employerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "Profile not found",
			Error:   &models.ErrorInfo{Code: "NOT_FOUND"},
		})
		return
	}

	if err := h.employerRepo.Delete(profile.ID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to delete profile",
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Profile deleted successfully",
	})
}

// GetPublicEmployerProfile godoc
//	@Summary		Get public employer profile
//	@Description	Returns public-facing company info for the given employer ID
//	@Tags			Employers
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Employer Profile ID"
//	@Success		200	{object}	models.SuccessResponse{data=models.EmployerProfile}
//	@Failure		404	{object}	models.ErrorResponse
//	@Router			/employers/{id} [get]
func (h *EmployerHandler) GetPublicEmployerProfile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "not valid input",
			Error:   &models.ErrorInfo{Code: "INVALID_INPUT"},
		})
		return
	}
	profile, err := h.employerRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "Employer not found",
			Error:   &models.ErrorInfo{Code: "PROFILE_NOT_FOUND", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Public employer profile retrieved",
		Data:    profile,
	})
}
