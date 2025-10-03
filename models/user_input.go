package models

import "time"

type User struct {
	ID           int       `json:"id" example:"1"`
	Email        string    `json:"email" example:"user@example.com"`
	PasswordHash string    `json:"-"` // hidden from JSON
	Role         string    `json:"role" example:"job_seeker"`
	CreatedAt    time.Time `json:"created_at" example:"2025-04-14T10:18:32Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2025-04-14T10:18:32Z"`
}

type UserInput struct {
	Email    string `json:"email" binding:"required,email" example:"jane@jobportal.com"`
	Password string `json:"password" binding:"required,min=8" example:"Str0ngPass!"`
	Role     string `json:"role" binding:"required,oneof=job_seeker employer" example:"employer"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email" example:"jane@jobportal.com"`
	Password string `json:"password" binding:"required" example:"Str0ngPass!"`
}

type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // shortened JWT
	Role  string `json:"role" example:"job_seeker"`
}
