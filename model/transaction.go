package model

import "time"

type CategoryType string

const (
	CategoryTypeIncome CategoryType = "income" // income artinya pendapatan
	CategoryTypeExpense CategoryType = "expense"  // expense artinya pengeluaran
)

type Transaction struct {
	ID uint `json:"id"`

	UserID    uint    `json:"userID"`
	BalanceID uint    `json:"balanceID"`

	Type      BalanceType  `json:"type"`
	Amount    int64        `json:"amount"`
	Category  CategoryType `json:"category"`
	Description string	`json:"description"`

	CreatedAt time.Time	`json:"createdAt"`
}

type Summary struct {
    Income  int64
    Expense int64
    AllBalance int64
}

type TransactionFilter struct {
	StartDate time.Time `json:"start_date"`
	EndDate		time.Time `json:"end_date"`
	Type BalanceType	`json:"type"`
	Category CategoryType `json:"category"`
	Limit int `json:"limit"`
}

type PaginatedTransactions struct {
	Data []Transaction `json:"data"`
	Page int	`json:"page"`
	Limit int `json:"limit"`
	TotalItems int	`json:"total_items"`
	TotalPages int	`json:"total_pages"`
}