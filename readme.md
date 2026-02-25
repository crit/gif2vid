# gif2vid

`gif2vid` is a Go-based CLI tool that takes a directory of GIF and/or animated WebP files and combines them into a single H.264 MP4 video. It uses `ffmpeg` and `ffprobe` under the hood to handle media processing.

## Features

- **Multi-format Support**: Combine GIF and animated WebP files into one video.
- **Automatic Sizing**: Automatically calculates the maximum width and height across all input files to create a uniform canvas (rounded up to the nearest even number for H.264 compatibility).
- **Contain Fit**: Each input is scaled to fit the target dimensions without cropping, with configurable background padding (default: black).
- **High Compatibility**: Generates H.264 MP4 files with `yuv420p` pixel format and `+faststart` for broad device and web compatibility.
- **Configurable Quality**: Control frame rate (FPS), quality (CRF), and encoding speed (preset).
- **Deterministic Order**: The output video follows the exact order of the files provided in the command line.

## Prerequisites

- **FFmpeg**: Must be installed and available in your `PATH`.
- **FFprobe**: Must be installed and available in your `PATH`.

On macOS (using Homebrew):
```bash
brew install ffmpeg
```

## Installation

### Using Make

```bash
make install
```

This will run `go install ./cmd/gif2vid`.

### Manual Build

```bash
go build -o gif2vid ./cmd/gif2vid
```

## Usage

```bash
gif2vid [flags] <input_directory>
```

### Examples

**Basic Usage:**
Combine all GIFs in a directory into an output video.
```bash
gif2vid -o output.mp4 ./my_gifs
```

**Mixed Formats and Custom Settings:**
Combine GIF and WebP files from a directory with a specific frame rate and background color.
```bash
gif2vid -o result.mp4 --fps 60 --bg "#1a1a1a" --crf 18 ./media_dir
```

**Overwrite Existing File:**
```bash
gif2vid -o output.mp4 --overwrite ./input_dir
```

### Flags

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-o`, `--output` | **(Required)** Output MP4 file path. | |
| `--fps` | Frames per second for the output video. | `30` |
| `--crf` | x264 CRF quality (lower is better, typically 0–51). | `23` |
| `--preset` | x264 encoding preset (`ultrafast` to `placebo`). | `medium` |
| `--bg` | Background padding color (name or #RRGGBB). | `black` |
| `--overwrite` | Overwrite the output file if it already exists. | `false` |
| `--keep-temp` | Retain the temporary workspace for debugging. | `false` |
| `--tmp-dir` | Specify a custom temporary directory. | (OS temp) |
| `--concurrency`, `-j` | Number of parallel workers (segments generation). | (Num CPUs) |
| `--verbose` | Enable verbose logging. | `false` |

## Development

### Running Tests
```bash
make test
```

### Building for Debugging
```bash
make build-debug
```

## License
MIT (or as specified in the repository)
