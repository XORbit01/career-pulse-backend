package models

import "time"

// JobSeekerProfile represents a job seeker's profile
type JobSeekerProfile struct {
	ID              int       `json:"id" example:"1"`
	UserID          int       `json:"user_id" example:"6"`
	FirstName       string    `json:"first_name" example:"Ali"`
	LastName        string    `json:"last_name" example:"Khalil"`
	Headline        string    `json:"headline" example:"Junior Golang Backend Developer"`
	Summary         string    `json:"summary" example:"Passionate backend engineer with experience in RESTful APIs and microservices."`
	Phone           string    `json:"phone" example:"+96170123456"`
	Location        string    `json:"location" example:"Beirut, Lebanon"`
	ResumeURL       string    `json:"resume_url" example:"/uploads/resumes/ali_resume.pdf"`
	LogoUrl         string    `json:"logo_url" example:"/uploads/resumes/ali_pfp.jpeg"`
	Skills          []string  `json:"skills" example:["Go","PostgreSQL","Docker"]`
	ExperienceLevel string    `json:"experience_level" validate:"omitempty,oneof='Entry-level' 'Mid-level' 'Senior' 'Lead'" example:"Mid-level"`
	CreatedAt       time.Time `json:"created_at" example:"2025-04-14T10:18:32Z"`
	UpdatedAt       time.Time `json:"updated_at" example:"2025-04-14T10:18:32Z"`
}

// JobSeekerProfileInput represents the data needed to create/update a job seeker profile
type JobSeekerProfileInput struct {
	FirstName       string   `json:"first_name" binding:"required" example:"Ali"`
	LastName        string   `json:"last_name" binding:"required" example:"Khalil"`
	Headline        string   `json:"headline" example:"Junior Golang Backend Developer"`
	Summary         string   `json:"summary" example:"Experienced in Go, PostgreSQL, and Docker."`
	Phone           string   `json:"phone" example:"+96170123456"`
	Location        string   `json:"location" example:"Beirut, Lebanon"`
	ResumeURL       string   `json:"resume_url" example:"/uploads/resumes/ali_resume.pdf"`
	LogoUrl         string   `json:"logo_url" example:"/uploads/resumes/ali_pfp.jpeg"`
	Skills          []string `json:"skills" example:["Go","PostgreSQL","Docker"]`
	ExperienceLevel string   `json:"experience_level" validate:"omitempty,oneof='Entry-level' 'Mid-level' 'Senior' 'Lead'" example:"Mid-level"`
}
