package media

import (
	"context"
	"encoding/json"
	"fmt"

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
	stdout, stderr, err := r.Run(ctx, "ffprobe", args)
	if err != nil {
		return 0, 0, fmt.Errorf("ffprobe failed for %s: %v\n%s", input, err, string(stderr))
	}
	var pr ProbeResult
	if err := json.Unmarshal(stdout, &pr); err != nil {
		return 0, 0, fmt.Errorf("ffprobe parse failed for %s: %v\nstdout=%s", input, err, string(stdout))
	}
	for _, s := range pr.Streams {
		if s.Width > 0 && s.Height > 0 {
			return s.Width, s.Height, nil
		}
	}
	return 0, 0, fmt.Errorf("no valid visual stream found in %s", input)
}
