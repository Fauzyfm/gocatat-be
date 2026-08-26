package utils

import (
	"fmt"
	"manajemen-keuangan-api/model"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenJWT(user *model.User) (string, error) {
	expHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))

	if expHours == 0 {
		expHours = 24
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email": user.Email,
		"role": user.Role,
		"exp": time.Now().Add(time.Duration(expHours) * time.Hour).Unix(), 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString([]byte (os.Getenv("JWT_SECRET")))
}

func GenerateTokenVerification(user *model.User) (string, error) {
	expHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))

	if expHours == 0 {
		expHours = 24
	}

	claims := jwt.MapClaims{
		"email": user.Email,
		"exp": time.Now().Add(time.Duration(expHours) * time.Hour).Unix(), 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString([]byte (os.Getenv("JWT_SECRET")))
}

func ParseVerificationToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing tidak terduga")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("token tidak valid atau kadaluwarsa")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("claims tidak valid")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return "", fmt.Errorf("claims email tidak ditemukan")
	}

	return email, nil
}

func GenerateChangePasswordToken(user *model.User) (string, error) {
	expHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))

	if expHours == 0 {
		expHours = 24
	}
	claims := jwt.MapClaims{
		"email": user.Email,
		"exp": time.Now().Add(time.Duration(expHours) * time.Hour). Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte (os.Getenv("JWT_SECRET")))
}