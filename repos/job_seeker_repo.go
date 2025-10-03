package repos

import (
	"database/sql"
	"errors"
	"time"

	"github.com/XORbit01/jobseeker-backend/models"
	"github.com/lib/pq"
)

// JobSeekerRepository handles database operations for job seeker profiles
type JobSeekerRepository struct {
	db *sql.DB
}

// NewJobSeekerRepository creates a new JobSeekerRepository
func NewJobSeekerRepository(db *sql.DB) *JobSeekerRepository {
	return &JobSeekerRepository{db: db}
}

// Create creates a new job seeker profile
func (r *JobSeekerRepository) Create(userID int, profile models.JobSeekerProfileInput) (int, error) {
	query := `
	INSERT INTO job_seeker_profiles (
		user_id, first_name, last_name, headline, summary, phone, location, 
		resume_url, logo_url, skills, experience_level, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7,
		$8, $9, $10, $11, $12, $13
	)
	RETURNING id
	`

	now := time.Now()
	var id int
	err := r.db.QueryRow(query,
		userID,
		profile.FirstName,
		profile.LastName,
		profile.Headline,
		profile.Summary,
		profile.Phone,
		profile.Location,
		profile.ResumeURL,
		profile.LogoUrl,
		pq.Array(profile.Skills),
		profile.ExperienceLevel,
		now,
		now,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByUserID retrieves a job seeker profile by user ID
func (r *JobSeekerRepository) GetByUserID(userID int) (*models.JobSeekerProfile, error) {
	query := `
		SELECT id, user_id, first_name, last_name, headline, summary, phone, location, resume_url,logo_url,skills,experience_level, created_at, updated_at
		FROM job_seeker_profiles
		WHERE user_id = $1
	`

	var profile models.JobSeekerProfile
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Headline,
		&profile.Summary,
		&profile.Phone,
		&profile.Location,
		&profile.ResumeURL,
		&profile.LogoUrl,
		pq.Array(&profile.Skills),
		&profile.ExperienceLevel,
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

// GetByID retrieves a job seeker profile by ID
func (r *JobSeekerRepository) GetByID(id int) (*models.JobSeekerProfile, error) {
	query := `
		SELECT id, user_id, first_name, last_name, headline, summary, phone, location, resume_url,logo_url, skills,experience_level, created_at, updated_at
		FROM job_seeker_profiles
		WHERE id = $1
	`

	var profile models.JobSeekerProfile
	err := r.db.QueryRow(query, id).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Headline,
		&profile.Summary,
		&profile.Phone,
		&profile.Location,
		&profile.ResumeURL,
		&profile.LogoUrl,
		pq.Array(&profile.Skills),
		&profile.ExperienceLevel, // <-- add this
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

func (r *JobSeekerRepository) Update(id int, input models.JobSeekerProfileInput) error {
	query := `
	UPDATE job_seeker_profiles
	SET first_name = $1, last_name = $2, phone = $3, location = $4,
		headline = $5, summary = $6, resume_url = $7, logo_url = $8,
		skills = $9, experience_level = $10, updated_at = NOW()
	WHERE id = $11
`
	_, err := r.db.Exec(query,
		input.FirstName,
		input.LastName,
		input.Phone,
		input.Location,
		input.Headline,
		input.Summary,
		input.ResumeURL,
		input.LogoUrl,
		pq.Array(input.Skills),
		input.ExperienceLevel,
		id,
	)
	return err
}

// Delete deletes a job seeker profile
func (r *JobSeekerRepository) Delete(id int) error {
	query := `DELETE FROM job_seeker_profiles WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
