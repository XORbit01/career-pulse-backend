package models

// UnifiedPublicEmployerProfile for public employer output
type UnifiedPublicEmployerProfile struct {
	UserID      int    `json:"id" example:"456"`
	ProfileType string `json:"profileType" example:"employer"`
	CompanyName string `json:"company_name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Industry    string `json:"industry"`
	CompanySize string `json:"company_size"`
	Website     string `json:"website"`
	LogoURL     string `json:"logo_url"`
}

// UnifiedPublicJobSeekerProfile for public job seeker output
type UnifiedPublicJobSeekerProfile struct {
	UserID          int      `json:"id" example:"123"`
	ProfileType     string   `json:"profileType" example:"job_seeker"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	Headline        string   `json:"headline"`
	Bio             string   `json:"bio"`
	Location        string   `json:"location"`
	ExperienceLevel string   `json:"experience_level"`
	Website         string   `json:"website"`
	ResumeURL       string   `json:"resume_url"`
	LogoURL         string   `json:"logo_url"`
	Skills          []string `json:"skills"`
}
