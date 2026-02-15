package service

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupStorageTestDir(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "storage_test_*")
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

// createTestJPEGImage creates a simple JPEG image for testing
func createTestJPEGImage(t *testing.T, width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with a solid color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}

	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	require.NoError(t, err)
	return buf.Bytes()
}

// createTestPNGImage creates a simple PNG image for testing
func createTestPNGImage(t *testing.T, width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 50, G: 100, B: 150, A: 255})
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	require.NoError(t, err)
	return buf.Bytes()
}

// createMockMultipartHeader creates a mock multipart.FileHeader for testing
func createMockMultipartHeader(filename, contentType string, content []byte) *multipart.FileHeader {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", contentType)

	return &multipart.FileHeader{
		Filename: filename,
		Header:   header,
		Size:     int64(len(content)),
	}
}

func TestNewStorageService(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	assert.NotNil(t, svc)

	// Verify directories were created
	_, err := os.Stat(filepath.Join(tempDir, "images"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(tempDir, "thumbnails"))
	assert.NoError(t, err)
}

func TestStorageService_UploadImageFromReader(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	imageData := createTestJPEGImage(t, 400, 300)
	reader := bytes.NewReader(imageData)

	info, err := svc.UploadImageFromReader(ctx, reader, "test_image.jpg", "image/jpeg", "items")
	require.NoError(t, err)

	assert.NotEmpty(t, info.ID)
	assert.Equal(t, "test_image.jpg", info.OriginalName)
	assert.Equal(t, "image/jpeg", info.MimeType)
	assert.True(t, info.Size > 0)
	assert.Contains(t, info.URL, "http://localhost:8080/storage/images/")
	assert.Equal(t, 400, info.Width)
	assert.Equal(t, 300, info.Height)
}

func TestStorageService_UploadImageFromReader_PNG(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	imageData := createTestPNGImage(t, 200, 150)
	reader := bytes.NewReader(imageData)

	info, err := svc.UploadImageFromReader(ctx, reader, "test.png", "image/png", "customers")
	require.NoError(t, err)

	assert.Equal(t, "image/png", info.MimeType)
	assert.True(t, strings.HasSuffix(info.Filename, ".png"))
	assert.Equal(t, 200, info.Width)
	assert.Equal(t, 150, info.Height)
}

func TestStorageService_UploadImageFromReader_UnsupportedType(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	reader := strings.NewReader("not an image")

	_, err := svc.UploadImageFromReader(ctx, reader, "test.txt", "text/plain", "items")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported file type")
}

func TestStorageService_GetImage(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	// Upload an image first
	imageData := createTestJPEGImage(t, 100, 100)
	reader := bytes.NewReader(imageData)

	uploadedInfo, err := svc.UploadImageFromReader(ctx, reader, "test.jpg", "image/jpeg", "items")
	require.NoError(t, err)

	// Get the image
	file, info, err := svc.GetImage(ctx, uploadedInfo.ID)
	require.NoError(t, err)
	defer file.Close()

	assert.Equal(t, uploadedInfo.ID, info.ID)
	assert.Equal(t, "image/jpeg", info.MimeType)
	assert.True(t, info.Size > 0)

	// Verify content
	content, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestStorageService_GetImage_NotFound(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	_, _, err := svc.GetImage(ctx, "items/nonexistent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestStorageService_GetImage_PathTraversal(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	_, _, err := svc.GetImage(ctx, "../../../etc/passwd")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

func TestStorageService_GetThumbnail(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	// Upload an image (thumbnail is created automatically)
	imageData := createTestJPEGImage(t, 400, 400)
	reader := bytes.NewReader(imageData)

	uploadedInfo, err := svc.UploadImageFromReader(ctx, reader, "test.jpg", "image/jpeg", "items")
	require.NoError(t, err)

	// Get the thumbnail
	file, info, err := svc.GetThumbnail(ctx, uploadedInfo.ID)
	require.NoError(t, err)
	defer file.Close()

	assert.Equal(t, uploadedInfo.ID, info.ID)
	assert.True(t, info.Size > 0)

	// Verify content
	content, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestStorageService_DeleteImage(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	// Upload an image
	imageData := createTestJPEGImage(t, 100, 100)
	reader := bytes.NewReader(imageData)

	uploadedInfo, err := svc.UploadImageFromReader(ctx, reader, "test.jpg", "image/jpeg", "items")
	require.NoError(t, err)

	// Verify image exists
	file, _, err := svc.GetImage(ctx, uploadedInfo.ID)
	require.NoError(t, err)
	file.Close() // Important: close file before attempting delete (especially on Windows)

	// Delete the image
	err = svc.DeleteImage(ctx, uploadedInfo.ID)
	require.NoError(t, err)

	// Verify image is deleted
	_, _, err = svc.GetImage(ctx, uploadedInfo.ID)
	assert.Error(t, err)
}

func TestStorageService_DeleteImage_PathTraversal(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	err := svc.DeleteImage(ctx, "../../../important-file")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

func TestStorageService_ListImages(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	// Upload multiple images
	for i := 0; i < 3; i++ {
		imageData := createTestJPEGImage(t, 100, 100)
		reader := bytes.NewReader(imageData)
		_, err := svc.UploadImageFromReader(ctx, reader, "test.jpg", "image/jpeg", "items")
		require.NoError(t, err)
	}

	// List images
	images, err := svc.ListImages(ctx, "items")
	require.NoError(t, err)
	assert.Len(t, images, 3)

	for _, img := range images {
		assert.NotEmpty(t, img.ID)
		assert.Equal(t, "image/jpeg", img.MimeType)
		assert.NotEmpty(t, img.URL)
		assert.NotEmpty(t, img.ThumbnailURL)
	}
}

func TestStorageService_ListImages_EmptyCategory(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	images, err := svc.ListImages(ctx, "empty_category")
	require.NoError(t, err)
	assert.Empty(t, images)
}

func TestStorageService_ListImages_DifferentCategories(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	// Upload to different categories
	imageData := createTestJPEGImage(t, 100, 100)

	_, err := svc.UploadImageFromReader(ctx, bytes.NewReader(imageData), "item.jpg", "image/jpeg", "items")
	require.NoError(t, err)

	_, err = svc.UploadImageFromReader(ctx, bytes.NewReader(imageData), "customer.jpg", "image/jpeg", "customers")
	require.NoError(t, err)

	// List only items
	itemImages, err := svc.ListImages(ctx, "items")
	require.NoError(t, err)
	assert.Len(t, itemImages, 1)

	// List only customers
	customerImages, err := svc.ListImages(ctx, "customers")
	require.NoError(t, err)
	assert.Len(t, customerImages, 1)
}

func TestStorageService_GetImageURL(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")

	url := svc.GetImageURL("items/abc123")
	assert.Equal(t, "http://localhost:8080/storage/images/items/abc123", url)
}

func TestStorageService_GetThumbnailURL(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")

	url := svc.GetThumbnailURL("items/abc123")
	assert.Equal(t, "http://localhost:8080/storage/thumbnails/items/abc123", url)
}

func TestImageInfo_Structure(t *testing.T) {
	info := ImageInfo{
		ID:           "items/abc123",
		Filename:     "abc123.jpg",
		OriginalName: "my_photo.jpg",
		MimeType:     "image/jpeg",
		Size:         12345,
		URL:          "http://example.com/images/items/abc123",
		ThumbnailURL: "http://example.com/thumbnails/items/abc123",
		Width:        800,
		Height:       600,
	}

	assert.Equal(t, "items/abc123", info.ID)
	assert.Equal(t, "abc123.jpg", info.Filename)
	assert.Equal(t, "my_photo.jpg", info.OriginalName)
	assert.Equal(t, "image/jpeg", info.MimeType)
	assert.Equal(t, int64(12345), info.Size)
	assert.Equal(t, 800, info.Width)
	assert.Equal(t, 600, info.Height)
}

func TestStorageService_UploadImage_LargeImage(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	// Create a larger image to test thumbnail generation
	imageData := createTestJPEGImage(t, 1000, 800)
	reader := bytes.NewReader(imageData)

	info, err := svc.UploadImageFromReader(ctx, reader, "large.jpg", "image/jpeg", "items")
	require.NoError(t, err)

	assert.Equal(t, 1000, info.Width)
	assert.Equal(t, 800, info.Height)
	assert.NotEmpty(t, info.ThumbnailURL)

	// Verify thumbnail was created
	_, thumbInfo, err := svc.GetThumbnail(ctx, info.ID)
	require.NoError(t, err)
	assert.True(t, thumbInfo.Size > 0)
}

func TestStorageService_AllowedMimeTypes(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	tests := []struct {
		mimeType    string
		shouldAllow bool
	}{
		{"image/jpeg", true},
		{"image/png", true},
		{"image/gif", true},
		{"image/webp", true},
		{"text/plain", false},
		{"application/pdf", false},
		{"image/bmp", false},
	}

	for _, tt := range tests {
		t.Run(tt.mimeType, func(t *testing.T) {
			var data []byte
			if tt.mimeType == "image/jpeg" {
				data = createTestJPEGImage(t, 50, 50)
			} else if tt.mimeType == "image/png" {
				data = createTestPNGImage(t, 50, 50)
			} else {
				data = []byte("test content")
			}

			_, err := svc.UploadImageFromReader(ctx, bytes.NewReader(data), "test", tt.mimeType, "test")

			if tt.shouldAllow {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "unsupported")
			}
		})
	}
}

func TestStorageService_UploadCreatesCategories(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	imageData := createTestJPEGImage(t, 50, 50)
	reader := bytes.NewReader(imageData)

	// Upload to a new category
	_, err := svc.UploadImageFromReader(ctx, reader, "test.jpg", "image/jpeg", "new_category")
	require.NoError(t, err)

	// Verify directories were created
	_, err = os.Stat(filepath.Join(tempDir, "images", "new_category"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(tempDir, "thumbnails", "new_category"))
	assert.NoError(t, err)
}

func TestStorageService_DeleteNonExistent(t *testing.T) {
	tempDir, cleanup := setupStorageTestDir(t)
	defer cleanup()

	svc := NewStorageService(tempDir, "http://localhost:8080/storage")
	ctx := context.Background()

	// Deleting non-existent image should not error
	err := svc.DeleteImage(ctx, "items/nonexistent")
	assert.NoError(t, err)
}

func TestConstants(t *testing.T) {
	assert.Equal(t, int64(10*1024*1024), int64(MaxFileSize))
	assert.Equal(t, uint(200), uint(ThumbnailWidth))
	assert.Equal(t, uint(200), uint(ThumbnailHeight))
}
