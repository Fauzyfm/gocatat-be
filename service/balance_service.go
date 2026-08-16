package service

import (
	"context"
	"fmt"
	"manajemen-keuangan-api/model"
	"manajemen-keuangan-api/repository"
)

type BalanceService struct {
	repoBalance repository.BalanceRepository
}

func NewBalanceService(repoBalance repository.BalanceRepository) *BalanceService {
	return &BalanceService{repoBalance: repoBalance}
}

func (s *BalanceService) CreateBalance(ctx context.Context, UserID uint, Wallet string, Type model.BalanceType) (*model.Balance, error){

	if Wallet == "" || Type == "" {
		return nil, fmt.Errorf("Wallet dan Type harus terisi!")
	}

	if (Type != model.BalanceTypeCash && Type != model.BalanceTypeNonCash) {
		return nil, fmt.Errorf("Type hanya bisa cash / nonCash")
	}

	newBalance := &model.Balance{
		UserID: UserID,
		Wallet: Wallet,
		Type: Type,
	}

	if err := s.repoBalance.CreateBalance(ctx, newBalance); err != nil {
		return nil, err
	}

	return newBalance, nil

}

func (s *BalanceService) GetAllBalanceByUserID(ctx context.Context, UserID uint) ([]model.Balance, error) {


	balances, err := s.repoBalance.GetAllBalanceByUserID(ctx, UserID)
	if err != nil {
		return nil, err
	}

	return balances, nil

}


func (s *BalanceService) GetBalanceByID(ctx context.Context, id uint, userID uint) (*model.Balance, error) {
	return s.repoBalance.GetBalanceByIDAndUserID(ctx, id, userID)
}

func (s *BalanceService) UpdateBalance(ctx context.Context, balance *model.Balance) error {
	existing, err := s.repoBalance.GetBalanceByIDAndUserID(ctx, balance.ID, balance.UserID)
	if err != nil {
		return err
	}

	balance.ID = existing.ID
	return s.repoBalance.UpdateBalance(ctx, balance)
}

func (s *BalanceService) DeleteBalance(ctx context.Context, id uint, userID uint) error {
	return s.repoBalance.DeleteBalance(ctx, id, userID)
}