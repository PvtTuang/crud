package authentication

import (
	"database/sql"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *User) error
	FindByID(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByEmail(email string) (*User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.Exec(`INSERT INTO users (username ,password_hash , email , created_at , updated_at )
	 VALUES ($1, $2, $3, NOW(), NOW())`, user.Username, user.PasswordHash, user.Email,
	)
	if err != nil {
		return err
	}

	return tx.Commit()

}

func (r *userRepository) FindByID(id uuid.UUID) (*User, error) {
	var user User
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	row := tx.QueryRow(`SELECT id , username ,password_hash , email , created_at , updated_at FROM users WHERE id = $1`, id)
	err = row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return &User{}, err
	}
	return &user, err
}

func (r *userRepository) FindByUsername(username string) (*User, error) {
	var user User
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	row := tx.QueryRow(`SELECT id , username ,password_hash , email , created_at , updated_at FROM users WHERE username = $1`, username)
	err = row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return &User{}, err
	}
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*User, error) {
	var user User
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	row := tx.QueryRow(`SELECT id , username ,password_hash , email , created_at , updated_at FROM users WHERE email = $1`, email)
	err = row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return &User{}, err
	}
	return &user, err
}
