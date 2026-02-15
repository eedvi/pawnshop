-- Drop notification tables
DROP TABLE IF EXISTS internal_notifications;
DROP TABLE IF EXISTS customer_notification_preferences;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS notification_templates;

-- Drop types
DROP TYPE IF EXISTS notification_status;
DROP TYPE IF EXISTS notification_channel;
DROP TYPE IF EXISTS notification_type;
