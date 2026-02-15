package domain

import "time"

// Transfer status constants
const (
	TransferStatusPending   = "pending"
	TransferStatusInTransit = "in_transit"
	TransferStatusCompleted = "completed"
	TransferStatusCancelled = "cancelled"
)

// ItemTransfer represents a transfer of an item between branches
type ItemTransfer struct {
	ID             int64  `json:"id"`
	TransferNumber string `json:"transfer_number"`

	// Item being transferred
	ItemID int64 `json:"item_id"`
	Item   *Item `json:"item,omitempty"`

	// Branches
	FromBranchID int64   `json:"from_branch_id"`
	FromBranch   *Branch `json:"from_branch,omitempty"`
	ToBranchID   int64   `json:"to_branch_id"`
	ToBranch     *Branch `json:"to_branch,omitempty"`

	// Status
	Status string `json:"status"`

	// Users involved
	RequestedBy int64 `json:"requested_by"`
	Requester   *User `json:"requester,omitempty"`
	ApprovedBy  *int64 `json:"approved_by,omitempty"`
	Approver    *User  `json:"approver,omitempty"`
	ReceivedBy  *int64 `json:"received_by,omitempty"`
	Receiver    *User  `json:"receiver,omitempty"`

	// Dates
	RequestedAt  time.Time  `json:"requested_at"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	ShippedAt    *time.Time `json:"shipped_at,omitempty"`
	ReceivedAt   *time.Time `json:"received_at,omitempty"`
	CancelledAt  *time.Time `json:"cancelled_at,omitempty"`

	// Notes
	RequestNotes        string  `json:"request_notes,omitempty"`
	ApprovalNotes       string  `json:"approval_notes,omitempty"`
	ReceiptNotes        string  `json:"receipt_notes,omitempty"`
	CancellationReason  string  `json:"cancellation_reason,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsPending checks if transfer is pending approval
func (t *ItemTransfer) IsPending() bool {
	return t.Status == TransferStatusPending
}

// IsInTransit checks if transfer is in transit
func (t *ItemTransfer) IsInTransit() bool {
	return t.Status == TransferStatusInTransit
}

// IsCompleted checks if transfer is completed
func (t *ItemTransfer) IsCompleted() bool {
	return t.Status == TransferStatusCompleted
}

// IsCancelled checks if transfer is cancelled
func (t *ItemTransfer) IsCancelled() bool {
	return t.Status == TransferStatusCancelled
}

// CanApprove checks if transfer can be approved
func (t *ItemTransfer) CanApprove() bool {
	return t.Status == TransferStatusPending
}

// CanShip checks if transfer can be shipped
func (t *ItemTransfer) CanShip() bool {
	return t.Status == TransferStatusPending
}

// CanReceive checks if transfer can be received
func (t *ItemTransfer) CanReceive() bool {
	return t.Status == TransferStatusInTransit
}

// CanCancel checks if transfer can be cancelled
func (t *ItemTransfer) CanCancel() bool {
	return t.Status == TransferStatusPending || t.Status == TransferStatusInTransit
}
