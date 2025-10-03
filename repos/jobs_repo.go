package repos

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/XORbit01/jobseeker-backend/models"
	"github.com/lib/pq"
)

// JobRepository handles database operations for jobs
type JobRepository struct {
	db *sql.DB
}

// NewJobRepository creates a new JobRepository
func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

// Create creates a new job
func (r *JobRepository) Create(employerID int, job models.JobInput) (int, error) {
	query := `
		INSERT INTO jobs (
			employer_id, title, description, location, job_type, 
			salary_min, salary_max, experience_level, required_skills, 
			category, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	now := time.Now()
	status := job.Status
	if status == "" {
		status = "active"
	}

	var id int
	err := r.db.QueryRow(query,
		employerID,
		job.Title,
		job.Description,
		job.Location,
		job.JobType,
		job.SalaryMin,
		job.SalaryMax,
		job.ExperienceLevel,
		pq.Array(job.RequiredSkills),
		job.Category,
		status,
		now,
		now,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByID retrieves a job by ID
func (r *JobRepository) GetByID(id int) (*models.Job, error) {
	query := `
		SELECT j.id, e.user_id AS employer_user_id, j.title, j.description, j.location, 
			j.job_type, j.salary_min, j.salary_max, j.experience_level,
			j.required_skills, j.status, j.created_at, j.updated_at,
			e.company_name, j.category, e.logo_url
		FROM jobs j
		JOIN employer_profiles e ON j.employer_id = e.id
		WHERE j.id = $1
	`

	var job models.Job
	var category sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&job.ID,
		&job.EmployerID,
		&job.Title,
		&job.Description,
		&job.Location,
		&job.JobType,
		&job.SalaryMin,
		&job.SalaryMax,
		&job.ExperienceLevel,
		pq.Array(&job.RequiredSkills),
		&job.Status,
		&job.CreatedAt,
		&job.UpdatedAt,
		&job.CompanyName,
		&category,
		&job.LogoURL,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("job not found")
		}
		return nil, err
	}

	if category.Valid {
		job.Category = category.String
	}

	return &job, nil
}

// Update updates a job
func (r *JobRepository) Update(id int, job models.JobInput) error {
	query := `
		UPDATE jobs
		SET title = $1, description = $2, location = $3, job_type = $4,
			salary_min = $5, salary_max = $6, experience_level = $7,
			required_skills = $8, category = $9, status = $10, updated_at = $11
		WHERE id = $12
	`

	status := job.Status
	if status == "" {
		status = "active"
	}

	_, err := r.db.Exec(query,
		job.Title,
		job.Description,
		job.Location,
		job.JobType,
		job.SalaryMin,
		job.SalaryMax,
		job.ExperienceLevel,
		pq.Array(job.RequiredSkills),
		job.Category,
		status,
		time.Now(),
		id,
	)

	return err
}

// Delete deletes a job
func (r *JobRepository) Delete(id int) error {
	query := `DELETE FROM jobs WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// IsJobOwner checks if an employer owns a job
func (r *JobRepository) IsJobOwner(jobID, employerID int) (bool, error) {
	query := `SELECT COUNT(*) FROM jobs WHERE id = $1 AND employer_id = $2`

	var count int
	err := r.db.QueryRow(query, jobID, employerID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// SearchJobs searches for jobs with filters
func (r *JobRepository) SearchJobs(params models.JobSearchParams) ([]*models.Job, int, error) {
	whereConditions := []string{"j.status = 'active'"}
	args := []any{}
	argCount := 1

	if params.Title != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("j.title ILIKE $%d", argCount))
		args = append(args, "%"+params.Title+"%")
		argCount++
	}

	if params.Location != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("j.location ILIKE $%d", argCount))
		args = append(args, "%"+params.Location+"%")
		argCount++
	}

	if params.JobType != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("j.job_type = $%d", argCount))
		args = append(args, params.JobType)
		argCount++
	}

	if params.MinSalary > 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("j.salary_min >= $%d", argCount))
		args = append(args, params.MinSalary)
		argCount++
	}

	if params.ExperienceLevel != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("j.experience_level = $%d", argCount))
		args = append(args, params.ExperienceLevel)
		argCount++
	}

	if params.EmployerID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("e.user_id = $%d", argCount))
		args = append(args, *params.EmployerID)
		argCount++
	}

	if params.Category != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("j.category = $%d", argCount))
		args = append(args, params.Category)
		argCount++
	}

	if len(params.Skills) > 0 {
		for _, skill := range params.Skills {
			whereConditions = append(whereConditions, fmt.Sprintf("$%d = ANY(j.required_skills)", argCount))
			args = append(args, skill)
			argCount++
		}
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM jobs j
		JOIN employer_profiles e ON j.employer_id = e.id
		%s
	`, whereClause)

	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	query := fmt.Sprintf(`
		SELECT j.id, j.employer_id, j.title, j.description, j.location, 
			   j.job_type, j.salary_min, j.salary_max, j.experience_level,
			   j.required_skills, j.status, j.created_at, j.updated_at,
			   e.company_name, j.category, e.logo_url
		FROM jobs j
		JOIN employer_profiles e ON j.employer_id = e.id
		%s
		ORDER BY j.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		var job models.Job
		var salaryMin, salaryMax sql.NullFloat64
		var category sql.NullString

		err := rows.Scan(
			&job.ID,
			&job.EmployerID,
			&job.Title,
			&job.Description,
			&job.Location,
			&job.JobType,
			&salaryMin,
			&salaryMax,
			&job.ExperienceLevel,
			pq.Array(&job.RequiredSkills),
			&job.Status,
			&job.CreatedAt,
			&job.UpdatedAt,
			&job.CompanyName,
			&category,
			&job.LogoURL,
		)
		if err != nil {
			return nil, 0, err
		}
		if salaryMin.Valid {
			job.SalaryMin = &salaryMin.Float64
		}
		if salaryMax.Valid {
			job.SalaryMax = &salaryMax.Float64
		}
		if category.Valid {
			job.Category = category.String
		}
		jobs = append(jobs, &job)
	}

	return jobs, total, rows.Err()
}

// GetByEmployerID retrieves jobs posted by a specific employer
func (r *JobRepository) GetByEmployerID(employerID int, page, limit int) ([]*models.Job, int, error) {
	offset := (page - 1) * limit

	// Total count query
	var total int
	countQuery := `SELECT COUNT(*) FROM jobs WHERE employer_id = $1`
	err := r.db.QueryRow(countQuery, employerID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT j.id, j.employer_id, j.title, j.description, j.location, 
			   j.job_type, j.salary_min, j.salary_max, j.experience_level,
			   j.required_skills, j.status, j.created_at, j.updated_at,
			   e.company_name, j.category, e.logo_url
		FROM jobs j
		JOIN employer_profiles e ON j.employer_id = e.id
		WHERE j.employer_id = $1
		ORDER BY j.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, employerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		var job models.Job
		var salaryMin, salaryMax sql.NullFloat64
		var category sql.NullString

		err := rows.Scan(
			&job.ID,
			&job.EmployerID,
			&job.Title,
			&job.Description,
			&job.Location,
			&job.JobType,
			&salaryMin,
			&salaryMax,
			&job.ExperienceLevel,
			pq.Array(&job.RequiredSkills),
			&job.Status,
			&job.CreatedAt,
			&job.UpdatedAt,
			&job.CompanyName,
			&category,
			&job.LogoURL,
		)
		if err != nil {
			return nil, 0, err
		}
		if salaryMin.Valid {
			job.SalaryMin = &salaryMin.Float64
		}
		if salaryMax.Valid {
			job.SalaryMax = &salaryMax.Float64
		}
		if category.Valid {
			job.Category = category.String
		}
		jobs = append(jobs, &job)
	}

	return jobs, total, nil
}
