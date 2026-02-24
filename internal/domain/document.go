package domain

import (
	"time"
)

// DocumentType represents the type of document
type DocumentType string

const (
	DocumentTypeLoanContract       DocumentType = "loan_contract"
	DocumentTypeLoanReceipt        DocumentType = "loan_receipt"
	DocumentTypePaymentReceipt     DocumentType = "payment_receipt"
	DocumentTypeSaleReceipt        DocumentType = "sale_receipt"
	DocumentTypeConfiscationNotice DocumentType = "confiscation_notice"
	DocumentTypeOther              DocumentType = "other"
)

// Document represents a generated document
type Document struct {
	ID             int64        `json:"id"`
	BranchID       int64        `json:"branch_id"`
	DocumentType   DocumentType `json:"document_type"`
	DocumentNumber string       `json:"document_number"`

	// Reference to related entity
	ReferenceType string `json:"reference_type"` // loan, payment, sale
	ReferenceID   int64  `json:"reference_id"`

	// File info
	FilePath  string `json:"file_path,omitempty"`
	FileURL   string `json:"file_url,omitempty"`
	FileSize  int    `json:"file_size,omitempty"`
	MimeType  string `json:"mime_type"`

	// Content hash for integrity
	ContentHash string `json:"content_hash,omitempty"`

	// Audit
	CreatedBy int64     `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Branch *Branch `json:"branch,omitempty"`
}

// TableName returns the database table name
func (Document) TableName() string {
	return "documents"
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID       int64  `json:"id"`
	BranchID *int64 `json:"branch_id,omitempty"`
	UserID   *int64 `json:"user_id,omitempty"`

	// Action details
	Action      string  `json:"action"`
	EntityType  string  `json:"entity_type"`
	EntityID    *int64  `json:"entity_id,omitempty"`
	Description *string `json:"description,omitempty"`

	// Data
	OldValues interface{} `json:"old_values,omitempty"`
	NewValues interface{} `json:"new_values,omitempty"`
	IPAddress string      `json:"ip_address,omitempty"`
	UserAgent string      `json:"user_agent,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`

	// Relations (populated from JOINs, not from foreign keys)
	User       *User   `json:"user,omitempty"`
	UserName   *string `json:"user_name,omitempty"`
	BranchName *string `json:"branch_name,omitempty"`
}

// TableName returns the database table name
func (AuditLog) TableName() string {
	return "audit_logs"
}

// Setting represents a system setting
type Setting struct {
	ID          int64       `json:"id"`
	BranchID    *int64      `json:"branch_id,omitempty"` // NULL for global settings
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Description string      `json:"description,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName returns the database table name
func (Setting) TableName() string {
	return "settings"
}

// RefreshToken represents a refresh token for auth
type RefreshToken struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	TokenHash  string     `json:"-"` // Never expose token hash
	DeviceInfo string     `json:"device_info,omitempty"`
	IPAddress  string     `json:"ip_address,omitempty"`
	ExpiresAt  time.Time  `json:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// TableName returns the database table name
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsExpired checks if the token is expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsRevoked checks if the token is revoked
func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

// IsValid checks if the token is valid
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.IsRevoked()
}
