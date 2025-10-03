package repos

import (
	"database/sql"
	"errors"
	"time"

	"github.com/XORbit01/jobseeker-backend/models"
)

// ApplicationRepository handles database operations for job applications
type ApplicationRepository struct {
	db *sql.DB
}

// NewApplicationRepository creates a new ApplicationRepository
func NewApplicationRepository(db *sql.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

// Create creates a new job application
func (r *ApplicationRepository) Create(jobSeekerID int, application models.ApplicationInput) (int, error) {
	query := `
		INSERT INTO applications (job_id, job_seeker_id, cover_letter, status, created_at, updated_at)
		VALUES ($1, $2, $3, 'pending', $4, $4)
		RETURNING id
	`

	now := time.Now()
	var id int
	err := r.db.QueryRow(query,
		application.JobID,
		jobSeekerID,
		application.CoverLetter,
		now,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByID retrieves an application by ID
func (r *ApplicationRepository) GetByID(id int) (*models.Application, error) {
	query := `
		SELECT a.id, a.job_id, a.job_seeker_id, a.cover_letter, a.status, a.created_at, a.updated_at,
			   j.title as job_title, e.company_name, js.first_name, js.last_name
		FROM applications a
		JOIN jobs j ON a.job_id = j.id
		JOIN employer_profiles e ON j.employer_id = e.id
		JOIN job_seeker_profiles js ON a.job_seeker_id = js.id
		WHERE a.id = $1
	`

	var app models.Application
	err := r.db.QueryRow(query, id).Scan(
		&app.ID,
		&app.JobID,
		&app.JobSeekerID,
		&app.CoverLetter,
		&app.Status,
		&app.CreatedAt,
		&app.UpdatedAt,
		&app.JobTitle,
		&app.CompanyName,
		&app.FirstName,
		&app.LastName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("application not found")
		}
		return nil, err
	}

	return &app, nil
}

// GetByJobSeekerID retrieves applications by job seeker ID
func (r *ApplicationRepository) GetByJobSeekerID(jobSeekerID int, page, limit int) ([]*models.Application, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM applications WHERE job_seeker_id = $1`
	err := r.db.QueryRow(countQuery, jobSeekerID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get applications with pagination
	query := `
		SELECT a.id, a.job_id, a.job_seeker_id, a.cover_letter, a.status, a.created_at, a.updated_at,
			   j.title as job_title, e.company_name, js.first_name, js.last_name
		FROM applications a
		JOIN jobs j ON a.job_id = j.id
		JOIN employer_profiles e ON j.employer_id = e.id
		JOIN job_seeker_profiles js ON a.job_seeker_id = js.id
		WHERE a.job_seeker_id = $1
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, jobSeekerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	applications := make([]*models.Application, 0)
	for rows.Next() {
		var app models.Application
		err := rows.Scan(
			&app.ID,
			&app.JobID,
			&app.JobSeekerID,
			&app.CoverLetter,
			&app.Status,
			&app.CreatedAt,
			&app.UpdatedAt,
			&app.JobTitle,
			&app.CompanyName,
			&app.FirstName,
			&app.LastName,
		)
		if err != nil {
			return nil, 0, err
		}
		applications = append(applications, &app)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

// GetByJobID retrieves applications by job ID
func (r *ApplicationRepository) GetByJobID(jobID int, page, limit int) ([]*models.Application, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM applications WHERE job_id = $1`
	err := r.db.QueryRow(countQuery, jobID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get applications with pagination
	query := `
		SELECT a.id, a.job_id, js.user_id, a.cover_letter, a.status, a.created_at, a.updated_at,
			   j.title as job_title, e.company_name, js.first_name, js.last_name, js.logo_url, js.resume_url
		FROM applications a
		JOIN jobs j ON a.job_id = j.id
		JOIN employer_profiles e ON j.employer_id = e.id
		JOIN job_seeker_profiles js ON a.job_seeker_id = js.id
		WHERE a.job_id = $1
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, jobID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	applications := make([]*models.Application, 0)
	for rows.Next() {
		var app models.Application
		err := rows.Scan(
			&app.ID,
			&app.JobID,
			&app.JobSeekerID,
			&app.CoverLetter,
			&app.Status,
			&app.CreatedAt,
			&app.UpdatedAt,
			&app.JobTitle,
			&app.CompanyName,
			&app.FirstName,
			&app.LastName,
			&app.LogoURL,
			&app.ResumeURL,
		)
		if err != nil {
			return nil, 0, err
		}
		applications = append(applications, &app)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

// UpdateStatus updates an application's status
func (r *ApplicationRepository) UpdateStatus(id int, status string) error {
	query := `
		UPDATE applications
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

// Delete deletes an application
func (r *ApplicationRepository) Delete(id int) error {
	query := `DELETE FROM applications WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// IsJobSeekerApplication checks if an application belongs to a job seeker
func (r *ApplicationRepository) IsJobSeekerApplication(applicationID, jobSeekerID int) (bool, error) {
	query := `SELECT COUNT(*) FROM applications WHERE id = $1 AND job_seeker_id = $2`

	var count int
	err := r.db.QueryRow(query, applicationID, jobSeekerID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// IsJobOwnerApplication checks if an application is for a job owned by an employer
func (r *ApplicationRepository) IsJobOwnerApplication(applicationID, employerID int) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM applications a
		JOIN jobs j ON a.job_id = j.id
		WHERE a.id = $1 AND j.employer_id = $2
	`

	var count int
	err := r.db.QueryRow(query, applicationID, employerID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
