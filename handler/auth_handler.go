package handler

import (
	"manajemen-keuangan-api/model"
	"manajemen-keuangan-api/service"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	svcAuth *service.AuthService
}

func NewAuthHandler(svcAuth *service.AuthService) *AuthHandler {
	return &AuthHandler{svcAuth: svcAuth}
}

// getCookieSecure membaca COOKIE_SECURE dari .env
func getCookieSecure() bool {
	return strings.ToLower(os.Getenv("COOKIE_SECURE")) == "true"
}

// getCookieSameSite membaca COOKIE_SAMESITE dari .env
func getCookieSameSite() string {
	sameSite := os.Getenv("COOKIE_SAMESITE")
	switch strings.ToLower(sameSite) {
	case "lax":
		return "Lax"
	case "none":
		return "None"
	case "strict":
		return "Strict"
	default:
		return "Strict"
	}
}

// getCookieDomain membaca COOKIE_DOMAIN dari .env
func getCookieDomain() string {
	return os.Getenv("COOKIE_DOMAIN")
}

// endpoint: /api/v1/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var inputRegister struct {
		UserName string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
		Role model.RoleType `json:"role"`
	}

	if err := c.BodyParser(&inputRegister); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"messagge": "Format Tidak Valid",
		})
	}

	if inputRegister.Role == "" {
		inputRegister.Role = model.RoleUser
	}

	user, err := h.svcAuth.Register(c.Context(), inputRegister.UserName, inputRegister.Email, inputRegister.Password, inputRegister.ConfirmPassword, inputRegister.Role)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	} 

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"messagge": "Register Berhasil",
		"data": user,
	})
}

// endpoint: /api/v1/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var inputLogin struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	
	if err := c.BodyParser(&inputLogin); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	token, err := h.svcAuth.Login(c.Context(), inputLogin.Email, inputLogin.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		MaxAge:   60 * 60 * 24,
		HTTPOnly: true,
		Secure:   getCookieSecure(),
		SameSite: getCookieSameSite(),
		Domain:   getCookieDomain(),
	})

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"messagge": "Login Berhasil",
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {

    // Hapus kedua cookie
    c.Cookie(&fiber.Cookie{
        Name:     "access_token",
        Value:    "",
        MaxAge:   -1,
        HTTPOnly: true,
        Secure:   getCookieSecure(),
        SameSite: getCookieSameSite(),
        Domain:   getCookieDomain(),
    })

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Logout berhasil",
    })
}