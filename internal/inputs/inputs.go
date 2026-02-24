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

// ValidateAndAbs ensures inputs exist as files and returns absolute cleaned paths.
func ValidateAndAbs(paths []string) ([]string, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no inputs provided")
	}
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		st, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("input not found: %s", p)
		}
		if !st.Mode().IsRegular() {
			return nil, fmt.Errorf("input is not a regular file: %s", p)
		}
		ext := strings.ToLower(filepath.Ext(p))
		if !allowedExt[ext] {
			// future: allow-any optional; for now, enforce
			return nil, fmt.Errorf("unsupported extension %s for %s (allowed: .gif,.webp)", ext, p)
		}
		ap, err := filepath.Abs(p)
		if err != nil {
			return nil, err
		}
		out = append(out, filepath.Clean(ap))
	}
	return out, nil
}
