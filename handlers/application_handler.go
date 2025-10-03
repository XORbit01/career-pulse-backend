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

// ApplicationHandler handles application-related routes
type ApplicationHandler struct {
	applicationRepo *repos.ApplicationRepository
	jobSeekerRepo   *repos.JobSeekerRepository
	employerRepo    *repos.EmployerRepository
	jobRepo         *repos.JobRepository
}

// NewApplicationHandler creates a new ApplicationHandler
func NewApplicationHandler(db *sql.DB) *ApplicationHandler {
	return &ApplicationHandler{
		applicationRepo: repos.NewApplicationRepository(db),
		jobSeekerRepo:   repos.NewJobSeekerRepository(db),
		employerRepo:    repos.NewEmployerRepository(db),
		jobRepo:         repos.NewJobRepository(db),
	}
}

// RegisterApplicationRoutes registers application routes
func RegisterApplicationRoutes(router *gin.RouterGroup, db *sql.DB) {
	handler := NewApplicationHandler(db)

	// Job seeker routes
	jobSeekerGroup := router.Group("/")
	jobSeekerGroup.Use(middleware.RoleMiddleware("job_seeker"))
	{
		jobSeekerGroup.POST("", handler.CreateApplication)
		jobSeekerGroup.GET("/job-seeker", handler.GetJobSeekerApplications)
		jobSeekerGroup.DELETE("/:id", handler.DeleteApplication)
	}

	// Employer routes
	employerGroup := router.Group("/")
	employerGroup.Use(middleware.RoleMiddleware("employer"))
	{
		employerGroup.GET("/job/:jobId", handler.GetJobApplications)
		employerGroup.PUT("/:id/status", handler.UpdateApplicationStatus)
	}

	// Common routes (accessible by both roles)
	router.GET("/:id", handler.GetApplication)
}

// CreateApplication godoc
//
//	@Summary		Apply for a job
//	@Description	Submit an application to a job. Requires role: job_seeker
//	@Tags			Applications
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.ApplicationInput	true	"Application input"
//	@Success		201		{object}	models.SuccessResponse{data=models.Application}
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		403		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/applications [post]
func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "unauthorized"})
		return
	}

	// Get job seeker profile
	jobSeeker, err := h.jobSeekerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "job seeker profile not found, please create one first"})
		return
	}

	var input models.ApplicationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	_, err = h.jobRepo.GetByID(input.JobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "job not found"})
		return
	}

	// Create application
	applicationID, err := h.applicationRepo.Create(jobSeeker.ID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	// Get created application
	application, err := h.applicationRepo.GetByID(applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, application)
}

// GetApplication godoc
//
//	@Summary		Get application by ID
//	@Description	Retrieve a job application. Requires role: job_seeker (only own) or employer (only own job's applications).
//	@Tags			Applications
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Application ID"
//	@Success		200	{object}	models.SuccessResponse{data=models.Application}
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		403	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/applications/{id} [get]
func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "unauthorized"})
		return
	}

	userRole, _ := c.Get("userRole")
	role := userRole.(string)

	applicationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid application ID"})
		return
	}

	// Check if user has access to this application
	switch role {
	case "job_seeker":
		jobSeeker, err := h.jobSeekerRepo.GetByUserID(userID.(int))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "job seeker profile not found"})
			return
		}

		isOwner, err := h.applicationRepo.IsJobSeekerApplication(applicationID, jobSeeker.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
			return
		}
		if !isOwner {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "you don't have permission to view this application"})
			return
		}
	case "employer":
		employer, err := h.employerRepo.GetByUserID(userID.(int))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "employer profile not found"})
			return
		}

		isOwner, err := h.applicationRepo.IsJobOwnerApplication(applicationID, employer.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
			return
		}
		if !isOwner {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "you don't have permission to view this application"})
			return
		}
	}

	// Get application
	application, err := h.applicationRepo.GetByID(applicationID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "application not found"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true, Data: application})
}

// GetJobSeekerApplications godoc
//
//	@Summary		Get applications of the current job seeker
//	@Description	Returns all job applications submitted by the current authenticated job seeker.
//	@Tags			Applications
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page	query		int	false	"Page number"		default(1)
//	@Param			limit	query		int	false	"Results per page"	default(10)
//	@Success		200		{object}	models.PaginatedResponse{data=[]models.Application}
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/applications/job-seeker [get]
func (h *ApplicationHandler) GetJobSeekerApplications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "unauthorized"})
		return
	}

	// Get job seeker profile
	jobSeeker, err := h.jobSeekerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "job seeker profile not found"})
		return
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Get applications
	applications, total, err := h.applicationRepo.GetByJobSeekerID(jobSeeker.ID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       applications,
		Page:       page,
		TotalPages: totalPages,
		TotalItems: total,
		Limit:      limit,
	})
}

// GetJobApplications godoc
//
//	@Summary		Get applications for a specific job
//	@Description	Returns all job applications submitted to a specific job owned by the current employer.
//	@Tags			Applications
//	@Security		BearerAuth
//	@Produce		json
//	@Param			jobId	path		int	true	"Job ID"
//	@Param			page	query		int	false	"Page number"		default(1)
//	@Param			limit	query		int	false	"Results per page"	default(10)
//	@Success		200		{object}	models.PaginatedResponse{data=[]models.Application}
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		403		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/applications/job/{jobId} [get]
func (h *ApplicationHandler) GetJobApplications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "unauthorized"})
		return
	}

	// Get employer
	employer, err := h.employerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "employer profile not found"})
		return
	}

	jobID, err := strconv.Atoi(c.Param("jobId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid job ID"})
		return
	}

	// Check if job belongs to this employer
	isOwner, err := h.jobRepo.IsJobOwner(jobID, employer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "you don't have permission to view applications for this job"})
		return
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Get applications
	applications, total, err := h.applicationRepo.GetByJobID(jobID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       applications,
		Page:       page,
		TotalPages: totalPages,
		TotalItems: total,
		Limit:      limit,
	})
}

// UpdateApplicationStatus godoc
//
//	@Summary		Update application status
//	@Description	Allows employers to update the status of applications for their own jobs.
//	@Tags			Applications
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Application ID"
//	@Param			input	body		models.ApplicationStatusInput	true	"New status"
//	@Success		200		{object}	models.SuccessResponse{data=models.Application}
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		403		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/applications/{id}/status [put]
func (h *ApplicationHandler) UpdateApplicationStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "unauthorized"})
		return
	}

	// Get employer profile
	employer, err := h.employerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "employer profile not found"})
		return
	}

	applicationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid application ID"})
		return
	}

	// Check if application is for a job owned by this employer
	isOwner, err := h.applicationRepo.IsJobOwnerApplication(applicationID, employer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "you don't have permission to update this application"})
		return
	}

	var input models.ApplicationStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	// Update application status
	err = h.applicationRepo.UpdateStatus(applicationID, input.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	// Get updated application
	application, err := h.applicationRepo.GetByID(applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, application)
}

// DeleteApplication godoc
//
//	@Summary		Delete an application
//	@Description	Allows job seekers to delete their own application. Requires role: job_seeker
//	@Tags			Applications
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Application ID"
//	@Success		200	{object}	models.SuccessResponse{data=nil}
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		403	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/applications/{id} [delete]
func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "unauthorized"})
		return
	}

	// Get job seeker profile
	jobSeeker, err := h.jobSeekerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "job seeker profile not found"})
		return
	}

	applicationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid application ID"})
		return
	}

	// Check if application belongs to this job seeker
	isOwner, err := h.applicationRepo.IsJobSeekerApplication(applicationID, jobSeeker.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "you don't have permission to delete this application"})
		return
	}

	// Delete application
	err = h.applicationRepo.Delete(applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true, Message: "application deleted successfully"})
}
