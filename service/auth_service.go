package service

import (
	"context"
	"fmt"
	"manajemen-keuangan-api/model"
	"manajemen-keuangan-api/repository"
	"manajemen-keuangan-api/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) Register(ctx context.Context, userName string, email string, password string, confirmPassword string, role model.RoleType) (*model.User, error) {

	if userName == "" || email == "" || password == "" {
		return nil, fmt.Errorf("username, email, password wajib untuk di isi!")
	}

	if password != confirmPassword {
		return nil, fmt.Errorf("password dan confirm password harus sama!")
	}

	if len(password) < 8 {
		return nil, fmt.Errorf("Password minimal 8 karakter!")
	}


	checkEmail, _ := s.authRepo.GetByEmail(ctx, email)
	if checkEmail != nil {
		return nil, fmt.Errorf("email sudah pernah terdaftar!")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte (password), 12)
	if err != nil {
		return nil, fmt.Errorf("gagal memproses password!")
	}

	user := &model.User{
		UserName: userName,
		Email: email,
		Password: string(hash),
		Role: role,
	}


	if err := s.authRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("gagal menyimpan user!")
	}

	return user, nil

}


func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error){

	user, err := s.authRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("email atau password salah!")
	}

	if err := bcrypt.CompareHashAndPassword([]byte (user.Password), []byte (password)); err != nil {
		return "", fmt.Errorf("email atau password salalh!")
	}

	accessToken, err := utils.GenerateTokenJWT(user)
	if err != nil {
		 return "", err
	}


	return	accessToken, nil

}

func (s *AuthService) GetMe(ctx context.Context, userID uint) (*model.User, error) {
	exitingUser, err := s.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("User tidak ditemukan!")
	}

	user := &model.User{
		ID: exitingUser.ID,
		UserName: exitingUser.UserName,
		Role: exitingUser.Role,
	}

	return user, nil
}