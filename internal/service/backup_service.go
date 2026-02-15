package service

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"pawnshop/internal/config"
)

// BackupInfo contains information about a backup
type BackupInfo struct {
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	Compressed  bool      `json:"compressed"`
	Description string    `json:"description,omitempty"`
}

// BackupService defines the interface for backup operations
type BackupService interface {
	// CreateBackup creates a new database backup
	CreateBackup(ctx context.Context, description string) (*BackupInfo, error)

	// RestoreBackup restores a database from a backup file
	RestoreBackup(ctx context.Context, filename string) error

	// ListBackups lists all available backups
	ListBackups(ctx context.Context) ([]*BackupInfo, error)

	// DeleteBackup deletes a backup file
	DeleteBackup(ctx context.Context, filename string) error

	// GetBackup retrieves a backup file for download
	GetBackup(ctx context.Context, filename string) (io.ReadCloser, *BackupInfo, error)

	// ScheduledBackup performs a scheduled backup (for cron jobs)
	ScheduledBackup(ctx context.Context) (*BackupInfo, error)

	// CleanupOldBackups removes backups older than the retention period
	CleanupOldBackups(ctx context.Context, retentionDays int) (int, error)
}

type backupService struct {
	dbConfig  *config.DatabaseConfig
	backupDir string
}

// NewBackupService creates a new backup service
func NewBackupService(dbConfig *config.DatabaseConfig, backupDir string) BackupService {
	// Ensure backup directory exists
	os.MkdirAll(backupDir, 0755)

	return &backupService{
		dbConfig:  dbConfig,
		backupDir: backupDir,
	}
}

func (s *backupService) CreateBackup(ctx context.Context, description string) (*BackupInfo, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("pawnshop_backup_%s.sql.gz", timestamp)
	filepath := filepath.Join(s.backupDir, filename)

	// Create pg_dump command
	cmd := exec.CommandContext(ctx, "pg_dump",
		"-h", s.dbConfig.Host,
		"-p", fmt.Sprintf("%d", s.dbConfig.Port),
		"-U", s.dbConfig.User,
		"-d", s.dbConfig.DBName,
		"-F", "p", // Plain text format
		"--no-owner",
		"--no-privileges",
	)

	// Set password via environment variable
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", s.dbConfig.Password))

	// Create output file with gzip compression
	outFile, err := os.Create(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup file: %w", err)
	}
	defer outFile.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	// Set comment in gzip header
	gzWriter.Comment = description
	gzWriter.ModTime = time.Now()

	// Pipe pg_dump output to gzip
	cmd.Stdout = gzWriter

	// Capture stderr
	var stderr strings.Builder
	cmd.Stderr = &stderr

	// Run the backup
	if err := cmd.Run(); err != nil {
		os.Remove(filepath) // Clean up partial file
		return nil, fmt.Errorf("pg_dump failed: %s - %w", stderr.String(), err)
	}

	// Ensure gzip is flushed
	if err := gzWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// Get file info
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup file info: %w", err)
	}

	return &BackupInfo{
		Filename:    filename,
		Size:        fileInfo.Size(),
		CreatedAt:   time.Now(),
		Compressed:  true,
		Description: description,
	}, nil
}

func (s *backupService) RestoreBackup(ctx context.Context, filename string) error {
	filepath := filepath.Join(s.backupDir, filename)

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", filename)
	}

	// Open the backup file
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer file.Close()

	// Create gzip reader if file is compressed
	var reader io.Reader = file
	if strings.HasSuffix(filename, ".gz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	// Create psql command for restore
	cmd := exec.CommandContext(ctx, "psql",
		"-h", s.dbConfig.Host,
		"-p", fmt.Sprintf("%d", s.dbConfig.Port),
		"-U", s.dbConfig.User,
		"-d", s.dbConfig.DBName,
		"-v", "ON_ERROR_STOP=1",
	)

	// Set password via environment variable
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", s.dbConfig.Password))

	// Pipe backup content to psql
	cmd.Stdin = reader

	// Capture output
	var stderr strings.Builder
	cmd.Stderr = &stderr

	// Run the restore
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("restore failed: %s - %w", stderr.String(), err)
	}

	return nil
}

func (s *backupService) ListBackups(ctx context.Context) ([]*BackupInfo, error) {
	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []*BackupInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only list SQL backups
		name := entry.Name()
		if !strings.HasPrefix(name, "pawnshop_backup_") {
			continue
		}
		if !strings.HasSuffix(name, ".sql") && !strings.HasSuffix(name, ".sql.gz") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backup := &BackupInfo{
			Filename:   name,
			Size:       info.Size(),
			CreatedAt:  info.ModTime(),
			Compressed: strings.HasSuffix(name, ".gz"),
		}

		// Try to read description from gzip header
		if backup.Compressed {
			if desc := s.getGzipComment(filepath.Join(s.backupDir, name)); desc != "" {
				backup.Description = desc
			}
		}

		backups = append(backups, backup)
	}

	return backups, nil
}

func (s *backupService) getGzipComment(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		return ""
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return ""
	}
	defer gzReader.Close()

	return gzReader.Comment
}

func (s *backupService) DeleteBackup(ctx context.Context, filename string) error {
	filepath := filepath.Join(s.backupDir, filename)

	// Security check: ensure filename doesn't contain path traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return fmt.Errorf("invalid filename")
	}

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", filename)
	}

	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("failed to delete backup: %w", err)
	}

	return nil
}

func (s *backupService) GetBackup(ctx context.Context, filename string) (io.ReadCloser, *BackupInfo, error) {
	// Security check
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return nil, nil, fmt.Errorf("invalid filename")
	}

	filepath := filepath.Join(s.backupDir, filename)

	// Check if file exists
	fileInfo, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("backup file not found: %s", filename)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to stat backup file: %w", err)
	}

	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open backup file: %w", err)
	}

	info := &BackupInfo{
		Filename:   filename,
		Size:       fileInfo.Size(),
		CreatedAt:  fileInfo.ModTime(),
		Compressed: strings.HasSuffix(filename, ".gz"),
	}

	return file, info, nil
}

func (s *backupService) ScheduledBackup(ctx context.Context) (*BackupInfo, error) {
	description := fmt.Sprintf("Scheduled backup - %s", time.Now().Format("2006-01-02 15:04:05"))
	return s.CreateBackup(ctx, description)
}

func (s *backupService) CleanupOldBackups(ctx context.Context, retentionDays int) (int, error) {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		return 0, fmt.Errorf("failed to read backup directory: %w", err)
	}

	deleted := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, "pawnshop_backup_") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			filepath := filepath.Join(s.backupDir, name)
			if err := os.Remove(filepath); err == nil {
				deleted++
			}
		}
	}

	return deleted, nil
}
