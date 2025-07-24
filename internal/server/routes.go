package server

import (
	"os"
	"rent-application/apps/auth"
	"rent-application/apps/payment"
	"rent-application/apps/product"
	"rent-application/apps/upload"
	"rent-application/apps/users"
	config "rent-application/configs"
	"rent-application/internal/cloudinary"
	"rent-application/internal/database"
	"rent-application/internal/email"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func SetupRouter(
	app *fiber.App,
	cfg config.Config,
	log *zerolog.Logger,
) {
	// ===================== MIDDLEWARE GLOBAL =====================
	app.Use(recover.New())

	// Middleware Request ID
	app.Use(func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Locals("request_id", requestID)
		c.Set("X-Request-ID", requestID)
		return c.Next()
	})

	// Middleware Logger dengan zerolog
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := c.Locals("request_id").(string)
		logger := log.With().Str("request_id", requestID).Logger()

		// Simpan logger di context
		c.Locals("logger", &logger)

		err := c.Next()

		// Log hanya jika status < 400 (tidak error)
		status := c.Response().StatusCode()
		if status < 400 {
			logger.Info().
				Str("method", c.Method()).
				Str("path", c.Path()).
				Int("status", status).
				Dur("latency", time.Since(start)).
				Msg("request")
		}

		return err
	})

	// ===================== ROUTE HEALTHCHECK =====================
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
		})
	})

	// ===================== INISIALISASI SERVICE =====================
	// === Inisialisasi service dan controller ===
	// Inisialisasi dan daftarkan auth domain

	db := database.NewPostgreDatabase(cfg.DB)
	store := database.NewStore(db, cfg.DB)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	authQuery := auth.NewAuthQuery()
	authRepo := auth.NewAuthRepository(store, authQuery)
	emailClient := email.NewResendClient(email.ResendConfig{ApiKey: cfg.Email.ApiKey, From: cfg.Email.From}, &logger)
	authService := auth.NewAuthService(authRepo, emailClient)
	authController := auth.NewAuthController(authService)
	authController.Route(app)

	productQuery := product.NewProductQuery()
	productRepo := product.NewProductRepository(store, productQuery)
	productService := product.NewProductService(productRepo)
	categoryController := product.NewCategoryController(productService)
	categoryController.Route(app)

	cloudinaryCfg := cloudinary.NewCloudinaryCfg(cfg.Cloudinary)
	uploadService := upload.NewUploadService(cloudinaryCfg)
	uploadController := upload.NewUploadController(uploadService)
	uploadController.Route(app)

	paymentQuery := payment.NewPaymentQuery()
	paymentRepo := payment.NewPaymentRepository(store, paymentQuery, productQuery)
	paymentService := payment.NewPaymentService(paymentRepo, productRepo, authRepo)
	paymentController := payment.NewPaymentController(paymentService)
	paymentController.Route(app)

	userQuery := users.NewUsersQuery()
	userRepo := users.NewUsersRepository(store, userQuery)
	userService := users.NewUsersService(userRepo)
	userController := users.NewUserController(userService)
	userController.Route(app)

	// Inisialisasi modul lain (product, users, payment, dll) di sini
	// productController := ...
	// usersController := ...
	// paymentController := ...
	// ...

	// === Struktur grup route (aktifkan jika diperlukan) ===
	// api := app.Group("/api/v1")
	// authGroup := api.Group("/auth")
	// authGroup.Post("/login", authController.Login)
	// authGroup.Post("/register", authController.Register)
	// authGroup.Post("/refresh", authController.RefreshToken)

	// usersGroup := api.Group("/users")
	// usersGroup.Use(web.JWTMiddleware(cfg.JWT.Secret))
	// usersGroup.Get("/", usersController.GetAllUsers)
	// usersGroup.Get(":id", usersController.GetUserByID)
	// usersGroup.Put(":id", usersController.UpdateUser)

	// paymentGroup := api.Group("/payment")
	// paymentGroup.Post("/pay", paymentController.Pay)

	// ... grup lain sesuai kebutuhan

	// ===================== HANDLER 404 =====================
	// app.Use(func(c *fiber.Ctx) error {
	// 	return web.NotFoundResponse(c)
	// })
}

// // Contoh inisialisasi controller
// func initAuthController(cfg *config.Config, log *zerolog.Logger) *auth.Controller {
// 	authRepo := auth.NewRepository(database.GetDB())
// 	authService := auth.NewService(authRepo, cfg.JWT)
// 	return auth.NewController(authService, log)
// }

// func initUserController(cfg *config.Config, log *zerolog.Logger) *users.Controller {
// 	userRepo := users.NewRepository(database.GetDB())
// 	userService := users.NewService(userRepo, cfg.Cloudinary)
// 	return users.NewController(userService, log)
// }
