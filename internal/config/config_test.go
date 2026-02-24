package config

import (
	"flag"
	"testing"
)

func TestConfigFinalize(t *testing.T) {
	t.Run("missing output", func(t *testing.T) {
		cfg := &Config{}
		err := cfg.Finalize([]string{"a.gif"})
		if err == nil {
			t.Error("Finalize should fail without output")
		}
	})

	t.Run("valid", func(t *testing.T) {
		cfg := &Config{Output: "out.mp4"}
		err := cfg.Finalize([]string{"a.gif", "b.webp"})
		if err != nil {
			t.Errorf("Finalize failed: %v", err)
		}
		if len(cfg.Inputs) != 2 {
			t.Errorf("expected 2 inputs, got %d", len(cfg.Inputs))
		}
		if cfg.Inputs[0] != "a.gif" || cfg.Inputs[1] != "b.webp" {
			t.Errorf("inputs mismatch: %v", cfg.Inputs)
		}
	})
}

func TestAddFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := AddFlags(fs)
	err := fs.Parse([]string{"-o", "out.mp4", "--fps", "60", "input.gif"})
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
	if len(cfg.Inputs) != 1 || cfg.Inputs[0] != "input.gif" {
		t.Errorf("Inputs = %v; want [input.gif]", cfg.Inputs)
	}
}

func TestAddFlagsConcurrency(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := AddFlags(fs)
	err := fs.Parse([]string{"-o", "out.mp4", "-j", "4", "input.gif"})
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
