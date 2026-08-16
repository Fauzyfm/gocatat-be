package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(c *fiber.Ctx) error {

	token := c.Cookies("access_token")

    if token == "" {
        authHeader := c.Get("Authorization")
        if authHeader != "" {
            parts := strings.Split(authHeader, " ")
            if len(parts) == 2 && parts[0] == "Bearer" {
                token = parts[1]
            }
        }
    }


	// Validasi token — sama seperti sebelumnya
    parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fiber.NewError(401, "Algoritma token tidak valid")
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })


	if err != nil || !parsedToken.Valid {
        return c.Status(401).JSON(fiber.Map{
            "success": false,
            "message": "Token tidak valid atau sudah expired",
        })
    }


	claims := parsedToken.Claims.(jwt.MapClaims)
    c.Locals("userID", uint(claims["user_id"].(float64)))
    c.Locals("email",  claims["email"].(string))
    c.Locals("role",   claims["role"].(string))

	return c.Next()

}