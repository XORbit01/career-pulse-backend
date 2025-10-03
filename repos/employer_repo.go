package repos

import (
	"database/sql"
	"errors"
	"time"

	"github.com/XORbit01/jobseeker-backend/models"
)

// EmployerRepository handles database operations for employer profiles
type EmployerRepository struct {
	db *sql.DB
}

// NewEmployerRepository creates a new EmployerRepository
func NewEmployerRepository(db *sql.DB) *EmployerRepository {
	return &EmployerRepository{db: db}
}

// Create creates a new employer profile
func (r *EmployerRepository) Create(userID int, profile models.EmployerProfileInput) (int, error) {
	query := `
		INSERT INTO employer_profiles (
			user_id, company_name, industry, website, description, 
			logo_url, location, company_size, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)
		RETURNING id
	`

	now := time.Now()
	var id int
	err := r.db.QueryRow(query,
		userID,
		profile.CompanyName,
		profile.Industry,
		profile.Website,
		profile.Description,
		profile.LogoURL,
		profile.Location,
		profile.CompanySize,
		now,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByUserID retrieves an employer profile by user ID
func (r *EmployerRepository) GetByUserID(userID int) (*models.EmployerProfile, error) {
	query := `
		SELECT id, user_id, company_name, industry, website, description, 
			   logo_url, location, company_size, created_at, updated_at
		FROM employer_profiles
		WHERE user_id = $1
	`

	var profile models.EmployerProfile
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.CompanyName,
		&profile.Industry,
		&profile.Website,
		&profile.Description,
		&profile.LogoURL,
		&profile.Location,
		&profile.CompanySize,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}

	return &profile, nil
}

// GetByID retrieves an employer profile by ID
func (r *EmployerRepository) GetByID(id int) (*models.EmployerProfile, error) {
	query := `
		SELECT id, user_id, company_name, industry, website, description, 
			   logo_url, location, company_size, created_at, updated_at
		FROM employer_profiles
		WHERE id = $1
	`

	var profile models.EmployerProfile
	err := r.db.QueryRow(query, id).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.CompanyName,
		&profile.Industry,
		&profile.Website,
		&profile.Description,
		&profile.LogoURL,
		&profile.Location,
		&profile.CompanySize,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}

	return &profile, nil
}

// Update updates an employer profile
func (r *EmployerRepository) Update(id int, profile models.EmployerProfileInput) error {
	query := `
		UPDATE employer_profiles
		SET company_name = $1, industry = $2, website = $3, description = $4,
			logo_url = $5, location = $6, company_size = $7, updated_at = $8
		WHERE id = $9
	`

	_, err := r.db.Exec(query,
		profile.CompanyName,
		profile.Industry,
		profile.Website,
		profile.Description,
		profile.LogoURL,
		profile.Location,
		profile.CompanySize,
		time.Now(),
		id,
	)

	return err
}

// Delete deletes an employer profile
func (r *EmployerRepository) Delete(id int) error {
	query := `DELETE FROM employer_profiles WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
