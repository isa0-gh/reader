package api

import (
	"fmt"
	"strings"
)

const MaxImageSize = 10 * 1024 * 1024 // 10MB

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func ValidateImageUpload(contentType string, size int64) error {
	if !allowedImageTypes[strings.ToLower(contentType)] {
		return fmt.Errorf("invalid image format: must be JPEG, PNG, GIF, or WebP")
	}
	if size > MaxImageSize {
		return fmt.Errorf("image too large: maximum size is 10MB")
	}
	return nil
}
