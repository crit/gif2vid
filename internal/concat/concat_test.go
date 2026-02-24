package concat

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestEscapePathForConcat(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"simple.mp4", "simple.mp4"},
		{"it's.mp4", "it'\\''s.mp4"},
		{"multi'quote'test.mp4", "multi'\\''quote'\\''test.mp4"},
	}
	for _, tt := range tests {
		got := EscapePathForConcat(tt.in)
		// Since EscapePathForConcat may call filepath.Abs, we check if it ends with our expected escaped string.
		if !strings.HasSuffix(got, tt.want) {
			t.Errorf("EscapePathForConcat(%q) = %q; want suffix %q", tt.in, got, tt.want)
		}
	}
}

func TestEscapePathForConcat_Abs(t *testing.T) {
	absIn, _ := filepath.Abs("foo'bar.mp4")
	want := strings.ReplaceAll(absIn, "'", "'\\''")
	got := EscapePathForConcat(absIn)
	if got != want {
		t.Errorf("EscapePathForConcat(abs %q) = %q; want %q", absIn, got, want)
	}
}
