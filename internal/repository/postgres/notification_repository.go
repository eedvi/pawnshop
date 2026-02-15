package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// Notification Template Repository
type notificationTemplateRepository struct {
	db *DB
}

// NewNotificationTemplateRepository creates a new notification template repository
func NewNotificationTemplateRepository(db *DB) repository.NotificationTemplateRepository {
	return &notificationTemplateRepository{db: db}
}

func (r *notificationTemplateRepository) Create(ctx context.Context, template *domain.NotificationTemplate) error {
	query := `
		INSERT INTO notification_templates (notification_type, channel, name, subject, body_template, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		template.NotificationType,
		template.Channel,
		template.Name,
		template.Subject,
		template.BodyTemplate,
		template.IsActive,
	).Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt)
}

func (r *notificationTemplateRepository) GetByID(ctx context.Context, id int64) (*domain.NotificationTemplate, error) {
	query := `
		SELECT id, notification_type, channel, name, subject, body_template, is_active, created_at, updated_at
		FROM notification_templates
		WHERE id = $1`

	template := &domain.NotificationTemplate{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&template.ID,
		&template.NotificationType,
		&template.Channel,
		&template.Name,
		&template.Subject,
		&template.BodyTemplate,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return template, nil
}

func (r *notificationTemplateRepository) GetByTypeAndChannel(ctx context.Context, notificationType, channel string) (*domain.NotificationTemplate, error) {
	query := `
		SELECT id, notification_type, channel, name, subject, body_template, is_active, created_at, updated_at
		FROM notification_templates
		WHERE notification_type = $1 AND channel = $2 AND is_active = true`

	template := &domain.NotificationTemplate{}
	err := r.db.QueryRowContext(ctx, query, notificationType, channel).Scan(
		&template.ID,
		&template.NotificationType,
		&template.Channel,
		&template.Name,
		&template.Subject,
		&template.BodyTemplate,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return template, nil
}

func (r *notificationTemplateRepository) Update(ctx context.Context, template *domain.NotificationTemplate) error {
	query := `
		UPDATE notification_templates SET
			name = $2,
			subject = $3,
			body_template = $4,
			is_active = $5,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRowContext(ctx, query,
		template.ID,
		template.Name,
		template.Subject,
		template.BodyTemplate,
		template.IsActive,
	).Scan(&template.UpdatedAt)
}

func (r *notificationTemplateRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notification_templates WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *notificationTemplateRepository) List(ctx context.Context, includeInactive bool) ([]*domain.NotificationTemplate, error) {
	query := `
		SELECT id, notification_type, channel, name, subject, body_template, is_active, created_at, updated_at
		FROM notification_templates`

	if !includeInactive {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY notification_type, channel"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*domain.NotificationTemplate
	for rows.Next() {
		template := &domain.NotificationTemplate{}
		if err := rows.Scan(
			&template.ID,
			&template.NotificationType,
			&template.Channel,
			&template.Name,
			&template.Subject,
			&template.BodyTemplate,
			&template.IsActive,
			&template.CreatedAt,
			&template.UpdatedAt,
		); err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}

	return templates, rows.Err()
}

func (r *notificationTemplateRepository) ListByType(ctx context.Context, notificationType string) ([]*domain.NotificationTemplate, error) {
	query := `
		SELECT id, notification_type, channel, name, subject, body_template, is_active, created_at, updated_at
		FROM notification_templates
		WHERE notification_type = $1 AND is_active = true
		ORDER BY channel`

	rows, err := r.db.QueryContext(ctx, query, notificationType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*domain.NotificationTemplate
	for rows.Next() {
		template := &domain.NotificationTemplate{}
		if err := rows.Scan(
			&template.ID,
			&template.NotificationType,
			&template.Channel,
			&template.Name,
			&template.Subject,
			&template.BodyTemplate,
			&template.IsActive,
			&template.CreatedAt,
			&template.UpdatedAt,
		); err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}

	return templates, rows.Err()
}

// Notification Repository
type notificationRepository struct {
	db *DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *DB) repository.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	query := `
		INSERT INTO notifications (
			customer_id, branch_id, notification_type, channel,
			subject, body, reference_type, reference_id,
			status, scheduled_for
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		notification.CustomerID,
		notification.BranchID,
		notification.NotificationType,
		notification.Channel,
		notification.Subject,
		notification.Body,
		notification.ReferenceType,
		notification.ReferenceID,
		notification.Status,
		notification.ScheduledFor,
	).Scan(&notification.ID, &notification.CreatedAt, &notification.UpdatedAt)
}

func (r *notificationRepository) GetByID(ctx context.Context, id int64) (*domain.Notification, error) {
	query := `
		SELECT id, customer_id, branch_id, notification_type, channel,
			   subject, body, reference_type, reference_id,
			   status, scheduled_for, sent_at, delivered_at, failed_at,
			   failure_reason, retry_count, created_at, updated_at
		FROM notifications
		WHERE id = $1`

	notification := &domain.Notification{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID,
		&notification.CustomerID,
		&notification.BranchID,
		&notification.NotificationType,
		&notification.Channel,
		&notification.Subject,
		&notification.Body,
		&notification.ReferenceType,
		&notification.ReferenceID,
		&notification.Status,
		&notification.ScheduledFor,
		&notification.SentAt,
		&notification.DeliveredAt,
		&notification.FailedAt,
		&notification.FailureReason,
		&notification.RetryCount,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (r *notificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	query := `
		UPDATE notifications SET
			status = $2,
			scheduled_for = $3,
			sent_at = $4,
			delivered_at = $5,
			failed_at = $6,
			failure_reason = $7,
			retry_count = $8,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRowContext(ctx, query,
		notification.ID,
		notification.Status,
		notification.ScheduledFor,
		notification.SentAt,
		notification.DeliveredAt,
		notification.FailedAt,
		notification.FailureReason,
		notification.RetryCount,
	).Scan(&notification.UpdatedAt)
}

func (r *notificationRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notifications WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *notificationRepository) List(ctx context.Context, filter repository.NotificationFilter) ([]*domain.Notification, int64, error) {
	var conditions []string
	var args []interface{}
	argPos := 1

	if filter.CustomerID != nil {
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", argPos))
		args = append(args, *filter.CustomerID)
		argPos++
	}
	if filter.BranchID != nil {
		conditions = append(conditions, fmt.Sprintf("branch_id = $%d", argPos))
		args = append(args, *filter.BranchID)
		argPos++
	}
	if filter.NotificationType != nil {
		conditions = append(conditions, fmt.Sprintf("notification_type = $%d", argPos))
		args = append(args, *filter.NotificationType)
		argPos++
	}
	if filter.Channel != nil {
		conditions = append(conditions, fmt.Sprintf("channel = $%d", argPos))
		args = append(args, *filter.Channel)
		argPos++
	}
	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *filter.Status)
		argPos++
	}
	if filter.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argPos))
		args = append(args, *filter.DateFrom)
		argPos++
	}
	if filter.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argPos))
		args = append(args, *filter.DateTo)
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM notifications %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT id, customer_id, branch_id, notification_type, channel,
			   subject, body, reference_type, reference_id,
			   status, scheduled_for, sent_at, delivered_at, failed_at,
			   failure_reason, retry_count, created_at, updated_at
		FROM notifications
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argPos, argPos+1)

	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		notification := &domain.Notification{}
		if err := rows.Scan(
			&notification.ID,
			&notification.CustomerID,
			&notification.BranchID,
			&notification.NotificationType,
			&notification.Channel,
			&notification.Subject,
			&notification.Body,
			&notification.ReferenceType,
			&notification.ReferenceID,
			&notification.Status,
			&notification.ScheduledFor,
			&notification.SentAt,
			&notification.DeliveredAt,
			&notification.FailedAt,
			&notification.FailureReason,
			&notification.RetryCount,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, total, rows.Err()
}

func (r *notificationRepository) ListByCustomer(ctx context.Context, customerID int64, filter repository.NotificationFilter) ([]*domain.Notification, int64, error) {
	filter.CustomerID = &customerID
	return r.List(ctx, filter)
}

func (r *notificationRepository) ListPending(ctx context.Context, limit int) ([]*domain.Notification, error) {
	query := `
		SELECT id, customer_id, branch_id, notification_type, channel,
			   subject, body, reference_type, reference_id,
			   status, scheduled_for, sent_at, delivered_at, failed_at,
			   failure_reason, retry_count, created_at, updated_at
		FROM notifications
		WHERE status = 'pending' AND (scheduled_for IS NULL OR scheduled_for <= NOW())
		ORDER BY created_at ASC
		LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanNotifications(rows)
}

func (r *notificationRepository) ListScheduled(ctx context.Context, before time.Time, limit int) ([]*domain.Notification, error) {
	query := `
		SELECT id, customer_id, branch_id, notification_type, channel,
			   subject, body, reference_type, reference_id,
			   status, scheduled_for, sent_at, delivered_at, failed_at,
			   failure_reason, retry_count, created_at, updated_at
		FROM notifications
		WHERE status = 'pending' AND scheduled_for IS NOT NULL AND scheduled_for <= $1
		ORDER BY scheduled_for ASC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, before, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanNotifications(rows)
}

func (r *notificationRepository) ListFailed(ctx context.Context, maxRetries int, limit int) ([]*domain.Notification, error) {
	query := `
		SELECT id, customer_id, branch_id, notification_type, channel,
			   subject, body, reference_type, reference_id,
			   status, scheduled_for, sent_at, delivered_at, failed_at,
			   failure_reason, retry_count, created_at, updated_at
		FROM notifications
		WHERE status = 'failed' AND retry_count < $1
		ORDER BY failed_at ASC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, maxRetries, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanNotifications(rows)
}

func (r *notificationRepository) scanNotifications(rows *sql.Rows) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	for rows.Next() {
		notification := &domain.Notification{}
		if err := rows.Scan(
			&notification.ID,
			&notification.CustomerID,
			&notification.BranchID,
			&notification.NotificationType,
			&notification.Channel,
			&notification.Subject,
			&notification.Body,
			&notification.ReferenceType,
			&notification.ReferenceID,
			&notification.Status,
			&notification.ScheduledFor,
			&notification.SentAt,
			&notification.DeliveredAt,
			&notification.FailedAt,
			&notification.FailureReason,
			&notification.RetryCount,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, rows.Err()
}

func (r *notificationRepository) MarkAsSent(ctx context.Context, id int64) error {
	query := `
		UPDATE notifications SET
			status = 'sent',
			sent_at = NOW(),
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *notificationRepository) MarkAsDelivered(ctx context.Context, id int64) error {
	query := `
		UPDATE notifications SET
			status = 'delivered',
			delivered_at = NOW(),
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *notificationRepository) MarkAsFailed(ctx context.Context, id int64, reason string) error {
	query := `
		UPDATE notifications SET
			status = 'failed',
			failed_at = NOW(),
			failure_reason = $2,
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, reason)
	return err
}

func (r *notificationRepository) IncrementRetry(ctx context.Context, id int64) error {
	query := `
		UPDATE notifications SET
			retry_count = retry_count + 1,
			status = 'pending',
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *notificationRepository) Cancel(ctx context.Context, id int64) error {
	query := `
		UPDATE notifications SET
			status = 'cancelled',
			updated_at = NOW()
		WHERE id = $1 AND status = 'pending'`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *notificationRepository) GetStatsByCustomer(ctx context.Context, customerID int64) (*repository.NotificationStats, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE status = 'sent'),
			COUNT(*) FILTER (WHERE status = 'delivered'),
			COUNT(*) FILTER (WHERE status = 'failed'),
			COUNT(*) FILTER (WHERE status = 'pending')
		FROM notifications
		WHERE customer_id = $1`

	stats := &repository.NotificationStats{}
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(
		&stats.TotalSent,
		&stats.TotalDelivered,
		&stats.TotalFailed,
		&stats.TotalPending,
	)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *notificationRepository) GetStatsByBranch(ctx context.Context, branchID int64, dateFrom, dateTo time.Time) (*repository.NotificationStats, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE status = 'sent'),
			COUNT(*) FILTER (WHERE status = 'delivered'),
			COUNT(*) FILTER (WHERE status = 'failed'),
			COUNT(*) FILTER (WHERE status = 'pending')
		FROM notifications
		WHERE branch_id = $1 AND created_at >= $2 AND created_at <= $3`

	stats := &repository.NotificationStats{}
	err := r.db.QueryRowContext(ctx, query, branchID, dateFrom, dateTo).Scan(
		&stats.TotalSent,
		&stats.TotalDelivered,
		&stats.TotalFailed,
		&stats.TotalPending,
	)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// Customer Notification Preference Repository
type customerNotificationPreferenceRepository struct {
	db *DB
}

// NewCustomerNotificationPreferenceRepository creates a new customer notification preference repository
func NewCustomerNotificationPreferenceRepository(db *DB) repository.CustomerNotificationPreferenceRepository {
	return &customerNotificationPreferenceRepository{db: db}
}

func (r *customerNotificationPreferenceRepository) Create(ctx context.Context, pref *domain.CustomerNotificationPreference) error {
	query := `
		INSERT INTO customer_notification_preferences (customer_id, notification_type, channel, is_enabled)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		pref.CustomerID,
		pref.NotificationType,
		pref.Channel,
		pref.IsEnabled,
	).Scan(&pref.ID, &pref.CreatedAt, &pref.UpdatedAt)
}

func (r *customerNotificationPreferenceRepository) GetByID(ctx context.Context, id int64) (*domain.CustomerNotificationPreference, error) {
	query := `
		SELECT id, customer_id, notification_type, channel, is_enabled, created_at, updated_at
		FROM customer_notification_preferences
		WHERE id = $1`

	pref := &domain.CustomerNotificationPreference{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&pref.ID,
		&pref.CustomerID,
		&pref.NotificationType,
		&pref.Channel,
		&pref.IsEnabled,
		&pref.CreatedAt,
		&pref.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return pref, nil
}

func (r *customerNotificationPreferenceRepository) GetByCustomerTypeAndChannel(ctx context.Context, customerID int64, notificationType, channel string) (*domain.CustomerNotificationPreference, error) {
	query := `
		SELECT id, customer_id, notification_type, channel, is_enabled, created_at, updated_at
		FROM customer_notification_preferences
		WHERE customer_id = $1 AND notification_type = $2 AND channel = $3`

	pref := &domain.CustomerNotificationPreference{}
	err := r.db.QueryRowContext(ctx, query, customerID, notificationType, channel).Scan(
		&pref.ID,
		&pref.CustomerID,
		&pref.NotificationType,
		&pref.Channel,
		&pref.IsEnabled,
		&pref.CreatedAt,
		&pref.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return pref, nil
}

func (r *customerNotificationPreferenceRepository) Update(ctx context.Context, pref *domain.CustomerNotificationPreference) error {
	query := `
		UPDATE customer_notification_preferences SET
			is_enabled = $2,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRowContext(ctx, query, pref.ID, pref.IsEnabled).Scan(&pref.UpdatedAt)
}

func (r *customerNotificationPreferenceRepository) Upsert(ctx context.Context, pref *domain.CustomerNotificationPreference) error {
	query := `
		INSERT INTO customer_notification_preferences (customer_id, notification_type, channel, is_enabled)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (customer_id, notification_type, channel) DO UPDATE SET
			is_enabled = EXCLUDED.is_enabled,
			updated_at = NOW()
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		pref.CustomerID,
		pref.NotificationType,
		pref.Channel,
		pref.IsEnabled,
	).Scan(&pref.ID, &pref.CreatedAt, &pref.UpdatedAt)
}

func (r *customerNotificationPreferenceRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM customer_notification_preferences WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *customerNotificationPreferenceRepository) ListByCustomer(ctx context.Context, customerID int64) ([]*domain.CustomerNotificationPreference, error) {
	query := `
		SELECT id, customer_id, notification_type, channel, is_enabled, created_at, updated_at
		FROM customer_notification_preferences
		WHERE customer_id = $1
		ORDER BY notification_type, channel`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prefs []*domain.CustomerNotificationPreference
	for rows.Next() {
		pref := &domain.CustomerNotificationPreference{}
		if err := rows.Scan(
			&pref.ID,
			&pref.CustomerID,
			&pref.NotificationType,
			&pref.Channel,
			&pref.IsEnabled,
			&pref.CreatedAt,
			&pref.UpdatedAt,
		); err != nil {
			return nil, err
		}
		prefs = append(prefs, pref)
	}

	return prefs, rows.Err()
}

func (r *customerNotificationPreferenceRepository) IsEnabled(ctx context.Context, customerID int64, notificationType, channel string) (bool, error) {
	query := `
		SELECT is_enabled
		FROM customer_notification_preferences
		WHERE customer_id = $1 AND notification_type = $2 AND channel = $3`

	var isEnabled bool
	err := r.db.QueryRowContext(ctx, query, customerID, notificationType, channel).Scan(&isEnabled)
	if err == sql.ErrNoRows {
		// Default to enabled if no preference is set
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return isEnabled, nil
}

func (r *customerNotificationPreferenceRepository) BulkUpsert(ctx context.Context, customerID int64, prefs []*domain.CustomerNotificationPreference) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO customer_notification_preferences (customer_id, notification_type, channel, is_enabled)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (customer_id, notification_type, channel) DO UPDATE SET
			is_enabled = EXCLUDED.is_enabled,
			updated_at = NOW()
		RETURNING id, created_at, updated_at`

	for _, pref := range prefs {
		pref.CustomerID = customerID
		err = tx.QueryRowContext(ctx, query,
			pref.CustomerID,
			pref.NotificationType,
			pref.Channel,
			pref.IsEnabled,
		).Scan(&pref.ID, &pref.CreatedAt, &pref.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Internal Notification Repository
type internalNotificationRepository struct {
	db *DB
}

// NewInternalNotificationRepository creates a new internal notification repository
func NewInternalNotificationRepository(db *DB) repository.InternalNotificationRepository {
	return &internalNotificationRepository{db: db}
}

func (r *internalNotificationRepository) Create(ctx context.Context, notification *domain.InternalNotification) error {
	query := `
		INSERT INTO internal_notifications (
			user_id, branch_id, title, message, type,
			reference_type, reference_id, action_url
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query,
		notification.UserID,
		notification.BranchID,
		notification.Title,
		notification.Message,
		notification.Type,
		notification.ReferenceType,
		notification.ReferenceID,
		notification.ActionURL,
	).Scan(&notification.ID, &notification.CreatedAt)
}

func (r *internalNotificationRepository) GetByID(ctx context.Context, id int64) (*domain.InternalNotification, error) {
	query := `
		SELECT id, user_id, branch_id, title, message, type,
			   reference_type, reference_id, action_url, is_read, read_at, created_at
		FROM internal_notifications
		WHERE id = $1`

	notification := &domain.InternalNotification{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.BranchID,
		&notification.Title,
		&notification.Message,
		&notification.Type,
		&notification.ReferenceType,
		&notification.ReferenceID,
		&notification.ActionURL,
		&notification.IsRead,
		&notification.ReadAt,
		&notification.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (r *internalNotificationRepository) Update(ctx context.Context, notification *domain.InternalNotification) error {
	query := `
		UPDATE internal_notifications SET
			is_read = $2,
			read_at = $3
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, notification.ID, notification.IsRead, notification.ReadAt)
	return err
}

func (r *internalNotificationRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM internal_notifications WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *internalNotificationRepository) List(ctx context.Context, filter repository.InternalNotificationFilter) ([]*domain.InternalNotification, int64, error) {
	var conditions []string
	var args []interface{}
	argPos := 1

	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argPos))
		args = append(args, *filter.UserID)
		argPos++
	}
	if filter.BranchID != nil {
		conditions = append(conditions, fmt.Sprintf("branch_id = $%d", argPos))
		args = append(args, *filter.BranchID)
		argPos++
	}
	if filter.Type != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argPos))
		args = append(args, *filter.Type)
		argPos++
	}
	if filter.IsRead != nil {
		conditions = append(conditions, fmt.Sprintf("is_read = $%d", argPos))
		args = append(args, *filter.IsRead)
		argPos++
	}
	if filter.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argPos))
		args = append(args, *filter.DateFrom)
		argPos++
	}
	if filter.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argPos))
		args = append(args, *filter.DateTo)
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM internal_notifications %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT id, user_id, branch_id, title, message, type,
			   reference_type, reference_id, action_url, is_read, read_at, created_at
		FROM internal_notifications
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argPos, argPos+1)

	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notifications, err := r.scanInternalNotifications(rows)
	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *internalNotificationRepository) ListByUser(ctx context.Context, userID int64, filter repository.InternalNotificationFilter) ([]*domain.InternalNotification, int64, error) {
	filter.UserID = &userID
	return r.List(ctx, filter)
}

func (r *internalNotificationRepository) ListUnreadByUser(ctx context.Context, userID int64, limit int) ([]*domain.InternalNotification, error) {
	query := `
		SELECT id, user_id, branch_id, title, message, type,
			   reference_type, reference_id, action_url, is_read, read_at, created_at
		FROM internal_notifications
		WHERE user_id = $1 AND is_read = false
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanInternalNotifications(rows)
}

func (r *internalNotificationRepository) scanInternalNotifications(rows *sql.Rows) ([]*domain.InternalNotification, error) {
	var notifications []*domain.InternalNotification
	for rows.Next() {
		notification := &domain.InternalNotification{}
		if err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.BranchID,
			&notification.Title,
			&notification.Message,
			&notification.Type,
			&notification.ReferenceType,
			&notification.ReferenceID,
			&notification.ActionURL,
			&notification.IsRead,
			&notification.ReadAt,
			&notification.CreatedAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, rows.Err()
}

func (r *internalNotificationRepository) MarkAsRead(ctx context.Context, id int64) error {
	query := `
		UPDATE internal_notifications SET
			is_read = true,
			read_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *internalNotificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	query := `
		UPDATE internal_notifications SET
			is_read = true,
			read_at = NOW()
		WHERE user_id = $1 AND is_read = false`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *internalNotificationRepository) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM internal_notifications WHERE user_id = $1 AND is_read = false`

	var count int64
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *internalNotificationRepository) CreateBulk(ctx context.Context, notifications []*domain.InternalNotification) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO internal_notifications (
			user_id, branch_id, title, message, type,
			reference_type, reference_id, action_url
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	for _, notification := range notifications {
		err = tx.QueryRowContext(ctx, query,
			notification.UserID,
			notification.BranchID,
			notification.Title,
			notification.Message,
			notification.Type,
			notification.ReferenceType,
			notification.ReferenceID,
			notification.ActionURL,
		).Scan(&notification.ID, &notification.CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *internalNotificationRepository) DeleteOlderThan(ctx context.Context, olderThan time.Time) (int64, error) {
	query := `DELETE FROM internal_notifications WHERE created_at < $1`

	result, err := r.db.ExecContext(ctx, query, olderThan)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
