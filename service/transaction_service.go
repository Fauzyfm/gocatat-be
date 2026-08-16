package service

import (
	"context"
	"database/sql"
	"fmt"
	"manajemen-keuangan-api/model"
	"manajemen-keuangan-api/repository"
	"time"
)

type TransactionService struct {
	db          *sql.DB
	repoTrx     repository.TransactionRepository
	repoBalance repository.BalanceRepository
}

func NewTransactionService(
	db *sql.DB,
	repoTrx repository.TransactionRepository,
	repoBalance repository.BalanceRepository,
) *TransactionService {
	return &TransactionService{
		db:          db,
		repoTrx:     repoTrx,
		repoBalance: repoBalance,
	}
}

func (s *TransactionService) CreateTransaction(
	ctx context.Context,
	userID uint,
	balanceID uint,
	trxType model.BalanceType,
	amount int64,
	category model.CategoryType,
	description string,
) (*model.Transaction, error) {

	// validasi input
	if amount <= 0 {
		return nil, fmt.Errorf("amount harus lebih dari 0")
	}

	if description == "" {
		return nil, fmt.Errorf("Description wajib diisi!")
	}

	if trxType != model.BalanceTypeCash && trxType != model.BalanceTypeNonCash {
		return nil, fmt.Errorf("type hanya bisa cash / non cash")
	}

	if category != model.CategoryTypeIncome && category != model.CategoryTypeExpense {
		return nil, fmt.Errorf("category hanya bisa income / expense")
	}

	// cek balance milik user teresbut
	_, err := s.repoBalance.GetBalanceByIDAndUserID(ctx, balanceID, userID)
	if err != nil {
		return nil, err
	}

	// check apakah amount balance
	amountBlance, err := s.repoBalance.CheckAmountBalance(ctx, userID, balanceID)
	if err != nil {
		return nil, err
	}
	// cek saldo cukup untuk expense/pengeluaran
	if category == model.CategoryTypeExpense {
		if amountBlance < amount {
			return nil, fmt.Errorf("saldo tidak cukup: saldo Rp%d, dibutuhkan Rp%d", amountBlance, amount)
		}
	}

	// check point
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi databse: %w", err)
	}

	defer tx.Rollback()

	trx := &model.Transaction{
		UserID:      userID,
		BalanceID:   balanceID,
		Type:        trxType,
		Amount:      amount,
		Category:    category,
		Description: description,
	}

	if err := s.repoTrx.CreateTransactionTx(ctx, tx, trx); err != nil {
		return nil, err
	}

	// if err := s.repoBalance.UpdateBalanceTX(ctx, tx, balanceID, userID, delta); err != nil {
	// 	return nil, err
	// }

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("gagal menyimpan: %w", err)
	}

	return trx, nil

}

func (s *TransactionService) GetAllTransactionByUserID(ctx context.Context, userID uint) ([]model.Transaction, error) {
	return s.repoTrx.GetAllTransactionByUserID(ctx, userID)
}

func (s *TransactionService) GetTransactionByIDandUserID(ctx context.Context, id uint, userID uint) (*model.Transaction, error) {
	trx, err := s.repoTrx.GetTransactionByIDandUserID(ctx, id, userID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("transaksi tidak ditemukan")
	}

	return trx, err
}

func (s *TransactionService) GetSummary(ctx context.Context, userID uint, start, end time.Time) (model.Summary, error) {
	if end.Before(start) {
		return model.Summary{}, fmt.Errorf("tanggal akhir tidak boleh sebelum tanggal awal")
	}
	return s.repoTrx.GetSummary(ctx, userID, start, end)
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, trx *model.Transaction) error {
	// validasi enum via switch agar jelas nilai validnya
	switch trx.Type {
	case model.BalanceTypeCash, model.BalanceTypeNonCash:
	default:
		return fmt.Errorf("type tidak sesuai: hanya cash / non_cash")
	}
	switch trx.Category {
	case model.CategoryTypeIncome, model.CategoryTypeExpense:
	default:
		return fmt.Errorf("category tidak sesuai: hanya income / expense")
	}
	if trx.Amount <= 0 {
		return fmt.Errorf("amount harus lebih dari 0")
	}
	if trx.Description == "" {
		return fmt.Errorf("description wajib diisi")
	}

	// pastikan transaksi ada dan milik user; error asli tetap dibawa dengan %w
	if _, err := s.repoTrx.GetTransactionByIDandUserID(ctx, trx.ID, trx.UserID); err != nil {
		return fmt.Errorf("transaction tidak ditemukan: %w", err)
	}

	return s.repoTrx.UpdateTransaction(ctx, trx)
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, id uint, userID uint) error {
	return s.repoTrx.DeleteTransaction(ctx, id, userID)
}

func (s *TransactionService) GetAllTransactionPaginated(ctx context.Context, userid uint, page int, limit int) (model.PaginatedTransactions, error) {
	return s.repoTrx.GetAllTransactionPaginated(ctx, userid, page, limit)
}

// endOfDay mengembalikan akhir hari (23:59:59.999999999) dari sebuah tanggal.
// Dipakai agar filter end_date bersifat INKLUSIF untuk satu hari penuh,
// bukan hanya sampai jam 00:00:00 hari tersebut.
func endOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

func (s *TransactionService) FilterDateTransaction(ctx context.Context, userid uint, page int, limit int, typeTrx model.BalanceType, categoryType model.CategoryType, startDate time.Time, endDate time.Time) (model.PaginatedTransactions, error) {

	// validasi domain: nilai enum tidak valid ditolak eksplisit,
	// bukan dibiarkan menghasilkan hasil kosong secara senyap
	if typeTrx != "" && typeTrx != model.BalanceTypeCash && typeTrx != model.BalanceTypeNonCash {
		return model.PaginatedTransactions{}, fmt.Errorf("type hanya bisa cash / non_cash")
	}
	if categoryType != "" && categoryType != model.CategoryTypeIncome && categoryType != model.CategoryTypeExpense {
		return model.PaginatedTransactions{}, fmt.Errorf("category hanya bisa income / expense")
	}
	if !startDate.IsZero() && !endDate.IsZero() && endDate.Before(startDate) {
		return model.PaginatedTransactions{}, fmt.Errorf("tanggal akhir tidak boleh sebelum tanggal awal")
	}

	// endOfDay agar filter end_date bersifat inklusif satu hari penuh;
	// normalisasi page/limit ditangani oleh repository
	if !endDate.IsZero() {
		endDate = endOfDay(endDate)
	}

	filterTrx := model.TransactionFilter{
		StartDate: startDate,
		EndDate:   endDate,
		Type:      typeTrx,
		Category:  categoryType,
		Limit:     limit,
	}

	return s.repoTrx.FilterDateTransaction(ctx, userid, page, filterTrx)
}
