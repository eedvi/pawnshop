package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/internal/service"
)

// JobService contains dependencies for scheduled jobs
type JobService struct {
	loanRepo            repository.LoanRepository
	paymentRepo         repository.PaymentRepository
	customerRepo        repository.CustomerRepository
	notificationService service.NotificationService
	loyaltyService      service.LoyaltyService
	logger              zerolog.Logger
}

// NewJobService creates a new JobService
func NewJobService(
	loanRepo repository.LoanRepository,
	paymentRepo repository.PaymentRepository,
	customerRepo repository.CustomerRepository,
	notificationService service.NotificationService,
	loyaltyService service.LoyaltyService,
	logger zerolog.Logger,
) *JobService {
	return &JobService{
		loanRepo:            loanRepo,
		paymentRepo:         paymentRepo,
		customerRepo:        customerRepo,
		notificationService: notificationService,
		loyaltyService:      loyaltyService,
		logger:              logger,
	}
}

// ProcessOverdueLoans checks for overdue loans and updates their status
func (s *JobService) ProcessOverdueLoans(ctx context.Context) error {
	s.logger.Info().Msg("Processing overdue loans...")

	// Get all branches and process each
	// For simplicity, we'll process branch ID 0 which means all branches
	loans, err := s.loanRepo.GetOverdueLoans(ctx, 0)
	if err != nil {
		return err
	}

	now := time.Now()
	processed := 0
	defaulted := 0

	for _, loan := range loans {
		// Skip if already defaulted, paid, or confiscated
		if loan.Status == domain.LoanStatusDefaulted || loan.Status == domain.LoanStatusPaid || loan.Status == domain.LoanStatusConfiscated {
			continue
		}

		// Check if past due date
		if loan.DueDate.Before(now) {
			// If still active, mark as overdue
			if loan.Status == domain.LoanStatusActive {
				if err := s.loanRepo.UpdateStatus(ctx, loan.ID, domain.LoanStatusOverdue); err != nil {
					s.logger.Error().Err(err).Int64("loan_id", loan.ID).Msg("Failed to mark loan as overdue")
					continue
				}
				processed++
			}

			// Check if past grace period - should be defaulted
			gracePeriodEnd := loan.DueDate.AddDate(0, 0, loan.GracePeriodDays)
			if gracePeriodEnd.Before(now) && loan.Status == domain.LoanStatusOverdue {
				if err := s.loanRepo.UpdateStatus(ctx, loan.ID, domain.LoanStatusDefaulted); err != nil {
					s.logger.Error().Err(err).Int64("loan_id", loan.ID).Msg("Failed to mark loan as defaulted")
					continue
				}
				defaulted++
			}
		}
	}

	s.logger.Info().
		Int("marked_overdue", processed).
		Int("marked_defaulted", defaulted).
		Msg("Overdue loan processing completed")

	return nil
}

// CalculateLateFeesJob calculates late fees for overdue loans
func (s *JobService) CalculateLateFeesJob(ctx context.Context) error {
	s.logger.Info().Msg("Calculating late fees...")

	loans, err := s.loanRepo.GetOverdueLoans(ctx, 0)
	if err != nil {
		return err
	}

	now := time.Now()
	updated := 0

	for _, loan := range loans {
		// Skip if not overdue or already defaulted
		if loan.Status != domain.LoanStatusOverdue {
			continue
		}

		// Calculate days overdue
		daysOverdue := int(now.Sub(loan.DueDate).Hours() / 24)
		if daysOverdue <= 0 {
			continue
		}

		// Calculate late fee (daily rate * principal * days overdue)
		lateFee := loan.LateFeeRate / 100 * loan.LoanAmount * float64(daysOverdue)

		// Update loan if late fee changed significantly
		if lateFee > loan.LateFeeAmount+0.01 {
			loan.LateFeeAmount = lateFee
			if err := s.loanRepo.Update(ctx, loan); err != nil {
				s.logger.Error().Err(err).Int64("loan_id", loan.ID).Msg("Failed to update late fees")
				continue
			}
			updated++
		}
	}

	s.logger.Info().Int("updated", updated).Msg("Late fee calculation completed")
	return nil
}

// SendDueDateReminders sends reminders for loans approaching due date
func (s *JobService) SendDueDateReminders(ctx context.Context) error {
	s.logger.Info().Msg("Sending due date reminders...")

	// Get loans due in next 7 days
	params := repository.LoanListParams{
		Status: func() *domain.LoanStatus { status := domain.LoanStatusActive; return &status }(),
	}

	result, err := s.loanRepo.List(ctx, params)
	if err != nil {
		return err
	}

	now := time.Now()
	remindersSent := 0

	for _, loan := range result.Data {
		daysUntilDue := int(loan.DueDate.Sub(now).Hours() / 24)

		// Send reminder if due in 1, 3, or 7 days
		if daysUntilDue == 1 || daysUntilDue == 3 || daysUntilDue == 7 {
			// Get customer info
			customer, err := s.customerRepo.GetByID(ctx, loan.CustomerID)
			if err != nil || customer == nil {
				s.logger.Warn().Int64("customer_id", loan.CustomerID).Msg("Customer not found for reminder")
				continue
			}

			// Send notification using notification service
			if s.notificationService != nil {
				message := fmt.Sprintf("Su préstamo #%d vence en %d día(s). Monto pendiente: Q%.2f",
					loan.ID, daysUntilDue, loan.LoanAmount-loan.AmountPaid)

				_, err := s.notificationService.SendToCustomer(ctx, service.SendNotificationRequest{
					CustomerID:  loan.CustomerID,
					Type:        "loan_due_reminder",
					Title:       "Recordatorio de Vencimiento de Préstamo",
					Message:     message,
					Channel:     "sms", // Default to SMS
					ReferenceType: func() *string { t := "loan"; return &t }(),
					ReferenceID:   &loan.ID,
				})
				if err != nil {
					s.logger.Error().Err(err).Int64("loan_id", loan.ID).Msg("Failed to send reminder notification")
					continue
				}
			}

			s.logger.Info().
				Int64("loan_id", loan.ID).
				Int64("customer_id", loan.CustomerID).
				Int("days_until_due", daysUntilDue).
				Msg("Sent due date reminder")
			remindersSent++
		}
	}

	s.logger.Info().Int("reminders_sent", remindersSent).Msg("Due date reminder processing completed")
	return nil
}

// CalculateDailyInterest calculates daily interest for active loans
func (s *JobService) CalculateDailyInterest(ctx context.Context) error {
	s.logger.Info().Msg("Calculating daily interest...")

	params := repository.LoanListParams{
		Status: func() *domain.LoanStatus { status := domain.LoanStatusActive; return &status }(),
		PaginationParams: repository.PaginationParams{PerPage: 1000},
	}

	result, err := s.loanRepo.List(ctx, params)
	if err != nil {
		return err
	}

	updated := 0
	for i := range result.Data {
		loan := &result.Data[i]
		// Calculate daily interest rate
		dailyRate := loan.InterestRate / 100 / 365 // Annual rate to daily

		// Calculate interest for today
		dailyInterest := loan.LoanAmount * dailyRate

		// Add to accrued interest
		loan.InterestAmount += dailyInterest

		// Update the loan
		if err := s.loanRepo.Update(ctx, loan); err != nil {
			s.logger.Error().Err(err).Int64("loan_id", loan.ID).Msg("Failed to update loan interest")
			continue
		}
		updated++
	}

	s.logger.Info().
		Int("loans_processed", updated).
		Msg("Daily interest calculation completed")

	return nil
}

// SendOverdueNotifications sends notifications for overdue loans
func (s *JobService) SendOverdueNotifications(ctx context.Context) error {
	s.logger.Info().Msg("Sending overdue notifications...")

	loans, err := s.loanRepo.GetOverdueLoans(ctx, 0)
	if err != nil {
		return err
	}

	notificationsSent := 0
	for _, loan := range loans {
		if loan.Status != domain.LoanStatusOverdue {
			continue
		}

		// Get customer info
		customer, err := s.customerRepo.GetByID(ctx, loan.CustomerID)
		if err != nil || customer == nil {
			continue
		}

		// Calculate days overdue
		now := time.Now()
		daysOverdue := int(now.Sub(loan.DueDate).Hours() / 24)

		// Send notification using notification service
		if s.notificationService != nil {
			message := fmt.Sprintf("Su préstamo #%d está vencido por %d día(s). Por favor realice su pago lo antes posible para evitar cargos adicionales.",
				loan.ID, daysOverdue)

			_, err := s.notificationService.SendToCustomer(ctx, service.SendNotificationRequest{
				CustomerID:  loan.CustomerID,
				Type:        "loan_overdue",
				Title:       "Préstamo Vencido",
				Message:     message,
				Channel:     "sms",
				ReferenceType: func() *string { t := "loan"; return &t }(),
				ReferenceID:   &loan.ID,
			})
			if err != nil {
				s.logger.Error().Err(err).Int64("loan_id", loan.ID).Msg("Failed to send overdue notification")
				continue
			}
		}

		notificationsSent++
	}

	s.logger.Info().Int("notifications_sent", notificationsSent).Msg("Overdue notification processing completed")
	return nil
}

// CleanupExpiredSessions cleans up expired refresh tokens and sessions
func (s *JobService) CleanupExpiredSessions(ctx context.Context) error {
	s.logger.Info().Msg("Cleaning up expired sessions...")
	// TODO: Implement cleanup of expired refresh tokens
	// This would require RefreshTokenRepository.DeleteExpired()
	return nil
}

// GenerateDailyReport generates daily business reports
func (s *JobService) GenerateDailyReport(ctx context.Context) error {
	s.logger.Info().Msg("Generating daily report...")

	// Get yesterday's date range
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Second)

	dateFrom := startOfDay.Format("2006-01-02")
	dateTo := endOfDay.Format("2006-01-02")

	// Get loan statistics
	loanParams := repository.LoanListParams{
		PaginationParams: repository.PaginationParams{PerPage: 1000},
	}
	loans, _ := s.loanRepo.List(ctx, loanParams)

	// Get payment statistics
	paymentParams := repository.PaymentListParams{
		PaginationParams: repository.PaginationParams{PerPage: 1000},
		DateFrom:         &dateFrom,
		DateTo:           &dateTo,
	}
	payments, _ := s.paymentRepo.List(ctx, paymentParams)

	// Calculate totals
	var totalPayments float64
	for _, p := range payments.Data {
		if p.Status == domain.PaymentStatusCompleted {
			totalPayments += p.Amount
		}
	}

	s.logger.Info().
		Int("total_loans", loans.Total).
		Int("payments_today", len(payments.Data)).
		Float64("total_payments", totalPayments).
		Str("date", dateFrom).
		Msg("Daily report generated")

	return nil
}

// RegisterDefaultJobs registers all default scheduled jobs
func RegisterDefaultJobs(scheduler *Scheduler, jobService *JobService) {
	// Process overdue loans - run every hour
	scheduler.AddJob(&Job{
		Name:     "process_overdue_loans",
		Schedule: "hourly",
		Handler:  jobService.ProcessOverdueLoans,
		Enabled:  true,
	})

	// Calculate late fees - run every 6 hours
	scheduler.AddJob(&Job{
		Name:     "calculate_late_fees",
		Schedule: "every:6h",
		Handler:  jobService.CalculateLateFeesJob,
		Enabled:  true,
	})

	// Calculate daily interest - run every day at midnight
	scheduler.AddJob(&Job{
		Name:     "calculate_daily_interest",
		Schedule: "daily",
		Handler:  jobService.CalculateDailyInterest,
		Enabled:  true,
	})

	// Send due date reminders - run every day
	scheduler.AddJob(&Job{
		Name:     "send_due_date_reminders",
		Schedule: "daily",
		Handler:  jobService.SendDueDateReminders,
		Enabled:  true,
	})

	// Send overdue notifications - run every day
	scheduler.AddJob(&Job{
		Name:     "send_overdue_notifications",
		Schedule: "daily",
		Handler:  jobService.SendOverdueNotifications,
		Enabled:  true,
	})

	// Cleanup expired sessions - run every day
	scheduler.AddJob(&Job{
		Name:     "cleanup_expired_sessions",
		Schedule: "daily",
		Handler:  jobService.CleanupExpiredSessions,
		Enabled:  true,
	})

	// Generate daily report - run every day
	scheduler.AddJob(&Job{
		Name:     "generate_daily_report",
		Schedule: "daily",
		Handler:  jobService.GenerateDailyReport,
		Enabled:  true,
	})
}
