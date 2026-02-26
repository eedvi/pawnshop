package logger

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// BusinessLogger provides structured logging for business events
type BusinessLogger struct {
	logger zerolog.Logger
}

// NewBusinessLogger creates a new BusinessLogger
func NewBusinessLogger(logger zerolog.Logger) *BusinessLogger {
	return &BusinessLogger{logger: logger}
}

// getContextLogger creates a logger with context values
func (l *BusinessLogger) getContextLogger(ctx context.Context) *zerolog.Logger {
	logger := FromContext(ctx, l.logger)
	return &logger
}

// Loan Events

func (l *BusinessLogger) LoanCreated(ctx context.Context, loanID int64, customerID int64, amount float64, interestRate float64) {
	l.getContextLogger(ctx).Info().
		Int64("loan_id", loanID).
		Int64("customer_id", customerID).
		Float64("amount", amount).
		Float64("interest_rate", interestRate).
		Str("event_type", "loan_created").
		Msg("Loan created successfully")
}

func (l *BusinessLogger) LoanApproved(ctx context.Context, loanID int64, approvedBy int64) {
	l.getContextLogger(ctx).Info().
		Int64("loan_id", loanID).
		Int64("approved_by", approvedBy).
		Str("event_type", "loan_approved").
		Msg("Loan approved")
}

func (l *BusinessLogger) LoanRenewed(ctx context.Context, loanID int64, newDueDate string, extensionFee float64) {
	l.getContextLogger(ctx).Info().
		Int64("loan_id", loanID).
		Str("new_due_date", newDueDate).
		Float64("extension_fee", extensionFee).
		Str("event_type", "loan_renewed").
		Msg("Loan renewed")
}

func (l *BusinessLogger) LoanCompleted(ctx context.Context, loanID int64, totalPaid float64) {
	l.getContextLogger(ctx).Info().
		Int64("loan_id", loanID).
		Float64("total_paid", totalPaid).
		Str("event_type", "loan_completed").
		Msg("Loan completed and paid in full")
}

func (l *BusinessLogger) LoanDefaulted(ctx context.Context, loanID int64, daysOverdue int) {
	l.getContextLogger(ctx).Warn().
		Int64("loan_id", loanID).
		Int("days_overdue", daysOverdue).
		Str("event_type", "loan_defaulted").
		Msg("Loan defaulted")
}

// Payment Events

func (l *BusinessLogger) PaymentReceived(ctx context.Context, paymentID int64, loanID int64, amount float64, method string) {
	l.getContextLogger(ctx).Info().
		Int64("payment_id", paymentID).
		Int64("loan_id", loanID).
		Float64("amount", amount).
		Str("payment_method", method).
		Str("event_type", "payment_received").
		Msg("Payment received")
}

func (l *BusinessLogger) PaymentFailed(ctx context.Context, loanID int64, amount float64, reason string) {
	l.getContextLogger(ctx).Error().
		Int64("loan_id", loanID).
		Float64("amount", amount).
		Str("reason", reason).
		Str("event_type", "payment_failed").
		Msg("Payment processing failed")
}

// Sale Events

func (l *BusinessLogger) SaleCompleted(ctx context.Context, saleID int64, customerID int64, itemID int64, amount float64) {
	l.getContextLogger(ctx).Info().
		Int64("sale_id", saleID).
		Int64("customer_id", customerID).
		Int64("item_id", itemID).
		Float64("amount", amount).
		Str("event_type", "sale_completed").
		Msg("Sale completed successfully")
}

func (l *BusinessLogger) SaleRefunded(ctx context.Context, saleID int64, refundAmount float64, reason string) {
	l.getContextLogger(ctx).Info().
		Int64("sale_id", saleID).
		Float64("refund_amount", refundAmount).
		Str("reason", reason).
		Str("event_type", "sale_refunded").
		Msg("Sale refunded")
}

// Cash Session Events

func (l *BusinessLogger) CashSessionOpened(ctx context.Context, sessionID int64, registerID int64, openingAmount float64) {
	l.getContextLogger(ctx).Info().
		Int64("session_id", sessionID).
		Int64("register_id", registerID).
		Float64("opening_amount", openingAmount).
		Str("event_type", "cash_session_opened").
		Msg("Cash session opened")
}

func (l *BusinessLogger) CashSessionClosed(ctx context.Context, sessionID int64, closingAmount float64, difference float64) {
	l.getContextLogger(ctx).Info().
		Int64("session_id", sessionID).
		Float64("closing_amount", closingAmount).
		Float64("difference", difference).
		Str("event_type", "cash_session_closed").
		Msg("Cash session closed")
}

func (l *BusinessLogger) CashDiscrepancy(ctx context.Context, sessionID int64, expected float64, actual float64, difference float64) {
	l.getContextLogger(ctx).Warn().
		Int64("session_id", sessionID).
		Float64("expected", expected).
		Float64("actual", actual).
		Float64("difference", difference).
		Str("event_type", "cash_discrepancy").
		Msg("Cash discrepancy detected")
}

// Item Events

func (l *BusinessLogger) ItemReceived(ctx context.Context, itemID int64, category string, estimatedValue float64) {
	l.getContextLogger(ctx).Info().
		Int64("item_id", itemID).
		Str("category", category).
		Float64("estimated_value", estimatedValue).
		Str("event_type", "item_received").
		Msg("Item received for pawn")
}

func (l *BusinessLogger) ItemRedeemed(ctx context.Context, itemID int64, loanID int64) {
	l.getContextLogger(ctx).Info().
		Int64("item_id", itemID).
		Int64("loan_id", loanID).
		Str("event_type", "item_redeemed").
		Msg("Item redeemed by customer")
}

func (l *BusinessLogger) ItemSold(ctx context.Context, itemID int64, salePrice float64) {
	l.getContextLogger(ctx).Info().
		Int64("item_id", itemID).
		Float64("sale_price", salePrice).
		Str("event_type", "item_sold").
		Msg("Item sold")
}

// Authentication Events

func (l *BusinessLogger) UserLogin(ctx context.Context, userID int64, email string, ipAddress string) {
	l.getContextLogger(ctx).Info().
		Int64("user_id", userID).
		Str("email", email).
		Str("ip_address", ipAddress).
		Str("event_type", "user_login").
		Msg("User logged in successfully")
}

func (l *BusinessLogger) LoginFailed(ctx context.Context, email string, ipAddress string, reason string) {
	l.getContextLogger(ctx).Warn().
		Str("email", email).
		Str("ip_address", ipAddress).
		Str("reason", reason).
		Str("event_type", "login_failed").
		Msg("Login attempt failed")
}

func (l *BusinessLogger) UserLogout(ctx context.Context, userID int64) {
	l.getContextLogger(ctx).Info().
		Int64("user_id", userID).
		Str("event_type", "user_logout").
		Msg("User logged out")
}

// Performance Events

func (l *BusinessLogger) SlowQuery(ctx context.Context, query string, durationMS int64) {
	l.getContextLogger(ctx).Warn().
		Str("query", SanitizeSQL(query)).
		Int64("duration_ms", durationMS).
		Str("event_type", "slow_query").
		Msg("Slow database query detected")
}

func (l *BusinessLogger) HighMemoryUsage(ctx context.Context, memoryMB int64) {
	l.getContextLogger(ctx).Warn().
		Int64("memory_mb", memoryMB).
		Str("event_type", "high_memory_usage").
		Msg("High memory usage detected")
}

// Error Events

func (l *BusinessLogger) BusinessRuleViolation(ctx context.Context, rule string, details map[string]interface{}) {
	l.getContextLogger(ctx).Warn().
		Str("rule", rule).
		Interface("details", SanitizeMap(details)).
		Str("event_type", "business_rule_violation").
		Msg("Business rule violation detected")
}

func (l *BusinessLogger) ValidationError(ctx context.Context, field string, value interface{}, reason string) {
	l.getContextLogger(ctx).Warn().
		Str("field", field).
		Interface("value", value).
		Str("reason", reason).
		Str("event_type", "validation_error").
		Msg("Validation error")
}

// Helper function to log with default global logger
func LogBusinessEvent(ctx context.Context, eventType string, message string, fields map[string]interface{}) {
	logger := FromContext(ctx, log.Logger)
	event := logger.Info().Str("event_type", eventType)

	for key, value := range fields {
		switch v := value.(type) {
		case string:
			event.Str(key, v)
		case int:
			event.Int(key, v)
		case int64:
			event.Int64(key, v)
		case float64:
			event.Float64(key, v)
		case bool:
			event.Bool(key, v)
		default:
			event.Interface(key, v)
		}
	}

	event.Msg(message)
}
