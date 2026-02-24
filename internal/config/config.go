package config

import (
	"errors"
	"flag"
	"runtime"
)

// Config holds all CLI/configuration options.
type Config struct {
	Output      string
	FPS         int
	CRF         int
	Preset      string
	BG          string
	Overwrite   bool
	KeepTemp    bool
	TmpDir      string
	Verbose     bool
	Concurrency int
	Inputs      []string
}

// AddFlags defines CLI flags on the provided FlagSet and returns a pointer to Config.
func AddFlags(fs *flag.FlagSet) *Config {
	cfg := &Config{}
	fs.StringVar(&cfg.Output, "output", "", "Output MP4 file path (required)")
	fs.StringVar(&cfg.Output, "o", "", "Output MP4 file path (required) [shorthand]")
	fs.IntVar(&cfg.FPS, "fps", 30, "Frames per second")
	fs.IntVar(&cfg.CRF, "crf", 23, "x264 CRF quality (lower is better)")
	fs.StringVar(&cfg.Preset, "preset", "medium", "x264 preset (ultrafast..placebo)")
	fs.StringVar(&cfg.BG, "bg", "black", "Background color (name or #RRGGBB)")
	fs.BoolVar(&cfg.Overwrite, "overwrite", false, "Overwrite output if it exists")
	fs.BoolVar(&cfg.KeepTemp, "keep-temp", false, "Keep temporary workspace")
	fs.StringVar(&cfg.TmpDir, "tmp-dir", "", "Temporary directory to use")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Verbose logging")
	fs.IntVar(&cfg.Concurrency, "concurrency", 0, "Number of parallel workers (default: runtime.NumCPU())")
	fs.IntVar(&cfg.Concurrency, "j", 0, "Number of parallel workers (default: runtime.NumCPU()) [shorthand]")
	return cfg
}

// Finalize validates required flags and attaches positional args as inputs.
func (c *Config) Finalize(args []string) error {
	c.Inputs = append(c.Inputs[:0], args...)
	if c.Output == "" {
		return errors.New("-o/--output is required")
	}
	if c.Concurrency <= 0 {
		c.Concurrency = runtime.NumCPU()
	}
	return nil
}
