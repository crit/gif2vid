package inputs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var allowedExt = map[string]bool{
	".gif":  true,
	".webp": true,
}

// GetFilesFromDir scans the directory for GIF and WebP files and returns absolute cleaned paths.
func GetFilesFromDir(dirPath string) ([]string, error) {
	st, err := os.Stat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("input directory not found: %s", dirPath)
	}
	if !st.IsDir() {
		return nil, fmt.Errorf("input is not a directory: %s", dirPath)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var out []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if !allowedExt[ext] {
			continue
		}

		p := filepath.Join(dirPath, entry.Name())
		ap, err := filepath.Abs(p)
		if err != nil {
			return nil, err
		}
		out = append(out, filepath.Clean(ap))
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no supported files (.gif, .webp) found in: %s", dirPath)
	}

	return out, nil
}
