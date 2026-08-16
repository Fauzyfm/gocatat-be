package model

import "time"

type BalanceType string

const (
	BalanceTypeNonCash    BalanceType = "nonCash"
	BalanceTypeCash    BalanceType = "cash"
)

type Balance struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userID"`
	Wallet    string     `json:"wallet"`
	Type      BalanceType `json:"type"`
	Amount    int64       `json:"amount"`
	CreatedAt time.Time	  `json:"createdAt"`
	UpdateAt  time.Time	  `json:"updateAt"`	
}


