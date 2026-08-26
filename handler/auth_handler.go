package handler

import (
	"context"
	"encoding/json"
	"io"
	"manajemen-keuangan-api/config"
	"manajemen-keuangan-api/model"
	"manajemen-keuangan-api/service"
	"manajemen-keuangan-api/utils"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
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
	}

	if err := c.BodyParser(&inputRegister); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"messagge": "Format Tidak Valid",
		})
	}

	user, err := h.svcAuth.Register(c.Context(), inputRegister.UserName, strings.ToLower(inputRegister.Email), inputRegister.Password, inputRegister.ConfirmPassword)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}  


	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"messagge": "Register Berhasil, silahkan check email anda untuk melakukan verifikasi!",
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

	token, err := h.svcAuth.Login(c.Context(), strings.ToLower(inputLogin.Email), inputLogin.Password)
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
		Path: "/",
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
		Path: "/",
    })

    return c.JSON(fiber.Map{
        "success": true,
        "messagge": "Logout berhasil",
    })
}

// endpoint : /api/v1/verify?token=.....
func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "token tidak ditemukan",
		})
	}

	if err := h.svcAuth.VerifyEmail(c.Context(), token); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"messagge": "Verifikasi berhasil, silahkan logiin!",
	})
}	


// endpoint : /api/v1/resend-verificationToken
func (h *AuthHandler) ResendVerificationToken(c *fiber.Ctx) error {
	var inputBody struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&inputBody); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format tidak sesuai!",
		})
	}

	if err := h.svcAuth.ResendVerificationToken(c.Context(), inputBody.Email); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"messagge": "Verification Email sudah dikirim ulang, silahkan check email anda",
	})
}

// endpoint: /api/v1/verification-change-password
func (h *AuthHandler) VerificationChangePassword(c *fiber.Ctx) error {
	var inputBody struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&inputBody); err != nil {		
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format tidak sesuai",
		})
	}

	if err := h.svcAuth.ChangePasswordEmail(c.Context(), strings.ToLower(inputBody.Email)); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"messagge": "Silahkan check email anda untuk melakukan perubahan password",
	})

}


// endpoint: /api/v1/change-password {post}
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var inputBody struct {
		Token string `json:"token"`
		NewPassword string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.BodyParser(&inputBody); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format tidak sesuai",
		})
	}

	if err := h.svcAuth.ChangePassword(c.Context(), inputBody.Token, inputBody.NewPassword, inputBody.ConfirmPassword); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"messagge": "Password berhasil dirubah, silahkan melakukan login kembali",
	})
}

func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	state := utils.GenerateRandomString(24)
	verifier, challenge := utils.GeneratePKCE()

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   600, // 10 menit
		HTTPOnly: true,
		Secure:   getCookieSecure(),
		SameSite: getCookieSameSite(),
	})
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_verifier",
		Value:    verifier,
		MaxAge:   600,
		HTTPOnly: true,
		Secure:   getCookieSecure(),
		SameSite: getCookieSameSite(),
	})

	url := config.GoogleOauthConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.AccessTypeOffline,
	)
	return c.Redirect(url)
}


func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {

	var feURL string
	if os.Getenv("APP_ENV") == "development" {
		feURL = "http://localhost:8888"
	} else {
		feURL = os.Getenv("FE_URL")
	}

	// 1.Validasi state (CSRF)
	if c.Cookies("oauth_state") == "" || c.Cookies("oauth_state") != c.Query("state") {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "State tidak valid",
		})
	}

	code := c.Query("code")
	verifier := c.Cookies("oauth_verifier")

	// 2. Tukar code jadi token
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", verifier))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "gagal menukar authorization code",
		})
	}
	c.ClearCookie("oauth_state", "oauth_verifier")

	// 3. ambil profil user google
	client := config.GoogleOauthConfig.Client(context.Background(), token)

	res, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"messagge": "gagal mengambil profile user",
		})
	}

	defer res.Body.Close()
	body, _  := io.ReadAll(res.Body)
	var info model.GoogleUserInfo
	json.Unmarshal(body, &info)

	if !info.EmailVerified {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Email google belim terverifikasi",

		})
	}

	// 4. FIND and CREATE pada Table users
	user, err := h.svcAuth.GoogleUser(c.Context(), info.Name, info.Email)
	if err != nil {
		return c.Redirect(feURL + "/login?error=" + url.QueryEscape(err.Error()))
	}

	// 5. Terbitkan sesi 
	accessToken, err := utils.GenerateTokenJWT(user)
	if err != nil {
		return c.Redirect(feURL + "/login?error=gagal_membuat_token")
	}

	c.Cookie(&fiber.Cookie{
		Name: "access_token",
		Value: accessToken,
		MaxAge: 60 * 60 * 24,
		HTTPOnly: true,
		Secure: getCookieSecure(),
		SameSite: getCookieSameSite(),
		Domain: getCookieDomain(),
		Path: "/",
	})

	return c.Redirect(feURL)
}