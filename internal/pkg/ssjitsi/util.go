package ssjitsi

import (
	"regexp"
	"strings"
)

// SafeFilename creates a filesystem-safe filename
func SafeFilename(filename string) string {
	// Define invalid characters for most filesystems
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)

	// Replace invalid characters with underscore
	safe := invalidChars.ReplaceAllString(filename, "_")

	// Remove leading/trailing spaces and dots (Windows restriction)
	safe = strings.Trim(safe, " .")

	// Ensure the filename is not empty
	if safe == "" {
		return "unknown"
	}

	// Limit length (optional, 255 is common filesystem limit)
	if len(safe) > 255 {
		safe = safe[:255]
	}

	return safe
}
