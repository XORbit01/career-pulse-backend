package models

import "time"

type EmployerProfile struct {
	ID          int       `json:"id" example:"1"`
	UserID      int       `json:"user_id" example:"42"`
	CompanyName string    `json:"company_name" example:"Tech Innovations Inc."`
	Industry    string    `json:"industry" example:"Information Technology"`
	Website     string    `json:"website" example:"https://www.techinnovations.com"`
	Description string    `json:"description" example:"Leading software company focused on building scalable backend systems and cloud solutions."`
	LogoURL     string    `json:"logo_url" example:"/uploads/logos/company123.png"`
	Location    string    `json:"location" example:"San Francisco, CA"`
	CompanySize string    `json:"company_size" example:"50-200 employees"`
	CreatedAt   time.Time `json:"created_at" example:"2025-04-14T10:18:32Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-04-14T10:18:32Z"`
}

type EmployerProfileInput struct {
	CompanyName string `json:"company_name" binding:"required" example:"Tech Innovations Inc."`
	Industry    string `json:"industry" example:"Information Technology"`
	Website     string `json:"website" example:"https://www.techinnovations.com"`
	Description string `json:"description" example:"We build secure, scalable, and cloud-native applications."`
	LogoURL     string `json:"logo_url" example:"/uploads/logos/company123.png"`
	Location    string `json:"location" example:"San Francisco, CA"`
	CompanySize string `json:"company_size" example:"50-200 employees"`
}
