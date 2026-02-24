package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// AbsClean returns an absolute, cleaned path.
func AbsClean(p string) (string, error) {
	if p == "" {
		return "", fmt.Errorf("empty path")
	}
	ap, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return filepath.Clean(ap), nil
}

// MkTempWorkspace creates a temp directory for work. If base is empty, use default.
func MkTempWorkspace(base string) (string, error) {
	pattern := "gif2vid-*"
	if base == "" {
		return os.MkdirTemp("", pattern)
	}
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", err
	}
	return os.MkdirTemp(base, pattern)
}

// AtomicRename moves temp file to destination, optionally refusing overwrite.
func AtomicRename(tmpPath, dstPath string, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(dstPath); err == nil {
			return fmt.Errorf("output exists: %s (use --overwrite)", dstPath)
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	// Ensure target dir exists
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}
	// On POSIX, os.Rename is atomic across same filesystem.
	return os.Rename(tmpPath, dstPath)
}

// WriteFile writes data to path with mode, creating parents.
func WriteFile(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, mode)
}

// CopyFile copies src to dst creating parents.
func CopyFile(src, dst string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	return err
}
