package repository

import (
	"context"
	"database/sql"
	"fmt"
	"manajemen-keuangan-api/model"

)

type BalanceRepository interface {
	CreateBalance(ctx context.Context, balance *model.Balance) error
	GetAllBalanceByUserID(ctx context.Context ,userID uint) ([]model.Balance, error)
	GetBalanceByID(ctx context.Context, id uint) (*model.Balance, error)
	UpdateBalance(ctx context.Context, balance *model.Balance) error
	DeleteBalance(ctx context.Context, id uint, userID uint) error
	// UpdateBalanceTX(ctx context.Context, dbtx DBTX, balanceID uint, userID uint, delta int64) error
	GetBalanceByIDAndUserID(ctx context.Context, id uint, userID uint) (*model.Balance, error)
	CheckAmountBalance(ctx context.Context, userid uint, balace uint) (int64, error)
}

type PostgresBalanceRepository struct {
	db *sql.DB
}

func NewPostgresBalanceRepository(db *sql.DB) *PostgresBalanceRepository {
	return &PostgresBalanceRepository{db: db}
}

func (r *PostgresBalanceRepository) CreateBalance(ctx context.Context, balance *model.Balance) error {

	query := `INSERT INTO balances (user_id, wallet, type, created_at, update_at) VALUES ($1, $2, $3, now(), now())
	RETURNING id, created_at, update_at
	`

	err := r.db.QueryRowContext(ctx, query, balance.UserID, balance.Wallet, balance.Type,).Scan(&balance.ID, &balance.CreatedAt, &balance.UpdateAt)

	if err != nil {
		return fmt.Errorf("gagal membuat balance: %w", err)
	}

	return nil
}

func (r *PostgresBalanceRepository) GetAllBalanceByUserID(ctx context.Context, userID uint) ([]model.Balance, error) {	
    query := `
    SELECT
      b.id, b.user_id, b.wallet, b.type,
      COALESCE(SUM(
        CASE
          WHEN t.category = 'income' THEN t.amount
          WHEN t.category = 'expense' THEN -t.amount
          ELSE 0
        END
      ), 0) AS amount,
      b.created_at, b.update_at
    FROM balances b
    LEFT JOIN transactions t ON b.id = t.balance_id
    WHERE b.user_id = $1
    GROUP BY b.id, b.user_id, b.wallet, b.type, b.created_at, b.update_at
    ORDER BY b.created_at DESC
    `


	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar balance: %w", err)
	}

	defer rows.Close()

	var balances []model.Balance
	for rows.Next() {
		var b model.Balance
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.Wallet, &b.Type, &b.Amount, &b.CreatedAt, &b.UpdateAt, 
		); err != nil {
			return nil, fmt.Errorf("gagal membaca daftar balance: %w", err)
		}
		balances = append(balances, b)
	}


	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error saat iterasi balance: %w", err)
	}

	return balances, nil
}

func (r *PostgresBalanceRepository) GetBalanceByID(ctx context.Context, id uint) (*model.Balance, error) {
    query := `
    SELECT
      b.id, b.user_id, b.wallet, b.type,
      COALESCE(SUM(
        CASE
          WHEN t.category = 'income' THEN t.amount
          WHEN t.category = 'expense' THEN -t.amount
          ELSE 0
        END
      ), 0) AS amount,
      b.created_at, b.update_at
    FROM balances b
    LEFT JOIN transactions t ON b.id = t.balance_id
    WHERE b.id = $1
    GROUP BY b.id, b.user_id, b.wallet, b.type, b.created_at, b.update_at
    `

	var balance model.Balance
	err := r.db.QueryRowContext(ctx, query, id).Scan(&balance.ID, &balance.UserID, &balance.Wallet, &balance.Type, &balance.Amount, &balance.CreatedAt, &balance.UpdateAt )
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil balance: %w", err)
	}

	return &balance, nil
}

func (r *PostgresBalanceRepository) UpdateBalance(ctx context.Context, balance *model.Balance) error {
    query := `
        UPDATE balances
        SET wallet = $1, type = $2, update_at = now()
        WHERE id = $3 AND user_id = $4
    `
 
    result, err := r.db.ExecContext(
        ctx, query,
        balance.Wallet, balance.Type, balance.ID, balance.UserID,
    )
    if err != nil {
        return fmt.Errorf("gagal update balance: %w", err)
    }
 
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("gagal membaca rows affected: %w", err)
    }
    if rowsAffected == 0 {
        return fmt.Errorf("balance tidak ditemukan atau bukan milik user ini")
    }
 
    return nil
}

func (r *PostgresBalanceRepository) DeleteBalance(ctx context.Context, id uint, userID uint) error {
	// untuk proses delete row itu masih menggunakan hard delete, mksdnya dengan langsung menhapus 1 row yang akan di delete nya.
	query := `
		DELETE FROM balances WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("gagal menghapus balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal membaca rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("balance tidak ditemukan atau bukan milik user ini")
	}

	return nil
} 

// func (r *PostgresBalanceRepository) UpdateBalanceTX(ctx context.Context, dbtx DBTX, balanceID uint, userID uint, delta int64) error {
// 	query := `
// 		UPDATE balances
// 		SET amount = amount + $1, update_at = now()
// 		WHERE id = $2 AND user_id = $3
// 		AND (amount + $1) >= 0
// 	`

// 	result, err := dbtx.ExecContext(ctx, query, delta, balanceID, userID)

// 	if err != nil {
// 		return fmt.Errorf("gagal update saldo: %w", err)
// 	}

// 	rows, _ := result.RowsAffected()
// 	if rows == 0 {
// 		return fmt.Errorf("saldo tidak cukup atau balance tidak ditemukan")
// 	}

// 	return nil
// }

func (r *PostgresBalanceRepository) GetBalanceByIDAndUserID(ctx context.Context, id uint, userID uint) (*model.Balance, error) {
    query := `
    SELECT
      b.id, b.user_id, b.wallet, b.type,
      COALESCE(SUM(
        CASE
          WHEN t.category = 'income' THEN t.amount
          WHEN t.category = 'expense' THEN -t.amount
          ELSE 0
        END
      ), 0) AS amount,
      b.created_at, b.update_at
    FROM balances b
    LEFT JOIN transactions t ON b.id = t.balance_id
    WHERE b.id = $1 AND b.user_id = $2
    GROUP BY b.id, b.user_id, b.wallet, b.type, b.created_at, b.update_at
    `

	var b model.Balance
	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(&b.ID, &b.UserID, &b.Wallet, &b.Type, &b.Amount, &b.CreatedAt, &b.UpdateAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("balance tidak ditemukan atau bukan milik user ini")
	}
    if err != nil {
        return nil, fmt.Errorf("gagal mengambil balance: %w", err)
    }
    return &b, nil

}

func (r *PostgresBalanceRepository) CheckAmountBalance(ctx context.Context, userid uint, balanceid uint) (int64, error) {

	query := `
	SELECT COALESCE(SUM(
	CASE
		WHEN category = 'income' THEN amount
		WHEN category = 'expense' THEN -amount
		ELSE 0
	END
	), 0) as total_balance
	FROM transactions
	WHERE user_id = $1 AND balance_id = $2
	`

	var total_balance int64
	err := r.db.QueryRowContext(ctx, query, userid, balanceid).Scan(&total_balance)
	if err != nil {
		return 0, fmt.Errorf("gagal melakuakan pengambilan data total balance")
	}

	return total_balance, nil

}