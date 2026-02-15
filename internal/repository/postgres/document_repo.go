package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pawnshop/internal/domain"
)

// DocumentRepository implements repository.DocumentRepository
type DocumentRepository struct {
	db *DB
}

// NewDocumentRepository creates a new DocumentRepository
func NewDocumentRepository(db *DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// Create creates a new document
func (r *DocumentRepository) Create(ctx context.Context, doc *domain.Document) error {
	query := `
		INSERT INTO documents (branch_id, document_type, document_number, reference_type, reference_id, file_path, file_url, file_size, mime_type, content_hash, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query,
		doc.BranchID,
		doc.DocumentType,
		doc.DocumentNumber,
		doc.ReferenceType,
		doc.ReferenceID,
		NullString(doc.FilePath),
		NullString(doc.FileURL),
		doc.FileSize,
		doc.MimeType,
		NullString(doc.ContentHash),
		doc.CreatedBy,
	).Scan(&doc.ID, &doc.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	return nil
}

// GetByID retrieves a document by ID
func (r *DocumentRepository) GetByID(ctx context.Context, id int64) (*domain.Document, error) {
	query := `
		SELECT id, branch_id, document_type, document_number, reference_type, reference_id, file_path, file_url, file_size, mime_type, content_hash, created_by, created_at
		FROM documents
		WHERE id = $1
	`

	doc := &domain.Document{}
	var filePath, fileURL, contentHash sql.NullString
	var fileSize sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&doc.ID,
		&doc.BranchID,
		&doc.DocumentType,
		&doc.DocumentNumber,
		&doc.ReferenceType,
		&doc.ReferenceID,
		&filePath,
		&fileURL,
		&fileSize,
		&doc.MimeType,
		&contentHash,
		&doc.CreatedBy,
		&doc.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	doc.FilePath = StringPtr(filePath)
	doc.FileURL = StringPtr(fileURL)
	doc.ContentHash = StringPtr(contentHash)
	if fileSize.Valid {
		doc.FileSize = int(fileSize.Int64)
	}

	return doc, nil
}

// ListByReference retrieves documents by reference type and ID
func (r *DocumentRepository) ListByReference(ctx context.Context, refType string, refID int64) ([]*domain.Document, error) {
	query := `
		SELECT id, branch_id, document_type, document_number, reference_type, reference_id, file_path, file_url, file_size, mime_type, content_hash, created_by, created_at
		FROM documents
		WHERE reference_type = $1 AND reference_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, refType, refID)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	defer rows.Close()

	docs := []*domain.Document{}
	for rows.Next() {
		doc := &domain.Document{}
		var filePath, fileURL, contentHash sql.NullString
		var fileSize sql.NullInt64

		err := rows.Scan(
			&doc.ID,
			&doc.BranchID,
			&doc.DocumentType,
			&doc.DocumentNumber,
			&doc.ReferenceType,
			&doc.ReferenceID,
			&filePath,
			&fileURL,
			&fileSize,
			&doc.MimeType,
			&contentHash,
			&doc.CreatedBy,
			&doc.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}

		doc.FilePath = StringPtr(filePath)
		doc.FileURL = StringPtr(fileURL)
		doc.ContentHash = StringPtr(contentHash)
		if fileSize.Valid {
			doc.FileSize = int(fileSize.Int64)
		}

		docs = append(docs, doc)
	}

	return docs, nil
}
