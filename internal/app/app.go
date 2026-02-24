package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/crit/gif2vid/internal/config"
	"github.com/crit/gif2vid/internal/ffmpeg"
	"github.com/crit/gif2vid/internal/inputs"
	"github.com/crit/gif2vid/internal/pipeline"
)

// Run is the main orchestration entry point.
func Run(ctx context.Context, cfg *config.Config) error {
	// Check environment binaries early
	if _, err := ffmpeg.LookPath("ffmpeg"); err != nil {
		return err
	}
	if _, err := ffmpeg.LookPath("ffprobe"); err != nil {
		return err
	}
	if cfg.Verbose {
		fmt.Println("[gif2vid] ffmpeg/ffprobe found in PATH")
	}

	// Validate inputs
	absInputs, err := inputs.ValidateAndAbs(cfg.Inputs)
	if err != nil {
		return err
	}
	cfg.Inputs = absInputs

	// Ensure output parent exists (later we also check overwrite)
	if err := os.MkdirAll(filepath.Dir(cfg.Output), 0o755); err != nil {
		return err
	}

	r := ffmpeg.ExecRunner{}
	return pipeline.Run(ctx, r, cfg)
}
