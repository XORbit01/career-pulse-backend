// handlers/profile_handler.go
package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/XORbit01/jobseeker-backend/models"
	"github.com/XORbit01/jobseeker-backend/repos"
	"github.com/gin-gonic/gin"
)

type PublicProfileHandler struct {
	employerRepo  *repos.EmployerRepository
	jobSeekerRepo *repos.JobSeekerRepository
}

func RegisterPublicProfileRoutes(router *gin.RouterGroup, db *sql.DB) {
	h := &PublicProfileHandler{
		employerRepo:  repos.NewEmployerRepository(db),
		jobSeekerRepo: repos.NewJobSeekerRepository(db),
	}

	router.GET("/:id", h.GetUnifiedPublicProfile)
}

// GetUnifiedPublicProfile godoc
//	@Summary		Get public profile by user ID
//	@Description	Returns either a job seeker or employer profile based on user ID
//	@Tags			Public
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	models.SuccessResponse{data=models.UnifiedPublicJobSeekerProfile}
//	@Success		200	{object}	models.SuccessResponse{data=models.UnifiedPublicEmployerProfile}
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Router			/profile/{id} [get]
func (h *PublicProfileHandler) GetUnifiedPublicProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   &models.ErrorInfo{Code: "INVALID_ID"},
		})
		return
	}

	if js, err := h.jobSeekerRepo.GetByUserID(userID); err == nil {
		response := models.UnifiedPublicJobSeekerProfile{
			UserID:          js.UserID,
			ProfileType:     "job_seeker",
			FirstName:       js.FirstName,
			LastName:        js.LastName,
			Headline:        js.Headline,
			Bio:             js.Summary,
			Location:        js.Location,
			ExperienceLevel: "Mid-level",
			Website:         "",
			ResumeURL:       js.ResumeURL,
			LogoURL:         js.LogoUrl,
			Skills:          []string{"JavaScript", "React", "Node.js", "Python"},
		}
		c.JSON(http.StatusOK, models.SuccessResponse{
			Success: true,
			Message: "Job seeker profile",
			Data:    response,
		})
		return
	}

	if emp, err := h.employerRepo.GetByUserID(userID); err == nil {
		response := models.UnifiedPublicEmployerProfile{
			UserID:      emp.UserID,
			ProfileType: "employer",
			CompanyName: emp.CompanyName,
			Description: emp.Description,
			Location:    emp.Location,
			Industry:    emp.Industry,
			CompanySize: emp.CompanySize,
			Website:     emp.Website,
			LogoURL:     emp.LogoURL,
		}
		c.JSON(http.StatusOK, models.SuccessResponse{
			Success: true,
			Message: "Employer profile",
			Data:    response,
		})
		return
	}

	c.JSON(http.StatusNotFound, models.ErrorResponse{
		Success: false,
		Message: "Profile not found",
		Error:   &models.ErrorInfo{Code: "PROFILE_NOT_FOUND"},
	})
}
