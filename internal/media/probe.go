package media

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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
func Probe(ctx context.Context, r ffmpeg.Runner, input string) (int, int, error) {
	args := []string{
		"-v", "error",
		"-show_entries", "stream=width,height,codec_type",
		"-of", "json",
		input,
	}
	stdout, _, err := r.Run(ctx, "ffprobe", args)
	if err != nil {
		return probeFallback(ctx, r, input)
	}
	var pr ProbeResult
	if err := json.Unmarshal(stdout, &pr); err != nil {
		return probeFallback(ctx, r, input)
	}
	for _, s := range pr.Streams {
		if s.Width > 0 && s.Height > 0 {
			return s.Width, s.Height, nil
		}
	}
	return probeFallback(ctx, r, input)
}

func probeFallback(ctx context.Context, r ffmpeg.Runner, input string) (int, int, error) {
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
	_, stderr, err := r.Run(ctx, "ffmpeg", args)
	if err != nil {
		return 0, 0, fmt.Errorf("probe fallback failed for %s: %v\n%s", input, err, string(stderr))
	}

	// Now probe the generated frame
	probeArgs := []string{
		"-v", "error",
		"-show_entries", "stream=width,height",
		"-of", "json",
		tmpFile,
	}
	stdout, stderr, probeErr := r.Run(ctx, "ffprobe", probeArgs)
	if probeErr != nil {
		return 0, 0, fmt.Errorf("probe fallback ffprobe failed for %s: %v\n%s", input, probeErr, string(stderr))
	}
	var pr ProbeResult
	if err := json.Unmarshal(stdout, &pr); err != nil {
		return 0, 0, fmt.Errorf("probe fallback parse failed for %s: %v", input, err)
	}
	for _, s := range pr.Streams {
		if s.Width > 0 && s.Height > 0 {
			return s.Width, s.Height, nil
		}
	}
	return 0, 0, fmt.Errorf("no valid visual stream found in %s after fallback", input)
}
