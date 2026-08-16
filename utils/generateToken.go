package utils

import (
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