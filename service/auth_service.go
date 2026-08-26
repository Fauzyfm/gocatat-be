package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"manajemen-keuangan-api/model"
	"manajemen-keuangan-api/repository"
	"manajemen-keuangan-api/utils"

	"golang.org/x/crypto/bcrypt"
)

type EmailSender interface {
    SendVerification(to, token string) error
}

type AuthService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) Register(ctx context.Context, userName string, email string, password string, confirmPassword string) (*model.User, error) {

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

	claims := &model.User{
		UserName: userName,
		Email: email,
	}
	verificationToken, err := utils.GenerateTokenVerification(claims)
	if err != nil {
		return nil, err
	}

	var provider string
	if provider == "" {
		provider = "local"
	} 



	user := &model.User{
		UserName: userName,
		Email: email,
		Password: string(hash),
		Role: "user",
		Provider: provider,
		IsVerified: false,
		VerificationToken: verificationToken,
	}




	if err := s.authRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("gagal menyimpan user!")
	}

	go func() {
        if err := utils.SendVerificationEmail(user.Email, verificationToken); err != nil {
            log.Printf("[WARN] gagal kirim email verifikasi ke %s: %v", user.Email, err)
        }
    }()

	return user, nil

}


func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error){

	user, err := s.authRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("email atau password salah!")
	}

	if user.IsVerified != true {
		return "", fmt.Errorf("User kamu belum di verifikasi, mohon untuk check kembali email anda atau register kembali menggunakan email tersebut")
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


func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	
	email, err := utils.ParseVerificationToken(token)
	if err != nil {
		return fmt.Errorf("token tidak valid atau sudah kadaluwarsa, silahkan melakukan register kembali")
	}

	user, err := s.authRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan")
	}

	if user.VerificationToken != token {
		return fmt.Errorf("token verifikasi tidak cocok atau sudah digunakan sebelumnya")
	}

	if user.IsVerified {
		return fmt.Errorf("akun sudah terverifikasi, silahkan login")
	}

	return s.authRepo.UpdateVerifiedUser(ctx, user.ID)

}


func (s *AuthService) ResendVerificationToken(ctx context.Context, email string) error {
	user, err := s.authRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan!")
	}

	if user.IsVerified {
		return fmt.Errorf("user sudah terverifikasi, silahkan login!")
	}

	claims := &model.User{
		UserName: user.UserName,
		Email: user.Email,
	}
	NewVerificationToken, err := utils.GenerateTokenVerification(claims)
	if err != nil {
		return fmt.Errorf("gagal melakukan pembuatan verification token")
	}

	if err := s.authRepo.UpdateVerificationToken(ctx, user.ID, NewVerificationToken); err != nil {
		return fmt.Errorf("gagal menyimpen token baru: %w", err)
	}

	go func() {
        if err := utils.SendVerificationEmail(user.Email, NewVerificationToken); err != nil {
            log.Printf("[WARN] gagal kirim email verifikasi ke %s: %v", user.Email, err)
        }
    }()

    return nil

}


func (s *AuthService) ChangePasswordEmail(ctx context.Context, email string) error {
	user, err := s.authRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan!")
	}

	if user.IsVerified != true {
		return fmt.Errorf("pastikan user sudah terverifikasi terlebih dahulu!")
	}

	claims := &model.User{
		UserName: user.UserName,
		Email: user.Email,
	}
	TokenChangePass, err := utils.GenerateChangePasswordToken(claims)
	if err != nil {
		return fmt.Errorf("gagal pemmbuatan token: %w", err)
	}

	if err := s.authRepo.UpdateVerificationToken(ctx, user.ID, TokenChangePass); err != nil {
		return fmt.Errorf("gagal menyimpan token: %w", err)
	}

	go func (){
		if err := utils.SendChangePasswordEmail(email, TokenChangePass); err != nil {
            log.Printf("[WARN] gagal kirim email verifikasi ke %s: %v", user.Email, err)
		}
	}()

	return nil
}

func (s *AuthService) ChangePassword(ctx context.Context, token, NewPassword, ConfirmPassword string) error {

	email, err := utils.ParseVerificationToken(token)
	if err != nil {
		return fmt.Errorf("token sudah tidak valid atau kadaluwarsa, silahkan lakukan reset password ulang")
	}

	user, err := s.authRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan!")
	}

	if user.IsVerified != true {
		return fmt.Errorf("pastikan user sudah terverifikasi terlebih dahulu!")
	}

	if user.VerificationToken != token || user.VerificationToken == "" {
		return	fmt.Errorf("token tidak cocok atau sudah digunakan")
	}

	if NewPassword != ConfirmPassword {
		return fmt.Errorf("Password dan Confirm Password harus sama")
	}

	if len(NewPassword) < 8 {
		return fmt.Errorf("password minimal 8 karakter")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(NewPassword), 12)
	if err != nil {
		return fmt.Errorf("gagal memproses password")
	}

	return s.authRepo.UpdatePassword(ctx, string(hash), user.ID)
 
}


func (s *AuthService) GoogleUser(ctx context.Context, username, email string) (*model.User, error) {
	existingUser, err := s.authRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.authRepo.CreateGoogleUser(ctx, username, email)
		}
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	if existingUser.Provider != "google" {
		return nil, fmt.Errorf("Emmail sudah terdaftar dengan provider %s, silahkan login dengan metode tersebut", existingUser.Provider)
	}

	if !existingUser.IsVerified {
		if err := s.authRepo.MarkGoogleVerified(ctx, existingUser.ID); err != nil {
			return nil, fmt.Errorf("gagal memverifikasi user google: %w", err )
		}
	}

	return existingUser, nil

}