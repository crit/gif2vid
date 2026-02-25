package media

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/crit/gif2vid/internal/ffmpeg"
)

type mockRunner struct {
	ffmpeg.Runner
	mockRun func(ctx context.Context, name string, args []string) ([]byte, []byte, error)
}

func (m *mockRunner) Run(ctx context.Context, name string, args []string) ([]byte, []byte, error) {
	return m.mockRun(ctx, name, args)
}

func TestProbe(t *testing.T) {
	ctx := context.Background()

	t.Run("success video stream", func(t *testing.T) {
		mr := &mockRunner{
			mockRun: func(ctx context.Context, name string, args []string) ([]byte, []byte, error) {
				res := ProbeResult{
					Streams: []struct {
						CodecType string `json:"codec_type"`
						Width     int    `json:"width"`
						Height    int    `json:"height"`
					}{
						{CodecType: "video", Width: 100, Height: 200},
					},
				}
				b, _ := json.Marshal(res)
				return b, nil, nil
			},
		}
		w, h, err := Probe(ctx, mr, "test.gif")
		if err != nil {
			t.Fatalf("Probe failed: %v", err)
		}
		if w != 100 || h != 200 {
			t.Errorf("got %dx%d, want 100x200", w, h)
		}
	})

	t.Run("success image stream", func(t *testing.T) {
		// This test case simulates what happens when ffprobe returns an image stream.
		// Currently it might work because of the broad check in the loop,
		// but let's see if the -select_streams v:0 prevents it in real ffprobe.
		// In this mock, we only control what the runner returns.
		mr := &mockRunner{
			mockRun: func(ctx context.Context, name string, args []string) ([]byte, []byte, error) {
				res := ProbeResult{
					Streams: []struct {
						CodecType string `json:"codec_type"`
						Width     int    `json:"width"`
						Height    int    `json:"height"`
					}{
						{CodecType: "image2", Width: 300, Height: 400},
					},
				}
				b, _ := json.Marshal(res)
				return b, nil, nil
			},
		}
		w, h, err := Probe(ctx, mr, "test.webp")
		if err != nil {
			t.Fatalf("Probe failed: %v", err)
		}
		if w != 300 || h != 400 {
			t.Errorf("got %dx%d, want 300x400", w, h)
		}
	})

	t.Run("no streams", func(t *testing.T) {
		mr := &mockRunner{
			mockRun: func(ctx context.Context, name string, args []string) ([]byte, []byte, error) {
				res := ProbeResult{
					Streams: []struct {
						CodecType string `json:"codec_type"`
						Width     int    `json:"width"`
						Height    int    `json:"height"`
					}{},
				}
				b, _ := json.Marshal(res)
				return b, nil, nil
			},
		}
		_, _, err := Probe(ctx, mr, "test.gif")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("no dimensions triggers fallback", func(t *testing.T) {
		mr := &mockRunner{
			mockRun: func(ctx context.Context, name string, args []string) ([]byte, []byte, error) {
				if name == "ffprobe" && args[len(args)-1] != "test.gif" {
					res := ProbeResult{
						Streams: []struct {
							CodecType string `json:"codec_type"`
							Width     int    `json:"width"`
							Height    int    `json:"height"`
						}{
							{CodecType: "image2", Width: 700, Height: 800},
						},
					}
					b, _ := json.Marshal(res)
					return b, nil, nil
				}
				if name == "ffmpeg" {
					return nil, nil, nil
				}
				// Initial ffprobe returns no streams
				res := ProbeResult{Streams: nil}
				b, _ := json.Marshal(res)
				return b, nil, nil
			},
		}
		w, h, err := Probe(ctx, mr, "test.gif")
		if err != nil {
			t.Fatalf("Probe failed: %v", err)
		}
		if w != 700 || h != 800 {
			t.Errorf("got %dx%d, want 700x800", w, h)
		}
	})

	t.Run("ffprobe error triggers fallback", func(t *testing.T) {
		mr := &mockRunner{
			mockRun: func(ctx context.Context, name string, args []string) ([]byte, []byte, error) {
				if name == "ffprobe" && args[len(args)-1] != "test.gif" {
					// This is the probe on the temp file
					res := ProbeResult{
						Streams: []struct {
							CodecType string `json:"codec_type"`
							Width     int    `json:"width"`
							Height    int    `json:"height"`
						}{
							{CodecType: "image2", Width: 500, Height: 600},
						},
					}
					b, _ := json.Marshal(res)
					return b, nil, nil
				}
				if name == "ffmpeg" {
					return nil, nil, nil
				}
				// Initial ffprobe fails
				return nil, []byte("some error"), errors.New("failed")
			},
		}
		w, h, err := Probe(ctx, mr, "test.gif")
		if err != nil {
			t.Fatalf("Probe failed: %v", err)
		}
		if w != 500 || h != 600 {
			t.Errorf("got %dx%d, want 500x600", w, h)
		}
	})
}
