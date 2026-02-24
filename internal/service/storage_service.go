package service

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"image/jpeg"
)

const (
	// MaxFileSize is the maximum file size (10MB)
	MaxFileSize = 10 * 1024 * 1024
	// ThumbnailWidth is the width of thumbnails
	ThumbnailWidth = 200
	// ThumbnailHeight is the height of thumbnails
	ThumbnailHeight = 200
)

var allowedMimeTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// ImageInfo contains information about an uploaded image
type ImageInfo struct {
	ID           string `json:"id"`
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	Size         int64  `json:"size"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	Width        int    `json:"width,omitempty"`
	Height       int    `json:"height,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// StorageService defines the interface for file storage operations
type StorageService interface {
	// UploadImage uploads an image file
	UploadImage(ctx context.Context, file *multipart.FileHeader, category string) (*ImageInfo, error)

	// UploadImageFromReader uploads an image from an io.Reader
	UploadImageFromReader(ctx context.Context, reader io.Reader, filename, mimeType, category string) (*ImageInfo, error)

	// GetImage retrieves an image by ID
	GetImage(ctx context.Context, id string) (io.ReadCloser, *ImageInfo, error)

	// GetThumbnail retrieves a thumbnail by ID
	GetThumbnail(ctx context.Context, id string) (io.ReadCloser, *ImageInfo, error)

	// DeleteImage deletes an image and its thumbnail
	DeleteImage(ctx context.Context, id string) error

	// ListImages lists images in a category
	ListImages(ctx context.Context, category string) ([]*ImageInfo, error)

	// GetImageURL returns the URL for an image
	GetImageURL(id string) string

	// GetThumbnailURL returns the URL for a thumbnail
	GetThumbnailURL(id string) string
}

type storageService struct {
	baseDir string
	baseURL string
}

// NewStorageService creates a new storage service
func NewStorageService(baseDir, baseURL string) StorageService {
	// Create base directories
	os.MkdirAll(filepath.Join(baseDir, "images"), 0755)
	os.MkdirAll(filepath.Join(baseDir, "thumbnails"), 0755)

	return &storageService{
		baseDir: baseDir,
		baseURL: baseURL,
	}
}

func (s *storageService) UploadImage(ctx context.Context, file *multipart.FileHeader, category string) (*ImageInfo, error) {
	// Validate file size
	if file.Size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum of %d bytes", MaxFileSize)
	}

	// Validate mime type
	mimeType := file.Header.Get("Content-Type")
	ext, ok := allowedMimeTypes[mimeType]
	if !ok {
		return nil, fmt.Errorf("unsupported file type: %s", mimeType)
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	return s.uploadFromReader(ctx, src, file.Filename, mimeType, ext, category, file.Size)
}

func (s *storageService) UploadImageFromReader(ctx context.Context, reader io.Reader, filename, mimeType, category string) (*ImageInfo, error) {
	ext, ok := allowedMimeTypes[mimeType]
	if !ok {
		return nil, fmt.Errorf("unsupported file type: %s", mimeType)
	}

	return s.uploadFromReader(ctx, reader, filename, mimeType, ext, category, 0)
}

func (s *storageService) uploadFromReader(_ context.Context, reader io.Reader, originalName, mimeType, ext, category string, size int64) (*ImageInfo, error) {
	// Generate unique ID
	id := uuid.New().String()
	filename := id + ext

	// Create category directory
	categoryDir := filepath.Join(s.baseDir, "images", category)
	thumbnailDir := filepath.Join(s.baseDir, "thumbnails", category)
	os.MkdirAll(categoryDir, 0755)
	os.MkdirAll(thumbnailDir, 0755)

	// Save original image
	imagePath := filepath.Join(categoryDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	// Copy content
	written, err := io.Copy(dst, reader)
	dst.Close()
	if err != nil {
		os.Remove(imagePath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	if size == 0 {
		size = written
	}

	// Get image dimensions and create thumbnail
	width, height := 0, 0
	thumbnailURL := ""

	file, err := os.Open(imagePath)
	if err == nil {
		defer file.Close()
		img, _, err := image.Decode(file)
		if err == nil {
			bounds := img.Bounds()
			width = bounds.Dx()
			height = bounds.Dy()

			// Create thumbnail
			thumb := resize.Thumbnail(ThumbnailWidth, ThumbnailHeight, img, resize.Lanczos3)
			thumbPath := filepath.Join(thumbnailDir, filename)
			thumbFile, err := os.Create(thumbPath)
			if err == nil {
				jpeg.Encode(thumbFile, thumb, &jpeg.Options{Quality: 85})
				thumbFile.Close()
				thumbnailURL = s.GetThumbnailURL(filepath.Join(category, filename))
			}
		}
	}

	return &ImageInfo{
		ID:           filepath.Join(category, filename),
		Filename:     filename,
		OriginalName: originalName,
		MimeType:     mimeType,
		Size:         size,
		URL:          s.GetImageURL(filepath.Join(category, filename)),
		ThumbnailURL: thumbnailURL,
		Width:        width,
		Height:       height,
		CreatedAt:    time.Now(),
	}, nil
}

func (s *storageService) GetImage(ctx context.Context, id string) (io.ReadCloser, *ImageInfo, error) {
	return s.getFile(ctx, "images", id)
}

func (s *storageService) GetThumbnail(ctx context.Context, id string) (io.ReadCloser, *ImageInfo, error) {
	return s.getFile(ctx, "thumbnails", id)
}

func (s *storageService) getFile(_ context.Context, subdir, id string) (io.ReadCloser, *ImageInfo, error) {
	// Security check
	if strings.Contains(id, "..") {
		return nil, nil, fmt.Errorf("invalid file ID")
	}

	dir := filepath.Join(s.baseDir, subdir, filepath.Dir(id))
	baseID := filepath.Base(id)

	var filePath string
	var filename string

	// First, try exact match (if ID includes extension)
	exactPath := filepath.Join(s.baseDir, subdir, id)
	if _, err := os.Stat(exactPath); err == nil {
		filePath = exactPath
		filename = baseID
	} else {
		// Fall back to prefix search (for IDs without extension)
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, nil, fmt.Errorf("file not found")
		}

		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, baseID+".") || name == baseID {
				filePath = filepath.Join(dir, name)
				filename = name
				break
			}
		}
	}

	if filePath == "" {
		return nil, nil, fmt.Errorf("file not found")
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("file not found")
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}

	// Determine mime type from extension
	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := "application/octet-stream"
	for mime, e := range allowedMimeTypes {
		if e == ext {
			mimeType = mime
			break
		}
	}

	info := &ImageInfo{
		ID:        id,
		Filename:  filename,
		MimeType:  mimeType,
		Size:      fileInfo.Size(),
		CreatedAt: fileInfo.ModTime(),
	}

	return file, info, nil
}

func (s *storageService) DeleteImage(ctx context.Context, id string) error {
	// Security check
	if strings.Contains(id, "..") {
		return fmt.Errorf("invalid file ID")
	}

	// Delete original
	imageDir := filepath.Join(s.baseDir, "images", filepath.Dir(id))
	thumbDir := filepath.Join(s.baseDir, "thumbnails", filepath.Dir(id))
	baseID := filepath.Base(id)

	// Find and delete files
	for _, dir := range []string{imageDir, thumbDir} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, baseID+".") {
				os.Remove(filepath.Join(dir, name))
			}
		}
	}

	return nil
}

func (s *storageService) ListImages(ctx context.Context, category string) ([]*ImageInfo, error) {
	dir := filepath.Join(s.baseDir, "images", category)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*ImageInfo{}, nil
		}
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var images []*ImageInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := filepath.Ext(name)
		id := strings.TrimSuffix(name, ext)

		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Determine mime type
		mimeType := "application/octet-stream"
		for mime, e := range allowedMimeTypes {
			if e == strings.ToLower(ext) {
				mimeType = mime
				break
			}
		}

		fullID := filepath.Join(category, id)
		images = append(images, &ImageInfo{
			ID:           fullID,
			Filename:     name,
			MimeType:     mimeType,
			Size:         info.Size(),
			URL:          s.GetImageURL(fullID),
			ThumbnailURL: s.GetThumbnailURL(fullID),
			CreatedAt:    info.ModTime(),
		})
	}

	return images, nil
}

func (s *storageService) GetImageURL(id string) string {
	return fmt.Sprintf("%s/images/%s", s.baseURL, id)
}

func (s *storageService) GetThumbnailURL(id string) string {
	return fmt.Sprintf("%s/thumbnails/%s", s.baseURL, id)
}
