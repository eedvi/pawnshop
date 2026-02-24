package service

import (
	"compress/gzip"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pawnshop/internal/config"
)

func setupBackupTestDir(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "backup_test_*")
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func createTestBackupFile(t *testing.T, dir, filename, content, comment string) {
	filepath := filepath.Join(dir, filename)

	file, err := os.Create(filepath)
	require.NoError(t, err)
	defer file.Close()

	if comment != "" || len(content) > 0 {
		gzWriter := gzip.NewWriter(file)
		gzWriter.Comment = comment
		gzWriter.ModTime = time.Now()
		gzWriter.Write([]byte(content))
		gzWriter.Close()
	}
}

func createTestSQLFile(t *testing.T, dir, filename, content string) {
	filepath := filepath.Join(dir, filename)
	err := os.WriteFile(filepath, []byte(content), 0644)
	require.NoError(t, err)
}

func TestNewBackupService(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "testdb",
	}

	// Create service with a new subdirectory
	backupDir := filepath.Join(tempDir, "backups")
	svc := NewBackupService(dbConfig, backupDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))

	assert.NotNil(t, svc)

	// Verify directory was created
	_, err := os.Stat(backupDir)
	assert.NoError(t, err)
}

func TestBackupService_ListBackups(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	// Create some test backup files
	createTestBackupFile(t, tempDir, "pawnshop_backup_20240101_120000.sql.gz", "SQL content", "Test backup 1")
	createTestBackupFile(t, tempDir, "pawnshop_backup_20240102_120000.sql.gz", "SQL content 2", "Test backup 2")
	createTestSQLFile(t, tempDir, "pawnshop_backup_20240103_120000.sql", "Plain SQL")
	createTestSQLFile(t, tempDir, "other_file.txt", "Not a backup") // Should be ignored
	os.Mkdir(filepath.Join(tempDir, "subdir"), 0755)                 // Should be ignored

	backups, err := svc.ListBackups(ctx)
	require.NoError(t, err)
	assert.Len(t, backups, 3)

	// Verify backup info
	backupNames := make(map[string]*BackupInfo)
	for _, b := range backups {
		backupNames[b.Filename] = b
	}

	assert.Contains(t, backupNames, "pawnshop_backup_20240101_120000.sql.gz")
	assert.Contains(t, backupNames, "pawnshop_backup_20240102_120000.sql.gz")
	assert.Contains(t, backupNames, "pawnshop_backup_20240103_120000.sql")

	// Verify compressed flag
	assert.True(t, backupNames["pawnshop_backup_20240101_120000.sql.gz"].Compressed)
	assert.False(t, backupNames["pawnshop_backup_20240103_120000.sql"].Compressed)

	// Verify description was read from gzip header
	assert.Equal(t, "Test backup 1", backupNames["pawnshop_backup_20240101_120000.sql.gz"].Description)
}

func TestBackupService_ListBackups_EmptyDir(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	backups, err := svc.ListBackups(ctx)
	require.NoError(t, err)
	assert.Empty(t, backups)
}

func TestBackupService_DeleteBackup(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	// Create a test backup file
	filename := "pawnshop_backup_20240101_120000.sql.gz"
	createTestBackupFile(t, tempDir, filename, "SQL content", "Test")

	// Verify file exists
	_, err := os.Stat(filepath.Join(tempDir, filename))
	require.NoError(t, err)

	// Delete the backup
	err = svc.DeleteBackup(ctx, filename)
	require.NoError(t, err)

	// Verify file is deleted
	_, err = os.Stat(filepath.Join(tempDir, filename))
	assert.True(t, os.IsNotExist(err))
}

func TestBackupService_DeleteBackup_NotFound(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	err := svc.DeleteBackup(ctx, "nonexistent.sql.gz")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestBackupService_DeleteBackup_PathTraversal(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	tests := []struct {
		name     string
		filename string
	}{
		{"double dots", "../etc/passwd"},
		{"forward slash", "subdir/backup.sql"},
		{"backslash", "subdir\\backup.sql"},
		{"hidden double dots", "..\\..\\important.sql"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.DeleteBackup(ctx, tt.filename)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid filename")
		})
	}
}

func TestBackupService_GetBackup(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	// Create a test backup file
	filename := "pawnshop_backup_20240101_120000.sql.gz"
	content := "SQL backup content here"
	createTestBackupFile(t, tempDir, filename, content, "Test description")

	// Get the backup
	reader, info, err := svc.GetBackup(ctx, filename)
	require.NoError(t, err)
	defer reader.Close()

	assert.Equal(t, filename, info.Filename)
	assert.True(t, info.Compressed)
	assert.True(t, info.Size > 0)
}

func TestBackupService_GetBackup_NotFound(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	reader, info, err := svc.GetBackup(ctx, "nonexistent.sql.gz")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Nil(t, reader)
	assert.Nil(t, info)
}

func TestBackupService_GetBackup_PathTraversal(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	tests := []struct {
		name     string
		filename string
	}{
		{"double dots", "../etc/passwd"},
		{"forward slash", "dir/backup.sql"},
		{"backslash", "dir\\backup.sql"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, info, err := svc.GetBackup(ctx, tt.filename)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid filename")
			assert.Nil(t, reader)
			assert.Nil(t, info)
		})
	}
}

func TestBackupService_CleanupOldBackups(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	// Create backup files with different ages
	// Note: We can't easily modify file times, so we create files and test the logic
	// Create files that should NOT be deleted (recent)
	createTestBackupFile(t, tempDir, "pawnshop_backup_recent.sql.gz", "content", "")

	// To properly test cleanup, we need to modify the file's mod time
	// Create an "old" backup and modify its timestamp
	oldFilename := "pawnshop_backup_old.sql.gz"
	createTestBackupFile(t, tempDir, oldFilename, "old content", "")

	oldFilePath := filepath.Join(tempDir, oldFilename)
	oldTime := time.Now().AddDate(0, 0, -31) // 31 days ago
	os.Chtimes(oldFilePath, oldTime, oldTime)

	// Also create a very old file
	veryOldFilename := "pawnshop_backup_very_old.sql.gz"
	createTestBackupFile(t, tempDir, veryOldFilename, "very old content", "")
	veryOldFilePath := filepath.Join(tempDir, veryOldFilename)
	veryOldTime := time.Now().AddDate(0, 0, -60) // 60 days ago
	os.Chtimes(veryOldFilePath, veryOldTime, veryOldTime)

	// Cleanup backups older than 30 days
	deleted, err := svc.CleanupOldBackups(ctx, 30)
	require.NoError(t, err)
	assert.Equal(t, 2, deleted)

	// Verify recent backup still exists
	_, err = os.Stat(filepath.Join(tempDir, "pawnshop_backup_recent.sql.gz"))
	assert.NoError(t, err)

	// Verify old backups are deleted
	_, err = os.Stat(oldFilePath)
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(veryOldFilePath)
	assert.True(t, os.IsNotExist(err))
}

func TestBackupService_CleanupOldBackups_NoOldBackups(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	// Create only recent backups
	createTestBackupFile(t, tempDir, "pawnshop_backup_20240101.sql.gz", "content", "")

	deleted, err := svc.CleanupOldBackups(ctx, 30)
	require.NoError(t, err)
	assert.Equal(t, 0, deleted)
}

func TestBackupService_CleanupOldBackups_SkipsNonBackups(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	// Create non-backup files with old timestamps
	err := os.WriteFile(filepath.Join(tempDir, "other_file.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	oldTime := time.Now().AddDate(0, 0, -60)
	os.Chtimes(filepath.Join(tempDir, "other_file.txt"), oldTime, oldTime)

	// Should not delete non-backup files
	deleted, err := svc.CleanupOldBackups(ctx, 30)
	require.NoError(t, err)
	assert.Equal(t, 0, deleted)

	// Verify file still exists
	_, err = os.Stat(filepath.Join(tempDir, "other_file.txt"))
	assert.NoError(t, err)
}

func TestBackupInfo_Structure(t *testing.T) {
	now := time.Now()
	info := BackupInfo{
		Filename:    "test_backup.sql.gz",
		Size:        12345,
		CreatedAt:   now,
		Compressed:  true,
		Description: "Test backup description",
	}

	assert.Equal(t, "test_backup.sql.gz", info.Filename)
	assert.Equal(t, int64(12345), info.Size)
	assert.Equal(t, now, info.CreatedAt)
	assert.True(t, info.Compressed)
	assert.Equal(t, "Test backup description", info.Description)
}

func TestBackupService_RestoreBackup_FileNotFound(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "testdb",
	}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	err := svc.RestoreBackup(ctx, "nonexistent_backup.sql.gz")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestBackupService_ListBackups_WithSubdirectories(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	// Create a backup and a subdirectory
	createTestBackupFile(t, tempDir, "pawnshop_backup_20240101.sql.gz", "content", "")

	// Create subdirectory that shouldn't be listed
	subdir := filepath.Join(tempDir, "pawnshop_backup_subdir")
	os.Mkdir(subdir, 0755)

	backups, err := svc.ListBackups(ctx)
	require.NoError(t, err)
	assert.Len(t, backups, 1)
	assert.Equal(t, "pawnshop_backup_20240101.sql.gz", backups[0].Filename)
}

func TestBackupService_GetBackup_NonCompressed(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled))
	ctx := context.Background()

	// Create a non-compressed backup
	filename := "pawnshop_backup_20240101_120000.sql"
	content := "SQL backup content"
	createTestSQLFile(t, tempDir, filename, content)

	reader, info, err := svc.GetBackup(ctx, filename)
	require.NoError(t, err)
	defer reader.Close()

	assert.Equal(t, filename, info.Filename)
	assert.False(t, info.Compressed)
}

func TestBackupService_getGzipComment(t *testing.T) {
	tempDir, cleanup := setupBackupTestDir(t)
	defer cleanup()

	dbConfig := &config.DatabaseConfig{}
	svc := NewBackupService(dbConfig, tempDir, zerolog.New(os.Stdout).Level(zerolog.Disabled)).(*backupService)

	// Create gzipped file with comment
	filename := "test_with_comment.sql.gz"
	comment := "This is a test comment"
	createTestBackupFile(t, tempDir, filename, "content", comment)

	result := svc.getGzipComment(filepath.Join(tempDir, filename))
	assert.Equal(t, comment, result)

	// Test non-existent file
	result = svc.getGzipComment(filepath.Join(tempDir, "nonexistent.gz"))
	assert.Empty(t, result)

	// Test non-gzip file
	createTestSQLFile(t, tempDir, "plain.sql", "content")
	result = svc.getGzipComment(filepath.Join(tempDir, "plain.sql"))
	assert.Empty(t, result)
}
