package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"pawnshop/internal/config"
	"pawnshop/internal/repository/postgres"
	"pawnshop/internal/scheduler"
	"pawnshop/internal/service"
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
		Str("app", cfg.App.Name+"-worker").
		Str("version", cfg.App.Version).
		Str("environment", cfg.App.Environment).
		Msg("Starting worker")

	// Connect to database
	db, err := postgres.NewDB(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()
	log.Info().Msg("Connected to database")

	// Initialize repositories
	loanRepo := postgres.NewLoanRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	customerRepo := postgres.NewCustomerRepository(db)
	userRepo := postgres.NewUserRepository(db)
	loyaltyRepo := postgres.NewLoyaltyRepository(db)

	// Initialize notification repositories
	notificationRepo := postgres.NewNotificationRepository(db)
	notificationTemplateRepo := postgres.NewNotificationTemplateRepository(db)
	notificationPreferenceRepo := postgres.NewCustomerNotificationPreferenceRepository(db)
	internalNotificationRepo := postgres.NewInternalNotificationRepository(db)

	// Initialize services
	notificationService := service.NewNotificationService(
		notificationRepo,
		notificationTemplateRepo,
		notificationPreferenceRepo,
		internalNotificationRepo,
		customerRepo,
		userRepo,
	)
	loyaltyService := service.NewLoyaltyService(customerRepo, loyaltyRepo)

	// Initialize scheduler
	sched := scheduler.New(log.Logger)

	// Initialize job service
	jobService := scheduler.NewJobService(
		loanRepo,
		paymentRepo,
		customerRepo,
		notificationService,
		loyaltyService,
		log.Logger,
	)

	// Register default jobs
	scheduler.RegisterDefaultJobs(sched, jobService)

	// Start scheduler
	sched.Start()
	log.Info().Msg("Worker started")

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down worker...")
	sched.Stop()
	log.Info().Msg("Worker stopped")
}

func setupLogger(debug bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
