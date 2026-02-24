package pipeline

import (
	"testing"

	"github.com/crit/gif2vid/internal/config"
)

func TestEven(t *testing.T) {
	tests := []struct {
		in   int
		want int
	}{
		{0, 0},
		{1, 2},
		{2, 2},
		{3, 4},
		{4, 4},
		{10, 10},
		{11, 12},
	}
	for _, tt := range tests {
		got := even(tt.in)
		if got != tt.want {
			t.Errorf("even(%d) = %d; want %d", tt.in, got, tt.want)
		}
	}
}

func TestBuildFilter(t *testing.T) {
	cfg := &config.Config{
		FPS: 30,
		BG:  "black",
	}
	got := BuildFilter(cfg, 1920, 1080)
	want := "fps=30,scale=1920:1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2:color=black,format=yuv420p"
	if got != want {
		t.Errorf("BuildFilter(...) = %q; want %q", got, want)
	}
}
