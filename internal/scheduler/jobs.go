package scheduler

import (
	"context"
	"fmt"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/internal/service"

	"github.com/rs/zerolog"
)

// JobService contains dependencies for scheduled jobs
type JobService struct {
	loanRepo            repository.LoanRepository
	itemRepo            repository.ItemRepository
	paymentRepo         repository.PaymentRepository
	customerRepo        repository.CustomerRepository
	notificationService service.NotificationService
	loyaltyService      service.LoyaltyService
	logger              zerolog.Logger
}

// NewJobService creates a new JobService
func NewJobService(
	loanRepo repository.LoanRepository,
	itemRepo repository.ItemRepository,
	paymentRepo repository.PaymentRepository,
	customerRepo repository.CustomerRepository,
	notificationService service.NotificationService,
	loyaltyService service.LoyaltyService,
	logger zerolog.Logger,
) *JobService {
	return &JobService{
		loanRepo:            loanRepo,
		itemRepo:            itemRepo,
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

	s.logger.Debug().Int("loans_found", len(loans)).Msg("Overdue loans retrieved")

	now := time.Now()
	processed := 0
	confiscated := 0
	skipped := 0

	for _, loan := range loans {
		s.logger.Debug().
			Int64("loan_id", loan.ID).
			Str("loan_number", loan.LoanNumber).
			Str("status", string(loan.Status)).
			Time("due_date", loan.DueDate.Time).
			Int("grace_period_days", loan.GracePeriodDays).
			Msg("Processing loan for status update")

		// Skip if already paid or confiscated
		if loan.Status == domain.LoanStatusPaid || loan.Status == domain.LoanStatusConfiscated {
			s.logger.Debug().
				Int64("loan_id", loan.ID).
				Str("status", string(loan.Status)).
				Msg("Skipping - already paid or confiscated")
			skipped++
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
				s.logger.Info().
					Int64("loan_id", loan.ID).
					Str("loan_number", loan.LoanNumber).
					Msg("Loan marked as overdue")
				processed++
			}

			// Check if past grace period - automatic confiscation
			// Add grace period days and set to end of day (23:59:59) to avoid timezone issues
			gracePeriodEnd := loan.DueDate.AddDate(0, 0, loan.GracePeriodDays)
			gracePeriodEndOfDay := time.Date(
				gracePeriodEnd.Year(),
				gracePeriodEnd.Month(),
				gracePeriodEnd.Day(),
				23, 59, 59, 0,
				time.Local, // Use local timezone
			)
			daysUntilConfiscation := int(gracePeriodEndOfDay.Sub(now).Hours() / 24)

			s.logger.Debug().
				Int64("loan_id", loan.ID).
				Time("grace_period_end", gracePeriodEndOfDay).
				Int("days_until_confiscation", daysUntilConfiscation).
				Bool("past_grace_period", gracePeriodEndOfDay.Before(now)).
				Str("current_status", string(loan.Status)).
				Msg("Checking confiscation eligibility")

			if gracePeriodEndOfDay.Before(now) && loan.Status == domain.LoanStatusOverdue {
				// Update loan status to confiscated
				if err := s.loanRepo.UpdateStatus(ctx, loan.ID, domain.LoanStatusConfiscated); err != nil {
					s.logger.Error().Err(err).Int64("loan_id", loan.ID).Msg("Failed to confiscate loan")
					continue
				}

				// Update item status to for_sale
				if err := s.itemRepo.UpdateStatus(ctx, loan.ItemID, domain.ItemStatusForSale); err != nil {
					s.logger.Error().Err(err).Int64("loan_id", loan.ID).Int64("item_id", loan.ItemID).Msg("Failed to update item status to for_sale")
					// Continue anyway - loan status was updated
				}

				s.logger.Info().
					Int64("loan_id", loan.ID).
					Str("loan_number", loan.LoanNumber).
					Int64("item_id", loan.ItemID).
					Msg("Loan automatically confiscated after grace period")

				confiscated++
			} else {
				s.logger.Debug().
					Int64("loan_id", loan.ID).
					Msg("Not yet eligible for confiscation")
			}
		} else {
			s.logger.Debug().
				Int64("loan_id", loan.ID).
				Msg("Skipping - not past due date yet")
			skipped++
		}
	}

	s.logger.Info().
		Int("marked_overdue", processed).
		Int("auto_confiscated", confiscated).
		Int("skipped", skipped).
		Int("total_processed", len(loans)).
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

	s.logger.Debug().Int("loans_found", len(loans)).Msg("Overdue loans retrieved")

	now := time.Now()
	updated := 0
	skipped := 0

	for _, loan := range loans {
		s.logger.Debug().
			Int64("loan_id", loan.ID).
			Str("loan_number", loan.LoanNumber).
			Str("status", string(loan.Status)).
			Time("due_date", loan.DueDate.Time).
			Float64("current_late_fee", loan.LateFeeAmount).
			Float64("late_fee_rate", loan.LateFeeRate).
			Float64("loan_amount", loan.LoanAmount).
			Msg("Processing loan for late fees")

		// Only calculate late fees for overdue loans (during grace period)
		// Once confiscated, late fees are frozen
		if loan.Status != domain.LoanStatusOverdue {
			s.logger.Debug().
				Int64("loan_id", loan.ID).
				Str("status", string(loan.Status)).
				Msg("Skipping - not in overdue status")
			skipped++
			continue
		}

		// Calculate days overdue
		daysOverdue := int(now.Sub(loan.DueDate.Time).Hours() / 24)
		if daysOverdue <= 0 {
			s.logger.Debug().
				Int64("loan_id", loan.ID).
				Int("days_overdue", daysOverdue).
				Msg("Skipping - not yet overdue")
			skipped++
			continue
		}

		// Calculate late fee (daily rate * principal * days overdue)
		lateFee := loan.LateFeeRate / 100 * loan.LoanAmount * float64(daysOverdue)

		s.logger.Debug().
			Int64("loan_id", loan.ID).
			Int("days_overdue", daysOverdue).
			Float64("calculated_late_fee", lateFee).
			Float64("current_late_fee", loan.LateFeeAmount).
			Float64("current_late_fee_remaining", loan.LateFeeRemaining).
			Float64("difference", lateFee-loan.LateFeeAmount).
			Msg("Late fee calculation")

		// Update loan if late fee changed significantly
		if lateFee > loan.LateFeeAmount+0.01 {
			// Calculate the increment to add to both fields
			increment := lateFee - loan.LateFeeAmount
			loan.LateFeeAmount = lateFee           // Total historical amount
			loan.LateFeeRemaining += increment     // Add increment to remaining (what's still owed)
			if err := s.loanRepo.Update(ctx, loan); err != nil {
				s.logger.Error().Err(err).Int64("loan_id", loan.ID).Msg("Failed to update late fees")
				continue
			}
			s.logger.Info().
				Int64("loan_id", loan.ID).
				Str("loan_number", loan.LoanNumber).
				Float64("new_late_fee", lateFee).
				Int("days_overdue", daysOverdue).
				Msg("Late fee updated")
			updated++
		} else {
			s.logger.Debug().
				Int64("loan_id", loan.ID).
				Msg("Skipping - late fee unchanged")
			skipped++
		}
	}

	s.logger.Info().
		Int("updated", updated).
		Int("skipped", skipped).
		Int("total_processed", len(loans)).
		Msg("Late fee calculation completed")
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
					CustomerID:    loan.CustomerID,
					Type:          "loan_due_reminder",
					Title:         "Recordatorio de Vencimiento de Préstamo",
					Message:       message,
					Channel:       "sms", // Default to SMS
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
		Status:           func() *domain.LoanStatus { status := domain.LoanStatusActive; return &status }(),
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
	now := time.Now()

	for _, loan := range loans {
		if loan.Status != domain.LoanStatusOverdue {
			continue
		}

		// Get customer info
		customer, err := s.customerRepo.GetByID(ctx, loan.CustomerID)
		if err != nil || customer == nil {
			continue
		}

		// Calculate days overdue and remaining grace period
		daysOverdue := int(now.Sub(loan.DueDate.Time).Hours() / 24)
		gracePeriodEnd := loan.DueDate.AddDate(0, 0, loan.GracePeriodDays)
		daysUntilConfiscation := int(gracePeriodEnd.Sub(now).Hours() / 24)

		// Send notification using notification service
		if s.notificationService != nil {
			var message string
			if daysUntilConfiscation > 0 {
				message = fmt.Sprintf("⚠️ Su préstamo %s está vencido por %d día(s). Le quedan %d día(s) para pagar antes de que el artículo sea confiscado y puesto en venta. Monto pendiente: Q%.2f + Mora: Q%.2f",
					loan.LoanNumber, daysOverdue, daysUntilConfiscation, loan.PrincipalRemaining+loan.InterestRemaining, loan.LateFeeAmount)
			} else {
				message = fmt.Sprintf("🚨 URGENTE: Su préstamo %s será confiscado HOY. Pague ahora para evitar la pérdida del artículo. Monto pendiente: Q%.2f + Mora: Q%.2f",
					loan.LoanNumber, loan.PrincipalRemaining+loan.InterestRemaining, loan.LateFeeAmount)
			}

			_, err := s.notificationService.SendToCustomer(ctx, service.SendNotificationRequest{
				CustomerID:    loan.CustomerID,
				Type:          "loan_overdue",
				Title:         "Préstamo Vencido - Riesgo de Confiscación",
				Message:       message,
				Channel:       "sms",
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
	// Process overdue loans - run every minute (dev: every:1m, prod: hourly)
	scheduler.AddJob(&Job{
		Name:     "process_overdue_loans",
		Schedule: "every:1m",
		Handler:  jobService.ProcessOverdueLoans,
		Enabled:  true,
	})

	// Calculate late fees - run every minute (dev: every:1m, prod: every:6h)
	scheduler.AddJob(&Job{
		Name:     "calculate_late_fees",
		Schedule: "every:1m",
		Handler:  jobService.CalculateLateFeesJob,
		Enabled:  true,
	})

	// Calculate daily interest - DISABLED for pawnshop model (simple interest)
	// In pawnshops, interest is calculated once at loan creation, not daily
	// Only late fees (mora) accumulate when overdue
	scheduler.AddJob(&Job{
		Name:     "calculate_daily_interest",
		Schedule: "every:1m",
		Handler:  jobService.CalculateDailyInterest,
		Enabled:  false, // Disabled: pawnshop uses simple interest model
	})

	// Send due date reminders - run every day
	scheduler.AddJob(&Job{
		Name:     "send_due_date_reminders",
		Schedule: "every:1m",
		Handler:  jobService.SendDueDateReminders,
		Enabled:  true,
	})

	// Send overdue notifications - run every day
	scheduler.AddJob(&Job{
		Name:     "send_overdue_notifications",
		Schedule: "every:1m",
		Handler:  jobService.SendOverdueNotifications,
		Enabled:  true,
	})

	// Cleanup expired sessions - run every day
	scheduler.AddJob(&Job{
		Name:     "cleanup_expired_sessions",
		Schedule: "every:1m",
		Handler:  jobService.CleanupExpiredSessions,
		Enabled:  true,
	})

	// Generate daily report - run every day
	scheduler.AddJob(&Job{
		Name:     "generate_daily_report",
		Schedule: "every:1m",
		Handler:  jobService.GenerateDailyReport,
		Enabled:  true,
	})
}
