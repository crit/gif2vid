package concat

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EscapePathForConcat escapes a file path for ffmpeg concat demuxer line: file '...'
// Only need to escape single quotes by replacing ' with '\” inside single-quoted string.
func EscapePathForConcat(p string) string {
	// ensure absolute for safety with -safe 0, but caller may handle it too
	ap := p
	if !filepath.IsAbs(ap) {
		if abs, err := filepath.Abs(ap); err == nil {
			ap = abs
		}
	}
	return strings.ReplaceAll(ap, "'", "'\\''")
}

// WriteConcatFile writes the concat list file with one absolute quoted path per line.
func WriteConcatFile(path string, segmentPaths []string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, p := range segmentPaths {
		line := fmt.Sprintf("file '%s'\n", EscapePathForConcat(p))
		if _, err := w.WriteString(line); err != nil {
			return err
		}
	}
	return w.Flush()
}
