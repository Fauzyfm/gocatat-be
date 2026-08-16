package repository

import (
	"context"
	"database/sql"
	"fmt"
	"manajemen-keuangan-api/model"
	"time"
)

type TransactionRepository interface {
	GetAllTransactionByUserID(ctx context.Context, userID uint) ([]model.Transaction, error)
	GetTransactionByIDandUserID(ctx context.Context, id uint, userID uint) (*model.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *model.Transaction) error
	DeleteTransaction(ctx context.Context, id uint, userID uint) error
	GetSummary(ctx context.Context, userID uint, start, end time.Time) (model.Summary, error)
	CreateTransactionTx(ctx context.Context, dbtx DBTX, transaction *model.Transaction) error
	FilterDateTransaction(ctx context.Context, userid uint, page int, filter model.TransactionFilter) (model.PaginatedTransactions, error)
	GetAllTransactionPaginated(ctx context.Context, userid uint, page int, limit int) (model.PaginatedTransactions, error)
}

type PostgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) GetAllTransactionByUserID(ctx context.Context, userID uint) ([]model.Transaction, error) {
	query := `
		SELECT id, user_id, balance_id, type, amount, category, description, created_at FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar transaction: %w", err)
	}

	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var trx model.Transaction
		if err := rows.Scan(&trx.ID, &trx.UserID, &trx.BalanceID, &trx.Type, &trx.Amount, &trx.Category, &trx.Description, &trx.CreatedAt); err != nil {
			return nil, fmt.Errorf("gagal membaca daftar transaction: %w", err)
		}
		transactions = append(transactions, trx)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error saat iteras transaction: %w", err)
	}

	return transactions, nil

}

func (r *PostgresTransactionRepository) GetTransactionByIDandUserID(ctx context.Context, id uint, userID uint) (*model.Transaction, error) {
	query := `
		SELECT id, user_id, balance_id, type, amount, category, description, created_at FROM transactions
		WHERE id = $1 AND user_id = $2
	`

	var trx model.Transaction
	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(&trx.ID, &trx.UserID, &trx.BalanceID, &trx.Type, &trx.Amount, &trx.Category, &trx.Description, &trx.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil transaction: %w", err)
	}

	return &trx, nil
}

func (r *PostgresTransactionRepository) UpdateTransaction(ctx context.Context, transaction *model.Transaction) error {
	query := `
		UPDATE transactions
		SET type = $1, amount = $2, category = $3, description = $4
		WHERE id = $5 AND user_id = $6
	`

	result, err := r.db.ExecContext(ctx, query, transaction.Type, transaction.Amount, transaction.Category, transaction.Description, transaction.ID, transaction.UserID)

	if err != nil {
		return fmt.Errorf("gagal update transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal membuat rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transaction tidak ditemukan atau bukan milik user ini")
	}

	return nil
}

func (r *PostgresTransactionRepository) DeleteTransaction(ctx context.Context, id uint, userID uint) error {
	query := `
		DELETE FROM transactions WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("gagal menghapus transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal membaca rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transaction tidak ditemukan atau bukan milik user ini")
	}

	return nil

}

func (r *PostgresTransactionRepository) GetSummary(ctx context.Context, userID uint, start, end time.Time) (model.Summary, error) {
	query := `
		SELECT
		COALESCE(SUM(CASE WHEN category = 'income' THEN amount ELSE 0 END), 0) AS income,
		COALESCE(SUM(CASE WHEN category = 'expense' THEN amount ELSE 0 END), 0) AS expense
		FROM transactions
		WHERE user_id = $1
			AND created_at >= $2
			AND created_at <= $3
	`

	var s model.Summary
	err := r.db.QueryRowContext(ctx, query, userID, start, end).Scan(&s.Income, &s.Expense)
	if err != nil {
		return model.Summary{}, fmt.Errorf("gagal mengambil summary: %w", err)
	}

	query2 := `
		SELECT COALESCE(SUM(
		CASE WHEN category = 'income' THEN amount
			WHEN category = 'expense' THEN -amount
			ELSE 0 END
		), 0) AS total_balances
		FROM transactions WHERE user_id = $1
	`

	err = r.db.QueryRowContext(ctx, query2, userID).Scan(&s.AllBalance)
	if err != nil {
		return model.Summary{}, fmt.Errorf("gagal mengambil data balances: %w", err)
	}

	return s, nil
}

func (r *PostgresTransactionRepository) CreateTransactionTx(
	ctx context.Context,
	dbtx DBTX,
	transaction *model.Transaction,
) error {
	query := `
        INSERT INTO transactions
            (user_id, balance_id, type, amount, category, description, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, now())
        RETURNING id, created_at
    `
	err := dbtx.QueryRowContext(ctx, query,
		transaction.UserID,
		transaction.BalanceID,
		transaction.Type,
		transaction.Amount,
		transaction.Category,
		transaction.Description,
	).Scan(&transaction.ID, &transaction.CreatedAt)

	if err != nil {
		return fmt.Errorf("gagal membuat transaksi: %w", err)
	}
	return nil
}

func (r *PostgresTransactionRepository) FilterDateTransaction(ctx context.Context, userid uint, page int, filter model.TransactionFilter) (model.PaginatedTransactions, error) {

	// normalisasi page & limit di repository (pertahanan terakhir):
	// mencegah panic "integer divide by zero" dan query berat jika
	// ada pemanggil yang lupa menormalisasi sebelum memanggil repo
	const (
		defaultLimit int = 20
		maxLimit     int = 100
	)
	if page <= 0 {
		page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = defaultLimit
	} else if filter.Limit > maxLimit {
		filter.Limit = maxLimit
	}

	// menghitung offset
	offset := (page - 1) * filter.Limit

	whereClause := " WHERE user_id = $1"
	args := []interface{}{userid}
	argPos := 2

	if filter.Type != "" {
		whereClause += fmt.Sprintf(" AND type = $%d", argPos)
		args = append(args, filter.Type)
		argPos++
	}
	if filter.Category != "" {
		whereClause += fmt.Sprintf(" AND category = $%d", argPos)
		args = append(args, filter.Category)
		argPos++
	}
	if !filter.StartDate.IsZero() {
		whereClause += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, filter.StartDate)
		argPos++
	}
	if !filter.EndDate.IsZero() {
		whereClause += fmt.Sprintf(" AND created_at <= $%d", argPos)
		args = append(args, filter.EndDate)
		argPos++
	}

	// menghitung total data yang match dengan filter
	countQuery := `
	SELECT COUNT(*) FROM transactions
	` + whereClause
	var totalItems int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return model.PaginatedTransactions{}, fmt.Errorf("gagal menghitung data transaction: %w", err)
	}

	// ambil data sesuai halaman dan juga match filter
	dataQuery := `
		SELECT id, user_id, balance_id, type, amount, category, description, created_at FROM transactions
	` + whereClause + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)

	dataArgs := append(args, filter.Limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return model.PaginatedTransactions{}, fmt.Errorf("Gagal mengambil data: %w", err)
	}

	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var trx model.Transaction
		if err := rows.Scan(&trx.ID, &trx.UserID, &trx.BalanceID, &trx.Type, &trx.Amount, &trx.Category, &trx.Description, &trx.CreatedAt); err != nil {
			return model.PaginatedTransactions{}, fmt.Errorf("gagal membaca data: %w", err)
		}
		transactions = append(transactions, trx)
	}

	if err := rows.Err(); err != nil {
		return model.PaginatedTransactions{}, fmt.Errorf("error saat iteras transaction: %w", err)
	}

	// menghitung page
	totalPages := (totalItems + filter.Limit - 1) / filter.Limit

	return model.PaginatedTransactions{
		Data:       transactions,
		Page:       page,
		Limit:      filter.Limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r *PostgresTransactionRepository) GetAllTransactionPaginated(ctx context.Context, userid uint, page int, limit int) (model.PaginatedTransactions, error) {

	// filter limit and page
	const (
		defaulLimit int = 20
		maxLimit    int = 100
	)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = defaulLimit
	} else if limit > maxLimit {
		limit = maxLimit
	}

	// menghitung offset
	offset := (page - 1) * limit

	var totalItems int
	countQuery := `
	SELECT COUNT(*) FROM transactions WHERE user_id = $1
	`
	if err := r.db.QueryRowContext(ctx, countQuery, userid).Scan(&totalItems); err != nil {
		return model.PaginatedTransactions{}, fmt.Errorf("gagal menghitung total transaksi %w", err)
	}

	dataQUERY := `
		SELECT id, user_id, balance_id,	type, amount, category, description, created_at FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, dataQUERY, userid, limit, offset)
	if err != nil {
		return model.PaginatedTransactions{}, fmt.Errorf("gagal mengambil data: %w", err)
	}

	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var trx model.Transaction
		if err := rows.Scan(&trx.ID, &trx.UserID, &trx.BalanceID, &trx.Type, &trx.Amount, &trx.Category, &trx.Description, &trx.CreatedAt); err != nil {
			return model.PaginatedTransactions{}, fmt.Errorf("gagal membaca data transactions: %w", err)
		}

		transactions = append(transactions, trx)
	}

	if err := rows.Err(); err != nil {
		return model.PaginatedTransactions{}, fmt.Errorf("error saat iterasi transaction: %w", err)
	}

	totalPages := (totalItems + limit - 1) / limit

	return model.PaginatedTransactions{
		Data:       transactions,
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}
