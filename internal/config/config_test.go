package config

import (
	"flag"
	"testing"
)

func TestConfigFinalize(t *testing.T) {
	t.Run("missing output", func(t *testing.T) {
		cfg := &Config{}
		err := cfg.Finalize([]string{"indir"})
		if err == nil {
			t.Error("Finalize should fail without output")
		}
	})

	t.Run("missing input directory", func(t *testing.T) {
		cfg := &Config{Output: "out.mp4"}
		err := cfg.Finalize([]string{})
		if err == nil {
			t.Error("Finalize should fail without input directory")
		}
	})

	t.Run("too many input arguments", func(t *testing.T) {
		cfg := &Config{Output: "out.mp4"}
		err := cfg.Finalize([]string{"dir1", "dir2"})
		if err == nil {
			t.Error("Finalize should fail with multiple input arguments")
		}
	})

	t.Run("valid", func(t *testing.T) {
		cfg := &Config{Output: "out.mp4"}
		err := cfg.Finalize([]string{"indir"})
		if err != nil {
			t.Errorf("Finalize failed: %v", err)
		}
		if cfg.InputDir != "indir" {
			t.Errorf("expected InputDir 'indir', got %q", cfg.InputDir)
		}
	})
}

func TestAddFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := AddFlags(fs)
	err := fs.Parse([]string{"-o", "out.mp4", "--fps", "60", "indir"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if err := cfg.Finalize(fs.Args()); err != nil {
		t.Fatalf("Finalize failed: %v", err)
	}
	if cfg.Output != "out.mp4" {
		t.Errorf("Output = %q; want out.mp4", cfg.Output)
	}
	if cfg.FPS != 60 {
		t.Errorf("FPS = %d; want 60", cfg.FPS)
	}
	if cfg.Concurrency <= 0 {
		t.Errorf("Concurrency = %d; want > 0 (default)", cfg.Concurrency)
	}
	if cfg.InputDir != "indir" {
		t.Errorf("InputDir = %q; want indir", cfg.InputDir)
	}
}

func TestAddFlagsConcurrency(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := AddFlags(fs)
	err := fs.Parse([]string{"-o", "out.mp4", "-j", "4", "indir"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if err := cfg.Finalize(fs.Args()); err != nil {
		t.Fatalf("Finalize failed: %v", err)
	}
	if cfg.Concurrency != 4 {
		t.Errorf("Concurrency = %d; want 4", cfg.Concurrency)
	}
}
