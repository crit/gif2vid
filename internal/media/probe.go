package media

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/crit/gif2vid/internal/config"
	"github.com/crit/gif2vid/internal/ffmpeg"
)

// ProbeResult captures parts of ffprobe JSON we care about.
type ProbeResult struct {
	Streams []struct {
		CodecType string `json:"codec_type"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
	} `json:"streams"`
}

// Probe returns the width and height of the first video stream in the file.
func Probe(ctx context.Context, r ffmpeg.Runner, cfg *config.Config, input string) (int, int, error) {
	args := []string{
		"-v", "error",
		"-show_entries", "stream=width,height,codec_type",
		"-of", "json",
		input,
	}
	stdout, _, err := r.Run(ctx, "ffprobe", args)
	if err != nil {
		return probeFallback(ctx, r, cfg, input)
	}
	var pr ProbeResult
	if err := json.Unmarshal(stdout, &pr); err != nil {
		return probeFallback(ctx, r, cfg, input)
	}
	for _, s := range pr.Streams {
		if s.Width > 0 && s.Height > 0 {
			return s.Width, s.Height, nil
		}
	}
	return probeFallback(ctx, r, cfg, input)
}

func probeFallback(ctx context.Context, r ffmpeg.Runner, cfg *config.Config, input string) (int, int, error) {
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "gif2vid_probe_fallback.png")
	defer os.Remove(tmpFile)

	args := []string{
		"-v", "error",
		"-y",
		"-i", input,
		"-vframes", "1",
		tmpFile,
	}
	_, _, err := r.Run(ctx, "ffmpeg", args)
	if err != nil {
		return probeMagickFallback(ctx, r, cfg, input)
	}

	// Now probe the generated frame
	probeArgs := []string{
		"-v", "error",
		"-show_entries", "stream=width,height",
		"-of", "json",
		tmpFile,
	}
	stdout, _, probeErr := r.Run(ctx, "ffprobe", probeArgs)
	if probeErr != nil {
		return probeMagickFallback(ctx, r, cfg, input)
	}
	var pr ProbeResult
	if err := json.Unmarshal(stdout, &pr); err != nil {
		return probeMagickFallback(ctx, r, cfg, input)
	}
	for _, s := range pr.Streams {
		if s.Width > 0 && s.Height > 0 {
			return s.Width, s.Height, nil
		}
	}
	return probeMagickFallback(ctx, r, cfg, input)
}

func probeMagickFallback(ctx context.Context, r ffmpeg.Runner, cfg *config.Config, input string) (int, int, error) {
	if cfg.MagickBin == "" {
		return 0, 0, fmt.Errorf("no valid visual stream found in %s and ImageMagick not available", input)
	}

	// Use identify to get dimensions: identify -format "%w %h" input.webp[0]
	// [0] ensures we only look at the first frame
	args := []string{"-format", "%w %h", input + "[0]"}
	bin := cfg.MagickBin
	if bin == "magick" {
		args = append([]string{"identify"}, args...)
	} else {
		bin = "identify"
	}

	stdout, stderr, err := r.Run(ctx, bin, args)
	if err != nil {
		return 0, 0, fmt.Errorf("magick identify failed for %s: %v\n%s", input, err, string(stderr))
	}

	fields := strings.Fields(string(stdout))
	if len(fields) < 2 {
		return 0, 0, fmt.Errorf("magick identify returned unexpected output for %s: %s", input, string(stdout))
	}

	w, errW := strconv.Atoi(fields[0])
	h, errH := strconv.Atoi(fields[1])
	if errW != nil || errH != nil {
		return 0, 0, fmt.Errorf("magick identify returned invalid dimensions for %s: %s", input, string(stdout))
	}

	return w, h, nil
}
