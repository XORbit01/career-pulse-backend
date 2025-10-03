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

// JobHandler handles job-related routes
type JobHandler struct {
	jobRepo      *repos.JobRepository
	employerRepo *repos.EmployerRepository
}

// NewJobHandler creates a new JobHandler
func NewJobHandler(db *sql.DB) *JobHandler {
	return &JobHandler{
		jobRepo:      repos.NewJobRepository(db),
		employerRepo: repos.NewEmployerRepository(db),
	}
}

// RegisterJobRoutesPrivate registers job-related routes
func RegisterJobRoutesPrivate(router *gin.RouterGroup, db *sql.DB) {
	handler := NewJobHandler(db)
	// Employer-only routes
	employerGroup := router.Group("/")
	employerGroup.GET("/employer/listings", handler.GetEmployerJobs)
	employerGroup.Use(middleware.RoleMiddleware("employer"))
	{
		employerGroup.POST("", handler.CreateJob)
		employerGroup.PUT("/:id", handler.UpdateJob)
		employerGroup.DELETE("/:id", handler.DeleteJob)
	}
}

// RegisterJobRoutes registers public job-related routes
func RegisterJobRoutes(router *gin.RouterGroup, db *sql.DB) {
	handler := NewJobHandler(db)

	// Public routes
	router.GET("", handler.SearchJobs)
	router.GET("/:id", handler.GetJob)
}

// @Summary		List jobs by the current employer
// @Description	Returns a paginated list of jobs created by the authenticated employer
// @Tags			Jobs
// @Security		BearerAuth
// @Produce		json
// @Param			page	query		int	false	"Page number"
// @Param			limit	query		int	false	"Results per page (max 100)"
// @Success		200		{object}	models.PaginatedResponse{data=[]models.Job}
// @Failure		401		{object}	models.ErrorResponse
// @Failure		500		{object}	models.ErrorResponse
// @Router			/jobs/employer/listings [get]
func (h *JobHandler) GetEmployerJobs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   &models.ErrorInfo{Code: "UNAUTHORIZED"},
		})
		return
	}

	employer, err := h.employerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Employer profile not found",
			Error:   &models.ErrorInfo{Code: "PROFILE_NOT_FOUND"},
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	jobs, total, err := h.jobRepo.GetByEmployerID(employer.ID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve jobs",
			Error:   &models.ErrorInfo{Code: "QUERY_FAILED", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Message:    "Employer jobs retrieved successfully",
		Data:       jobs,
		Page:       page,
		TotalPages: (total + limit - 1) / limit,
		TotalItems: total,
		Limit:      limit,
	})
}

// CreateJob godoc
//
//	@Summary		Create a new job posting
//	@Description	Employers can create a new job posting
//	@Tags			Jobs
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.JobInput	true	"Job input"
//	@Success		201		{object}	models.SuccessResponse{data=models.Job}
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/jobs [post]
func (h *JobHandler) CreateJob(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "unauthorized"})
		return
	}

	employer, err := h.employerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "employer profile not found, please create one first"})
		return
	}

	var input models.JobInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	jobID, err := h.jobRepo.Create(employer.ID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	job, err := h.jobRepo.GetByID(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Message: "Job created successfully",
		Data:    job,
	})
}

// GetJob godoc
//
//	@Summary		Get a job by ID
//	@Description	Retrieve a job posting by its ID
//	@Tags			Jobs
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Job ID"
//	@Success		200	{object}	models.SuccessResponse{data=models.Job}
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Router			/jobs/{id} [get]
func (h *JobHandler) GetJob(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid job ID"})
		return
	}

	job, err := h.jobRepo.GetByID(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "job not found"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Job retrieved successfully",
		Data:    job,
	})
}

// UpdateJob godoc
//
//	@Summary		Update a job posting
//	@Description	Employers can update their job postings
//	@Tags			Jobs
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Job ID"
//	@Param			input	body		models.JobInput	true	"Updated job input"
//	@Success		200		{object}	models.SuccessResponse{data=models.Job}
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		403		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/jobs/{id} [put]
func (h *JobHandler) UpdateJob(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "unauthorized"})
		return
	}

	jobID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid job ID"})
		return
	}

	employer, err := h.employerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "employer profile not found"})
		return
	}

	isOwner, err := h.jobRepo.IsJobOwner(jobID, employer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "you don't have permission to update this job"})
		return
	}

	var input models.JobInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	err = h.jobRepo.Update(jobID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	job, err := h.jobRepo.GetByID(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Job updated successfully",
		Data:    job,
	})
}

// DeleteJob godoc
//
//	@Summary		Delete a job posting
//	@Description	Employers can delete their job postings
//	@Tags			Jobs
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Job ID"
//	@Success		200	{object}	models.SuccessResponse
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		403	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/jobs/{id} [delete]
func (h *JobHandler) DeleteJob(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   &models.ErrorInfo{Code: "UNAUTHORIZED"},
		})
		return
	}

	jobID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid job ID",
			Error:   &models.ErrorInfo{Code: "INVALID_ID"},
		})
		return
	}

	employer, err := h.employerRepo.GetByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Employer profile not found",
			Error:   &models.ErrorInfo{Code: "PROFILE_NOT_FOUND"},
		})
		return
	}

	isOwner, err := h.jobRepo.IsJobOwner(jobID, employer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Error checking ownership",
			Error:   &models.ErrorInfo{Code: "OWNERSHIP_ERROR", Details: err.Error()},
		})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Success: false,
			Message: "Access denied: Not your job",
			Error:   &models.ErrorInfo{Code: "FORBIDDEN"},
		})
		return
	}

	err = h.jobRepo.Delete(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Error deleting job",
			Error:   &models.ErrorInfo{Code: "DELETE_ERROR", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Job deleted successfully",
	})
}

// SearchJobs godoc
//
//		@Summary		Search for jobs
//		@Description	Search public job listings using filters and pagination
//		@Tags			Jobs
//		@Accept			json
//		@Produce		json
//		@Param			title				query		string		false	"Job title"
//		@Param			location			query		string		false	"Location"
//		@Param			job_type			query		string		false	"Job type"	Enums(full_time, part_time, contract, internship, remote)
//		@Param			min_salary			query		number		false	"Minimum salary"
//		@Param			experience_level	query		string		false	"Experience level"	Enums(Entry-level, Mid-level, Senior, Lead)
//		@Param			skills				query		[]string	false	"Comma-separated skill list"
//		@Param			employer_user_id	query		int			false	"Employer user ID to filter jobs by company"
//		@Param			page				query		int			false	"Page number"
//		@Param			limit				query		int			false	"Results per page (max 100)"
//	    @Param          category			query       string      false    "category or industry"
//		@Success		200					{object}	models.PaginatedResponse{data=[]models.Job}
//		@Failure		400					{object}	models.ErrorResponse
//		@Failure		500					{object}	models.ErrorResponse
//		@Router			/jobs [get]
func (h *JobHandler) SearchJobs(c *gin.Context) {
	var params models.JobSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid search parameters",
			Error:   &models.ErrorInfo{Code: "INVALID_PARAMS", Details: err.Error()},
		})
		return
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	params.Skills = c.QueryArray("skills")

	jobs, total, err := h.jobRepo.SearchJobs(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to search jobs",
			Error:   &models.ErrorInfo{Code: "SEARCH_FAILED", Details: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Message:    "Jobs retrieved successfully",
		Data:       jobs,
		Page:       params.Page,
		TotalPages: (total + params.Limit - 1) / params.Limit,
		TotalItems: total,
		Limit:      params.Limit,
	})
}
