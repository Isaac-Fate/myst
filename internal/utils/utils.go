package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func ResolvePath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		// Get the home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		// Replace the tilde with the home directory
		path = filepath.Join(homeDir, path[2:])
	}

	// Convert to absolute path
	return filepath.Abs(path)
}
