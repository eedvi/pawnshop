package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"pawnshop/internal/config"
	"pawnshop/internal/handler"
	"pawnshop/internal/middleware"
	"pawnshop/internal/pdf"
	"pawnshop/internal/repository/postgres"
	"pawnshop/internal/service"
	"pawnshop/pkg/auth"
	"pawnshop/pkg/cache"
	"pawnshop/pkg/metrics"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Setup logger
	setupLogger(cfg.App.Debug)

	log.Info().
		Str("app", cfg.App.Name).
		Str("version", cfg.App.Version).
		Str("environment", cfg.App.Environment).
		Msg("Starting application")

	// Connect to database
	db, err := postgres.NewDB(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()
	log.Info().Msg("Connected to database")

	// Initialize Redis cache (optional - continues without cache if unavailable)
	var redisCache *cache.Cache
	redisCache, err = cache.New(cache.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		Prefix:   "pawnshop",
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to connect to Redis - caching disabled")
		redisCache = nil
	} else {
		defer redisCache.Close()
		log.Info().Msg("Connected to Redis cache")
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	branchRepo := postgres.NewBranchRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)
	customerRepo := postgres.NewCustomerRepository(db)
	itemRepo := postgres.NewItemRepository(db)
	loanRepo := postgres.NewLoanRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	saleRepo := postgres.NewSaleRepository(db)
	cashRegisterRepo := postgres.NewCashRegisterRepository(db)
	cashSessionRepo := postgres.NewCashSessionRepository(db)
	cashMovementRepo := postgres.NewCashMovementRepository(db)
	settingRepo := postgres.NewSettingRepository(db)
	auditRepo := postgres.NewAuditLogRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)

	// New repositories for transfers, expenses, and notifications
	transferRepo := postgres.NewTransferRepository(db)
	expenseRepo := postgres.NewExpenseRepository(db)
	expenseCategoryRepo := postgres.NewExpenseCategoryRepository(db)
	notificationRepo := postgres.NewNotificationRepository(db)
	notificationTemplateRepo := postgres.NewNotificationTemplateRepository(db)
	notificationPreferenceRepo := postgres.NewCustomerNotificationPreferenceRepository(db)
	internalNotificationRepo := postgres.NewInternalNotificationRepository(db)
	twoFactorRepo := postgres.NewTwoFactorRepository(db)
	loyaltyRepo := postgres.NewLoyaltyRepository(db)

	// Initialize auth components
	jwtManager := auth.NewJWTManager(auth.JWTConfig{
		Secret:          cfg.JWT.Secret,
		AccessTokenTTL:  cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.JWT.RefreshTokenTTL,
		Issuer:          cfg.JWT.Issuer,
	})
	passwordManager := auth.NewPasswordManager()

	// Initialize services
	authService := service.NewAuthService(userRepo, roleRepo, refreshTokenRepo, jwtManager, passwordManager, log.Logger)
	userService := service.NewUserService(userRepo, roleRepo, branchRepo, passwordManager)
	customerService := service.NewCustomerService(customerRepo, branchRepo)
	itemService := service.NewItemService(itemRepo, branchRepo, categoryRepo, customerRepo)
	loanService := service.NewLoanService(loanRepo, itemRepo, customerRepo, paymentRepo, log.Logger)
	paymentService := service.NewPaymentService(paymentRepo, loanRepo, customerRepo, itemRepo, log.Logger)
	saleService := service.NewSaleService(saleRepo, itemRepo, customerRepo, branchRepo)
	cashService := service.NewCashService(cashRegisterRepo, cashSessionRepo, cashMovementRepo, branchRepo)
	branchService := service.NewBranchService(branchRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// Use cached services when Redis is available
	roleService := service.NewCachedRoleService(roleRepo, redisCache)
	settingService := service.NewCachedSettingService(settingRepo, redisCache)
	auditService := service.NewAuditService(auditRepo)

	// New services for transfers, expenses, and notifications
	transferService := service.NewTransferService(transferRepo, itemRepo, branchRepo)
	expenseService := service.NewExpenseService(expenseRepo, expenseCategoryRepo, branchRepo)
	notificationService := service.NewNotificationService(
		notificationRepo,
		notificationTemplateRepo,
		notificationPreferenceRepo,
		internalNotificationRepo,
		customerRepo,
		userRepo,
	)
	twoFactorService := service.NewTwoFactorService(twoFactorRepo, userRepo, passwordManager, cfg.App.Name)
	loyaltyService := service.NewLoyaltyService(customerRepo, loyaltyRepo)

	// Initialize storage service
	storagePath := filepath.Join(".", "storage")
	storageBaseURL := "/storage" // URL path for serving images
	storageService := service.NewStorageService(storagePath, storageBaseURL)

	// Initialize backup service
	backupPath := filepath.Join(".", "backups")
	backupService := service.NewBackupService(&cfg.Database, backupPath, log.Logger)

	// Initialize PDF generator
	pdfGenerator := pdf.NewGenerator(cfg.App.Name, "", "")
	reportService := service.NewReportService(loanRepo, paymentRepo, saleRepo, customerRepo, itemRepo, pdfGenerator)

	// Initialize audit logger
	auditLogger := middleware.NewAuditLogger(auditService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, auditLogger, log.Logger)
	userHandler := handler.NewUserHandler(userService, auditLogger)
	customerHandler := handler.NewCustomerHandler(customerService, auditLogger)
	itemHandler := handler.NewItemHandler(itemService, auditLogger)
	loanHandler := handler.NewLoanHandler(loanService, auditLogger, log.Logger)
	paymentHandler := handler.NewPaymentHandler(paymentService, auditLogger, log.Logger)
	saleHandler := handler.NewSaleHandler(saleService, auditLogger)
	cashHandler := handler.NewCashHandler(cashService, auditLogger)
	branchHandler := handler.NewBranchHandler(branchService, auditLogger)
	categoryHandler := handler.NewCategoryHandler(categoryService, auditLogger)
	roleHandler := handler.NewRoleHandler(roleService, auditLogger)
	reportHandler := handler.NewReportHandler(reportService)
	settingHandler := handler.NewSettingHandler(settingService, auditLogger)
	auditHandler := handler.NewAuditHandler(auditService)

	// New handlers for transfers, expenses, and notifications
	transferHandler := handler.NewTransferHandler(transferService)
	expenseHandler := handler.NewExpenseHandler(expenseService, auditLogger)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	twoFactorHandler := handler.NewTwoFactorHandler(twoFactorService, userService)
	loyaltyHandler := handler.NewLoyaltyHandler(loyaltyService)
	storageHandler := handler.NewStorageHandler(storageService, itemService)
	backupHandler := handler.NewBackupHandler(backupService)

	// Initialize middleware
	loggingMiddleware := middleware.NewLoggingMiddleware(log.Logger)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager, userRepo, roleRepo, log.Logger)
	rateLimiter := middleware.NewRateLimiter(middleware.DefaultRateLimitConfig())
	loginRateLimiter := middleware.NewRateLimiter(middleware.LoginRateLimitConfig())

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:               cfg.App.Name,
		ReadTimeout:           cfg.Server.ReadTimeout,
		WriteTimeout:          cfg.Server.WriteTimeout,
		IdleTimeout:           cfg.Server.IdleTimeout,
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(loggingMiddleware.Logger())
	app.Use(loggingMiddleware.Recovery())
	app.Use(compress.New())
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.CORS(middleware.DefaultCORSConfig()))
	app.Use(rateLimiter.Middleware())
	app.Use(metrics.Middleware())

	// Metrics endpoint (Prometheus)
	metricsHandler := handler.NewMetricsHandler()
	metricsHandler.RegisterRoutes(app)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"version": cfg.App.Version,
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Swagger documentation
	app.Get("/swagger.json", func(c *fiber.Ctx) error {
		return c.SendFile("api/swagger.json")
	})

	// Swagger UI
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Type("html").SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Pawnshop API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: "/swagger.json",
            dom_id: '#swagger-ui',
            presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
            layout: "BaseLayout"
        });
    </script>
</body>
</html>`)
	})

	// API routes
	api := app.Group("/api/v1")

	// Apply login rate limiter to login endpoint
	api.Post("/auth/login", loginRateLimiter.Middleware(), authHandler.Login)

	// Register routes
	authHandler.RegisterRoutes(api, authMiddleware)
	userHandler.RegisterRoutes(api, authMiddleware)
	customerHandler.RegisterRoutes(api, authMiddleware)
	itemHandler.RegisterRoutes(api, authMiddleware)
	loanHandler.RegisterRoutes(api, authMiddleware)
	paymentHandler.RegisterRoutes(api, authMiddleware)
	saleHandler.RegisterRoutes(api, authMiddleware)
	cashHandler.RegisterRoutes(api, authMiddleware)
	branchHandler.RegisterRoutes(api, authMiddleware)
	categoryHandler.RegisterRoutes(api, authMiddleware)
	roleHandler.RegisterRoutes(api, authMiddleware)
	reportHandler.RegisterRoutes(api, authMiddleware)
	settingHandler.RegisterRoutes(api, authMiddleware)
	auditHandler.RegisterRoutes(api, authMiddleware)

	// New routes for transfers, expenses, notifications, and 2FA
	transferHandler.RegisterRoutes(api, authMiddleware)
	expenseHandler.RegisterRoutes(api, authMiddleware)
	notificationHandler.RegisterRoutes(api, authMiddleware)
	twoFactorHandler.RegisterRoutes(api, authMiddleware)
	loyaltyHandler.RegisterRoutes(api, authMiddleware)
	storageHandler.RegisterRoutes(app, api, authMiddleware)
	backupHandler.RegisterRoutes(api, authMiddleware)

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "The requested resource was not found",
			},
		})
	})

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		log.Info().Str("address", addr).Msg("Server starting")
		if err := app.Listen(addr); err != nil {
			log.Fatal().Err(err).Msg("Server error")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	if err := app.ShutdownWithTimeout(30 * time.Second); err != nil {
		log.Error().Err(err).Msg("Error during shutdown")
	}

	log.Info().Msg("Server stopped")
}

func setupLogger(debug bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Log the error with details
	requestID := c.Get("X-Request-ID")
	logEvent := log.Error().
		Err(err).
		Str("request_id", requestID).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Int("status", code)

	// Add user ID if available
	if userID := c.Locals("user_id"); userID != nil {
		if uid, ok := userID.(int64); ok {
			logEvent.Int64("user_id", uid)
		}
	}

	logEvent.Msg("Unhandled error in request")

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "ERROR",
			"message": message,
		},
	})
}
