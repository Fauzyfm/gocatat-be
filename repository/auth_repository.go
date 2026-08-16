package repository

import (
	"context"
	"database/sql"
	"fmt"
	"manajemen-keuangan-api/model"

)

type AuthRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id uint) (*model.User, error)
}

type PostgresAuthRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(db *sql.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}


func (r *PostgresAuthRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, email, password, role, created_at, update_at, deleted_at) VALUES ($1, $2, $3, $4, now(), now(), NULL)
		RETURNING id, created_at, update_at
	`

	err := r.db.QueryRowContext(ctx, query, user.UserName, user.Email, user.Password, user.Role).Scan(&user.ID, &user.CreatedAt, &user.UpdateAt)

	if err != nil {
		return fmt.Errorf("gagal memuat user: %w", err)
	}

	return nil
}

func (r *PostgresAuthRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
	SELECT id, username, email, password, role, created_at, update_at, deleted_at FROM users
          WHERE email = $1
	`
	
	var user model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.UserName, &user.Email, &user.Password, &user.Role,
		&user.CreatedAt, &user.UpdateAt, &user.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	return &user, nil
}

func (r *PostgresAuthRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	query := `
		SELECT id, username, email, password, role, created_at, update_at, deleted_at FROM users
		WHERE id = $1
	`


	var user model.User
err := r.db.QueryRowContext(ctx, query, id).Scan(
    &user.ID, &user.UserName, &user.Email, &user.Password, &user.Role,
    &user.CreatedAt, &user.UpdateAt, &user.DeletedAt,
)

	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	return &user, nil
}