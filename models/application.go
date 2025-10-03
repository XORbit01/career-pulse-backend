package models

import "time"

// Application represents a job application
type Application struct {
	ID          int       `json:"id" example:"1"`
	JobID       int       `json:"job_id" example:"101"`
	JobSeekerID int       `json:"job_seeker_id" example:"55"`
	CoverLetter string    `json:"cover_letter" example:"I am very excited to apply for this role. I believe my experience and passion make me a strong fit."`
	Status      string    `json:"status" example:"pending"`
	CreatedAt   time.Time `json:"created_at" example:"2025-04-14T10:18:32Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-04-14T10:18:32Z"`

	// Additional fields for detailed responses
	JobTitle    string `json:"job_title,omitempty" example:"Backend Engineer"`
	CompanyName string `json:"company_name,omitempty" example:"Tech Innovations Inc."`
	FirstName   string `json:"first_name,omitempty" example:"Ali"`
	LastName    string `json:"last_name,omitempty" example:"Khalil"`
	LogoURL     string `json:"logo_url"`
	ResumeURL   string `json:"resume_url" example:"/uploads/resumes/ali_resume.pdf"`
}

// ApplicationInput represents the data needed to create an application
type ApplicationInput struct {
	JobID       int    `json:"job_id" binding:"required" example:"101"`
	CoverLetter string `json:"cover_letter" example:"I'm highly motivated to join your team. Here's why I think I'd be a great fit..."`
}

// ApplicationStatusInput represents the data needed to update an application status
type ApplicationStatusInput struct {
	Status string `json:"status" binding:"required,oneof=pending reviewed interview rejected accepted" example:"interview"`
}
