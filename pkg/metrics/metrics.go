package metrics

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP Metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently in flight",
		},
	)

	// Database Metrics
	dbQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "table"},
	)

	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"operation", "table"},
	)

	dbConnectionsOpen = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_open",
			Help: "Number of open database connections",
		},
	)

	dbConnectionsInUse = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_in_use",
			Help: "Number of database connections in use",
		},
	)

	// Business Metrics
	loansCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loans_created_total",
			Help: "Total number of loans created",
		},
		[]string{"branch_id"},
	)

	loansAmount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loans_amount_total",
			Help: "Total amount of loans in currency",
		},
		[]string{"branch_id"},
	)

	paymentsReceived = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payments_received_total",
			Help: "Total number of payments received",
		},
		[]string{"branch_id", "method"},
	)

	paymentsAmount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payments_amount_total",
			Help: "Total amount of payments received",
		},
		[]string{"branch_id", "method"},
	)

	activeLoans = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "active_loans",
			Help: "Number of currently active loans",
		},
		[]string{"branch_id", "status"},
	)

	cashRegisterBalance = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cash_register_balance",
			Help: "Current cash register balance",
		},
		[]string{"branch_id", "register_id"},
	)

	// Worker Metrics
	jobExecutions = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "job_executions_total",
			Help: "Total number of job executions",
		},
		[]string{"job_name", "status"},
	)

	jobDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "job_duration_seconds",
			Help:    "Duration of job executions",
			Buckets: []float64{.1, .5, 1, 5, 10, 30, 60, 120, 300},
		},
		[]string{"job_name"},
	)

	// Notification Metrics
	notificationsSent = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notifications_sent_total",
			Help: "Total number of notifications sent",
		},
		[]string{"channel", "type", "status"},
	)
)

// Middleware returns a Fiber middleware for collecting HTTP metrics
func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		httpRequestsInFlight.Inc()

		// Process request
		err := c.Next()

		// Record metrics after request
		httpRequestsInFlight.Dec()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().StatusCode())
		path := c.Route().Path
		method := c.Method()

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}

// RecordDBQuery records database query metrics
func RecordDBQuery(operation, table string, duration time.Duration) {
	dbQueriesTotal.WithLabelValues(operation, table).Inc()
	dbQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// SetDBConnections sets database connection pool metrics
func SetDBConnections(open, inUse int) {
	dbConnectionsOpen.Set(float64(open))
	dbConnectionsInUse.Set(float64(inUse))
}

// RecordLoanCreated records loan creation metrics
func RecordLoanCreated(branchID int64, amount float64) {
	branch := strconv.FormatInt(branchID, 10)
	loansCreated.WithLabelValues(branch).Inc()
	loansAmount.WithLabelValues(branch).Add(amount)
}

// RecordPaymentReceived records payment metrics
func RecordPaymentReceived(branchID int64, method string, amount float64) {
	branch := strconv.FormatInt(branchID, 10)
	paymentsReceived.WithLabelValues(branch, method).Inc()
	paymentsAmount.WithLabelValues(branch, method).Add(amount)
}

// SetActiveLoans sets the gauge for active loans
func SetActiveLoans(branchID int64, status string, count int) {
	branch := strconv.FormatInt(branchID, 10)
	activeLoans.WithLabelValues(branch, status).Set(float64(count))
}

// SetCashRegisterBalance sets the cash register balance gauge
func SetCashRegisterBalance(branchID, registerID int64, balance float64) {
	branch := strconv.FormatInt(branchID, 10)
	register := strconv.FormatInt(registerID, 10)
	cashRegisterBalance.WithLabelValues(branch, register).Set(balance)
}

// RecordJobExecution records job execution metrics
func RecordJobExecution(jobName, status string, duration time.Duration) {
	jobExecutions.WithLabelValues(jobName, status).Inc()
	jobDuration.WithLabelValues(jobName).Observe(duration.Seconds())
}

// RecordNotificationSent records notification metrics
func RecordNotificationSent(channel, notificationType, status string) {
	notificationsSent.WithLabelValues(channel, notificationType, status).Inc()
}
