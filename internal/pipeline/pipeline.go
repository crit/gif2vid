package pipeline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/crit/gif2vid/internal/concat"
	"github.com/crit/gif2vid/internal/config"
	"github.com/crit/gif2vid/internal/ffmpeg"
	"github.com/crit/gif2vid/internal/media"
)

// even rounds up to the nearest even number.
func even(x int) int {
	if x%2 != 0 {
		return x + 1
	}
	return x
}

// BuildFilter builds the ffmpeg -vf filter string.
func BuildFilter(cfg *config.Config, targetW, targetH int) string {
	return fmt.Sprintf("fps=%d,scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2:color=%s,format=yuv420p",
		cfg.FPS, targetW, targetH, targetW, targetH, cfg.BG)
}

// Run executes the full pipeline.
func Run(ctx context.Context, r ffmpeg.Runner, cfg *config.Config) error {
	// Probe inputs and compute target canvas
	maxW, maxH := 0, 0
	for _, in := range cfg.Inputs {
		w, h, err := media.Probe(ctx, r, in)
		if err != nil {
			return err
		}
		if w > maxW {
			maxW = w
		}
		if h > maxH {
			maxH = h
		}
	}
	maxW = even(maxW)
	maxH = even(maxH)
	if maxW == 0 || maxH == 0 {
		return fmt.Errorf("failed to determine target dimensions")
	}

	// Temp workspace
	tmpDir, err := filepath.Abs(cfg.TmpDir)
	if cfg.TmpDir == "" {
		// pick a default temp under OS temp; we won't add util.MkTempWorkspace to avoid extra deps here
		tmpDir, err = filepath.Abs(filepath.Join(os.TempDir(), "gif2vid-work"))
	}
	if err != nil {
		return err
	}
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return err
	}

	segments := make([]string, len(cfg.Inputs))
	type job struct {
		index int
		input string
	}
	jobs := make(chan job, len(cfg.Inputs))
	for i, in := range cfg.Inputs {
		jobs <- job{index: i, input: in}
	}
	close(jobs)

	var wg sync.WaitGroup
	errs := make(chan error, len(cfg.Inputs))
	numWorkers := cfg.Concurrency
	if numWorkers > len(cfg.Inputs) {
		numWorkers = len(cfg.Inputs)
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				seg := filepath.Join(tmpDir, fmt.Sprintf("seg_%04d.mp4", j.index))
				args := []string{
					"-y", // segments may overwrite if re-run within workspace
					"-i", j.input,
					"-vf", BuildFilter(cfg, maxW, maxH),
					"-an",
					"-c:v", "libx264",
					"-preset", cfg.Preset,
					"-crf", fmt.Sprintf("%d", cfg.CRF),
					seg,
				}
				_, stderr, err := r.Run(ctx, "ffmpeg", args)
				if err != nil {
					errs <- fmt.Errorf("ffmpeg segment failed for %s:\ncmd: %s\n%s", j.input, ffmpeg.PrettyCmd("ffmpeg", args), string(stderr))
					return
				}
				segments[j.index] = seg
			}
		}()
	}
	wg.Wait()
	close(errs)

	if err, ok := <-errs; ok {
		return err
	}

	concatPath := filepath.Join(tmpDir, "concat.txt")
	if err := concat.WriteConcatFile(concatPath, segments); err != nil {
		return err
	}
	outTmp := filepath.Join(tmpDir, "out.tmp.mp4")
	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", concatPath,
		"-c:v", "libx264",
		"-preset", cfg.Preset,
		"-crf", fmt.Sprintf("%d", cfg.CRF),
		"-pix_fmt", "yuv420p",
		"-movflags", "+faststart",
		"-an",
		outTmp,
	}
	_, stderr, err := r.Run(ctx, "ffmpeg", args)
	if err != nil {
		return fmt.Errorf("ffmpeg concat failed:\ncmd: %s\n%s", ffmpeg.PrettyCmd("ffmpeg", args), string(stderr))
	}

	// Move to final output
	finalArgs := []string{"-y"}
	_ = finalArgs // not used; we use Go to move
	if !cfg.Overwrite {
		if _, err := os.Stat(cfg.Output); err == nil {
			return fmt.Errorf("output exists: %s (use --overwrite)", cfg.Output)
		}
	}
	if err := os.MkdirAll(filepath.Dir(cfg.Output), 0o755); err != nil {
		return err
	}
	if err := os.Rename(outTmp, cfg.Output); err != nil {
		return err
	}

	// Cleanup unless keep-temp
	if !cfg.KeepTemp {
		_ = os.RemoveAll(tmpDir)
	} else {
		fmt.Printf("[gif2vid] temp kept at: %s\n", tmpDir)
	}
	return nil
}
