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

// JobSeekerHandler handles job seeker profile routes
type JobSeekerHandler struct {
	jobSeekerRepo *repos.JobSeekerRepository
	userRepo      *repos.UserRepository
}

// NewJobSeekerHandler creates a new JobSeekerHandler
func NewJobSeekerHandler(db *sql.DB) *JobSeekerHandler {
	return &JobSeekerHandler{
		jobSeekerRepo: repos.NewJobSeekerRepository(db),
		userRepo:      repos.NewUserRepository(db),
	}
}

// RegisterJobSeekerRoutes registers job seeker routes with role protection
func RegisterJobSeekerRoutes(router *gin.RouterGroup, db *sql.DB) {
	handler := NewJobSeekerHandler(db)

	jobSeekerGroup := router.Group("/")
	router.GET("/:id", handler.GetPublicProfileByID)
	jobSeekerGroup.Use(middleware.RoleMiddleware("job_seeker"))
	{
		jobSeekerGroup.POST("/profile", handler.CreateProfile)
		jobSeekerGroup.GET("/profile", handler.GetProfile)
		jobSeekerGroup.PUT("/profile", handler.UpdateProfile)
		jobSeekerGroup.DELETE("/profile", handler.DeleteProfile)
	}
}

// CreateProfile godoc
//
//	@Summary		Create job seeker profile
//	@Description	Creates a new job seeker profile for the authenticated user
//	@Tags			Job Seekers
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.JobSeekerProfileInput	true	"Profile input"
//	@Success		201		{object}	models.SuccessResponse{data=models.JobSeekerProfile}
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/job-seekers/profile [post]
func (h *JobSeekerHandler) CreateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized access",
			Error:   &models.ErrorInfo{Code: "UNAUTHORIZED"},
		})
		return
	}

	// Verify that the user actually exists in the database
	if _, err := h.userRepo.GetByID(userID.(int)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "User not found - please log in again",
			Error:   &models.ErrorInfo{Code: "USER_NOT_FOUND"},
		})
		return
	}

	var input models.JobSeekerProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid input: " + err.Error(),
			Error:   &models.ErrorInfo{Code: "INVALID_INPUT"},
		})
		return
	}

	if _, err := h.jobSeekerRepo.GetByUserID(userID.(int)); err == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Profile already exists",
			Error:   &models.ErrorInfo{Code: "PROFILE_EXISTS"},
		})
		return
	}

	profileID, err := h.jobSeekerRepo.Create(userID.(int), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to create profile",
			Error:   &models.ErrorInfo{Code: "CREATE_FAILED", Details: err.Error()},
		})
		return
	}

	profile, err := h.jobSeekerRepo.GetByID(profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch created profile",
			Error:   &models.ErrorInfo{Code: "FETCH_FAILED", Details: err.Error()},
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
//
//	@Summary		Get job seeker profile
//	@Description	Retrieves the authenticated job seeker's profile
//	@Tags			Job Seekers
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	models.SuccessResponse{data=models.JobSeekerProfile}
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Router			/job-seekers/profile [get]
func (h *JobSeekerHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized access",
			Error:   &models.ErrorInfo{Code: "UNAUTHORIZED"},
		})
		return
	}

	profile, err := h.jobSeekerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "Profile not found",
			Error:   &models.ErrorInfo{Code: "PROFILE_NOT_FOUND"},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    profile,
	})
}

// UpdateProfile godoc
//
//	@Summary		Update job seeker profile
//	@Description	Updates the authenticated user's job seeker profile
//	@Tags			Job Seekers
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.JobSeekerProfileInput	true	"Profile input"
//	@Success		200		{object}	models.SuccessResponse{data=models.JobSeekerProfile}
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/job-seekers/profile [put]
func (h *JobSeekerHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized access",
			Error:   &models.ErrorInfo{Code: "UNAUTHORIZED"},
		})
		return
	}

	var input models.JobSeekerProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid input: " + err.Error(),
			Error:   &models.ErrorInfo{Code: "INVALID_INPUT", Details: err.Error()},
		})
		return
	}

	profile, err := h.jobSeekerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "Profile not found",
			Error:   &models.ErrorInfo{Code: "PROFILE_NOT_FOUND", Details: err.Error()},
		})
		return
	}

	if err := h.jobSeekerRepo.Update(profile.ID, input); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to update profile",
			Error:   &models.ErrorInfo{Code: "UPDATE_FAILED", Details: err.Error()},
		})
		return
	}

	updatedProfile, err := h.jobSeekerRepo.GetByID(profile.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch updated profile",
			Error:   &models.ErrorInfo{Code: "FETCH_FAILED", Details: err.Error()},
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
//
//	@Summary		Delete job seeker profile
//	@Description	Deletes the authenticated user's job seeker profile
//	@Tags			Job Seekers
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	models.SuccessResponse{data=nil}
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/job-seekers/profile [delete]
func (h *JobSeekerHandler) DeleteProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized access",
			Error:   &models.ErrorInfo{Code: "UNAUTHORIZED"},
		})
		return
	}

	profile, err := h.jobSeekerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "Profile not found",
			Error:   &models.ErrorInfo{Code: "PROFILE_NOT_FOUND", Details: err.Error()},
		})
		return
	}

	if err := h.jobSeekerRepo.Delete(profile.ID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to delete profile",
			Error:   &models.ErrorInfo{Code: "DELETE_FAILED", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Profile deleted successfully",
	})
}

// GetPublicProfileByID godoc
//
//	@Summary		Get public job seeker profile
//	@Description	Retrieves a public view of a job seeker's profile by ID
//	@Tags			Job Seekers
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Job Seeker Profile ID"
//	@Success		200	{object}	models.SuccessResponse{data=models.JobSeekerProfile}
//	@Failure		404	{object}	models.ErrorResponse
//	@Router			/job-seekers/{id} [get]
func (h *JobSeekerHandler) GetPublicProfileByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "not valid input",
			Error:   &models.ErrorInfo{Code: "INVALID_INPUT", Details: err.Error()},
		})
		return
	}
	profile, err := h.jobSeekerRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "Profile not found",
			Error:   &models.ErrorInfo{Code: "PROFILE_NOT_FOUND", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Public profile retrieved",
		Data:    profile,
	})
}
