package repos

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/XORbit01/jobseeker-backend/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user models.UserInput, passwordHash string) (int, error) {
	query := `
		INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id
	`

	now := time.Now()
	var id int
	err := r.db.QueryRow(query, user.Email, passwordHash, user.Role, now).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Update updates a user in the database
func (r *UserRepository) Update(id int, user models.User) error {
	query := `
		UPDATE users
		SET email = $1, role = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(query, user.Email, user.Role, time.Now(), id)
	return err
}

// Delete deletes a user from the database
func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
