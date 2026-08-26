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
	UpdateVerifiedUser(ctx context.Context, userid uint) error
	UpdateVerificationToken(ctx context.Context, userid uint, NewVerificationToken string) error
	UpdatePassword(ctx context.Context, NewPassword string, userid uint) error
	CreateGoogleUser(ctx context.Context, username, email string) (*model.User, error)
	MarkGoogleVerified(ctx context.Context, userid uint) error
}

type PostgresAuthRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(db *sql.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}


func (r *PostgresAuthRepository) Create(ctx context.Context, user *model.User) error {
    query := `
        INSERT INTO users (username, email, password, role, provider, is_verified, verification_token, created_at, update_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
        RETURNING id, created_at, update_at
    `

	err := r.db.QueryRowContext(ctx, query, user.UserName, user.Email, user.Password, user.Role, user.Provider, user.IsVerified, user.VerificationToken).Scan(&user.ID, &user.CreatedAt, &user.UpdateAt)

	if err != nil {
		return fmt.Errorf("gagal memuat user: %w", err)
	}

	return nil
}

func (r *PostgresAuthRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
SELECT id, username, email, password, role, provider, is_verified, verification_token, created_at, update_at, deleted_at 
FROM users WHERE email = $1
	`

	var user model.User
	var password sql.NullString // ← tambahkan ini

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.UserName, &user.Email, &password, &user.Role,
		&user.Provider, &user.IsVerified, &user.VerificationToken,
		&user.CreatedAt, &user.UpdateAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	if password.Valid {
		user.Password = password.String
	}
	// kalau tidak Valid, user.Password otomatis tetap "" (zero value) — aman

	return &user, nil
}

func (r *PostgresAuthRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	query := `
		SELECT id, username, email, password, role, created_at, update_at, deleted_at FROM users
		WHERE id = $1
	`


	var user model.User
	var password sql.NullString
err := r.db.QueryRowContext(ctx, query, id).Scan(
    &user.ID, &user.UserName, &user.Email, &user.Password, &user.Role,
    &user.CreatedAt, &user.UpdateAt, &user.DeletedAt,
)

	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	if password.Valid {
		user.Password = password.String
	}

	return &user, nil
}


func (r *PostgresAuthRepository) UpdateVerifiedUser(ctx context.Context, userid uint) error {
	query := `
		UPDATE users
		SET is_verified = true, verification_token = '', update_at = now()
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, userid)
	if err != nil {
		return fmt.Errorf("gagal memverifikasi user: %w", err)	
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user tidak ditemukan")
	}

	return nil

}

func (r *PostgresAuthRepository) UpdateVerificationToken(ctx context.Context, userid uint, NewVerificationToken string) error {
	query := `
		UPDATE users
		SET verification_token = $1
		WHERE id = $2
	`

	res, err := r.db.ExecContext(ctx, query, NewVerificationToken, userid)
	if err != nil {
		return fmt.Errorf("gagal update verification token: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user tidak ditemukan!")
	}

	return nil
}


func (r *PostgresAuthRepository) UpdatePassword(ctx context.Context, NewPassword string, userid uint) error {
	query := `
	UPDATE users
	SET password = $1, verification_token = '', update_at = now()
	WHERE id = $2
	`

	res, err := r.db.ExecContext(ctx, query, NewPassword, userid)
	if err != nil {
		return fmt.Errorf("gagal update password: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user tidak ditemukan")
	}

	return nil
}

func (r *PostgresAuthRepository) CreateGoogleUser(ctx context.Context, username, email string) (*model.User, error){
	query := `
        INSERT INTO users (username, email, password, role, provider, is_verified, verification_token, created_at, update_at)
        VALUES ($1, $2, NULL, 'user', 'google', true, '', now(), now())
        RETURNING id, created_at, update_at, username, email, role, is_verified, provider
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, username, email).Scan(&user.ID, &user.CreatedAt, &user.UpdateAt, &user.UserName, &user.Email, &user.Role, &user.IsVerified, &user.Provider,)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat user: %w", err)
	}

	return &user, nil
}

func (r *PostgresAuthRepository) MarkGoogleVerified(ctx context.Context, userid uint) error {
	query := `
		UPDATE users 
		SET is_verified = true, update_at = now()
		WHERE id = $1 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, userid)
	if err != nil {
		return fmt.Errorf("gagal mengupdate google verified")
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user tidak ditemukan")
	}

	return nil
}