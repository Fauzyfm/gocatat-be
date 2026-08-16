package main

import (
	"errors"
	"log"
	"manajemen-keuangan-api/config"
	"manajemen-keuangan-api/handler"
	"manajemen-keuangan-api/middleware"
	"manajemen-keuangan-api/repository"
	"manajemen-keuangan-api/service"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {

	// Load .env di paling awal sebelum apapun yang pakai os.Getenv
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	// Connect DB
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Gagal connect database") // Tidak print detail error (bisa berisi DSN/password)
	}

	// Dependency
	// Repository Layer
	userRepo := repository.NewPostgresAuthRepository(db)
	balanceRepo := repository.NewPostgresBalanceRepository(db)
	trxRepo := repository.NewPostgresTransactionRepository(db)

	// Service Layer
	authService := service.NewAuthService(userRepo)
	balanceService := service.NewBalanceService(balanceRepo)
	trxService := service.NewTransactionService(db, trxRepo, balanceRepo)

	// Handler Layer
	authHandler := handler.NewAuthHandler(authService)
	balanceHandler := handler.NewBalanceHandler(balanceService)
	trxHandler := handler.NewTransactionHandler(trxService)

	// Setup Go Fiber
	app := fiber.New(fiber.Config{
		AppName: "Manajemen Keuangan API v1.0 (gocatat)",

		// Error handler global — tangkap semua error yang tidak di-handle
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Untuk error 500, jangan expose detail internal ke client
			message := err.Error()
			if code == fiber.StatusInternalServerError {
				log.Printf("[ERROR] Internal: %v", err) // Log detail hanya ke server terminal
				message = "Terjadi kesalahan pada server"
			}

			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": message,
			})
		},
	})

	// Middleware Global
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} ${method} ${path} - ${latency}\n",
	}))

	// CORS — konfigurasi dari .env
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:3000"
	}
	corsMethods := os.Getenv("CORS_METHODS")
	if corsMethods == "" {
		corsMethods = "GET,POST,PUT,DELETE,PATCH"
	}
	corsHeaders := os.Getenv("CORS_HEADERS")
	if corsHeaders == "" {
		corsHeaders = "Content-Type,Authorization"
	}
	corsCredentials := strings.ToLower(os.Getenv("CORS_CREDENTIALS")) == "true"

	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     corsMethods,
		AllowHeaders:     corsHeaders,
		AllowCredentials: corsCredentials,
	}))

	// Rate Limiter — konfigurasi dari .env
	rateLimitMax, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_MAX"))
	if rateLimitMax <= 0 {
		rateLimitMax = 30
	}
	rateLimitExp, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_EXPIRATION_SECONDS"))
	if rateLimitExp <= 0 {
		rateLimitExp = 60
	}

	app.Use(limiter.New(limiter.Config{
		Max:        rateLimitMax,
		Expiration: time.Duration(rateLimitExp) * time.Second,
	}))

	// Health Check — untuk monitoring EasyPanel & uptime check
	app.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "ok"
		if err := db.Ping(); err != nil {
			dbStatus = "down"
		}

		status := "healthy"
		httpCode := 200
		if dbStatus != "ok" {
			status = "unhealthy"
			httpCode = 503
		}

		return c.Status(httpCode).JSON(fiber.Map{
			"status":   status,
			"database": dbStatus,
			"env":      os.Getenv("APP_ENV"),
		})
	})

	// Routes
	api := app.Group("/api/v1")

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	// Balance Routes
	balance := api.Group("/balance", middleware.AuthRequired)
	balance.Post("/", balanceHandler.CreateBalance)
	balance.Get("/", balanceHandler.GetAllBalanceByUserID)
	balance.Get("/:id", balanceHandler.GetBalanceByID)
	balance.Put("/:id", balanceHandler.UpdateBalance)
	balance.Delete("/:id", balanceHandler.DeleteBalance)

	// Transaction Routes
	trx := api.Group("/transaction", middleware.AuthRequired)
	trx.Post("/", trxHandler.CreateTransaction)
	trx.Get("/", trxHandler.GetAllTransactions)
	trx.Get("/summary", trxHandler.GetSummary)
	trx.Get("/:id", trxHandler.GetTransactionByIDandUserID)
	trx.Put("/:id", trxHandler.UpdateTransaction)
	trx.Delete("/:id", trxHandler.DeleteTransaction)

	// Get me Routes
	protected := api.Group("", middleware.AuthRequired)
	protected.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"user_id": c.Locals("userID"),
				"role":    c.Locals("role"),
				"email":   c.Locals("email"),
			},
		})
	})

	// Jalankan Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Log startup info yang aman (tidak print secrets)
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}
	log.Printf("Server berjalan di port %s [env: %s]", port, appEnv)
	log.Printf("CORS origins: %s", corsOrigins)
	log.Fatal(app.Listen(":" + port))

}
