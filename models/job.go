package models

import (
	"time"
)

// Job represents a job posting
type Job struct {
	ID              int       `json:"id" example:"1"`
	EmployerID      int       `json:"employer_id" example:"42"`
	Title           string    `json:"title" example:"Senior Golang Developer"`
	Description     string    `json:"description" example:"We're looking for a backend engineer experienced in Go, PostgreSQL, and distributed systems."`
	Location        string    `json:"location" example:"Remote"`
	JobType         string    `json:"job_type" example:"full_time"`
	SalaryMin       *float64  `json:"salary_min" example:"60000"`
	SalaryMax       *float64  `json:"salary_max" example:"90000"`
	ExperienceLevel string    `json:"experience_level" validate:"omitempty,oneof='Entry-level' 'Mid-level' 'Senior' 'Lead'" example:"Mid-level"`
	RequiredSkills  []string  `json:"required_skills" example:["Go","PostgreSQL","Docker"]`
	Status          string    `json:"status" example:"active"`
	CreatedAt       time.Time `json:"created_at" example:"2025-04-14T10:18:32Z"`
	UpdatedAt       time.Time `json:"updated_at" example:"2025-04-14T10:18:32Z"`
	CompanyName     string    `json:"company_name,omitempty" example:"Tech Innovations Inc."`
	Category        string    `json:"category,omitempty" example:"Engineering"`
	LogoURL         string    `json:"logo_url,omitempty" example:"/uploads/logos/company123.png"`
}

// JobInput represents the data needed to create/update a job
type JobInput struct {
	Title           string   `json:"title" binding:"required" example:"Senior Golang Developer"`
	Description     string   `json:"description" binding:"required" example:"Work on scalable systems, microservices, and DevOps pipelines."`
	Location        string   `json:"location" example:"Remote"`
	JobType         string   `json:"job_type" binding:"required,oneof=full_time part_time contract internship remote" example:"full_time"`
	SalaryMin       *float64 `json:"salary_min" example:"60000"`
	SalaryMax       *float64 `json:"salary_max" example:"90000"`
	ExperienceLevel string   `json:"experience_level" validate:"omitempty,oneof='Entry-level' 'Mid-level' 'Senior' 'Lead'" example:"Mid-level"`
	RequiredSkills  []string `json:"required_skills" example:["Go","PostgreSQL","Docker"]`
	Category        string   `json:"category" binding:"required" example:"Engineering"`
	Status          string   `json:"status" binding:"omitempty,oneof=active closed draft" example:"active"`
}

// JobSearchParams represents parameters for searching jobs
type JobSearchParams struct {
	Title           string   `form:"title" example:"Golang Developer"`
	Location        string   `form:"location" example:"Remote"`
	JobType         string   `form:"job_type" example:"full_time"`
	Skills          []string `form:"skills" example:"Go,PostgreSQL"`
	MinSalary       float64  `form:"min_salary" example:"50000"`
	ExperienceLevel string   `form:"experience_level" validate:"omitempty,oneof='Entry-level' 'Mid-level' 'Senior' 'Lead'" example:"Mid-level"`
	Category        string   `form:"category"  example:"Engineering"`
	Page            int      `form:"page,default=1" example:"1"`
	Limit           int      `form:"limit,default=10" example:"10"`
	EmployerID      *int     `form:"employer_id" example:"12"`
}
