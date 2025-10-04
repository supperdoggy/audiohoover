# AudioHoover

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**AudioHoover** is a professional-grade Go tool designed to collect and consolidate audio files from multiple playlist directories into a single destination directory, with intelligent duplicate detection and filtering.

## Features

- 📁 **Multi-Playlist Support**: Process multiple playlist folders in one run
- 🔍 **Smart Duplicate Detection**: Skip files that already exist in the destination
- 🎵 **Audio File Filtering**: Only process specified audio file formats (configurable)
- 🏃 **Dry-Run Mode**: Preview operations without actually copying files
- 📊 **Detailed Statistics**: Track files processed, copied, skipped, and errors
- ⚙️ **Flexible Configuration**: Configure via CLI flags or environment variables
- 🧪 **Well-Tested**: Comprehensive unit and integration tests
- 📝 **Production Logging**: Structured logging with zap for production use

## Installation

### From Source

```bash
git clone https://github.com/supperdoggy/audiohoover.git
cd audiohoover
make install
```

### Using Go Install

```bash
go install github.com/supperdoggy/audiohoover@latest
```

### Build Binary

```bash
make build
```

This will create an `audiohoover` binary in the current directory.

## Usage

### Basic Usage

```bash
# Using default settings (./playlists -> ./output)
audiohoover

# Specify custom directories
audiohoover -playlists /path/to/playlists -destination /path/to/output

# Dry-run mode (preview without copying)
audiohoover -dry-run

# Verbose logging
audiohoover -verbose
```

### CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-playlists` | `./playlists` | Root directory containing playlist folders |
| `-destination` | `./output` | Destination folder for collected audio files |
| `-dry-run` | `false` | Run without actually copying files |
| `-extensions` | `.mp3,.flac,.wav,.m4a,.aac,.ogg,.wma` | Comma-separated list of audio file extensions |
| `-verbose` | `false` | Enable verbose logging |
| `-version` | - | Show version information |

### Environment Variables

You can also configure AudioHoover using environment variables:

```bash
export PLAYLISTS_ROOT=/path/to/playlists
export DESTINATION_FOLDER=/path/to/output
export DRY_RUN=true
export AUDIO_EXTENSIONS=.mp3,.wav,.flac
export VERBOSE=true

audiohoover
```

**Note**: CLI flags take precedence over environment variables.

### Examples

#### Example 1: Collect all MP3 and FLAC files

```bash
audiohoover \
  -playlists ~/Music/Playlists \
  -destination ~/Music/Collection \
  -extensions .mp3,.flac
```

#### Example 2: Preview operation (dry-run)

```bash
audiohoover -dry-run -verbose
```

#### Example 3: Process with custom extensions

```bash
audiohoover \
  -playlists ./my-playlists \
  -destination ./output \
  -extensions .mp3,.wav,.m4a,.aac
```

## Directory Structure

Expected directory structure:

```
playlists/
├── playlist1/
│   ├── song1.mp3
│   ├── song2.flac
│   └── cover.jpg (ignored)
├── playlist2/
│   ├── track1.wav
│   └── track2.mp3
└── playlist3/
    ├── audio1.m4a
    └── readme.txt (ignored)
```

After running AudioHoover:

```
output/
├── song1.mp3
├── song2.flac
├── track1.wav
├── track2.mp3
└── audio1.m4a
```

## Development

### Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile targets)

### Building

```bash
make build
```

### Testing

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Run all tests
make test-all

# Generate coverage report
make coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run go vet
make vet

# Run linter (requires golangci-lint)
make lint

# Run all checks
make check
```

### Makefile Targets

Run `make help` to see all available targets:

```bash
make help
```

Available targets:
- `build` - Build the binary
- `install` - Install to $GOPATH/bin
- `test` - Run unit tests
- `test-integration` - Run integration tests
- `test-all` - Run all tests
- `coverage` - Generate coverage report
- `lint` - Run linter
- `fmt` - Format code
- `vet` - Run go vet
- `clean` - Remove build artifacts
- `run` - Build and run
- `check` - Run all checks

## Project Structure

```
audiohoover/
├── main.go                      # Entry point with CLI
├── pkg/
│   ├── config/
│   │   ├── config.go           # Configuration management
│   │   └── config_test.go      # Config tests
│   ├── service/
│   │   ├── service.go          # Core business logic
│   │   └── service_test.go     # Service tests
│   └── utils/
│       ├── util.go             # Utility functions
│       └── util_test.go        # Utility tests
├── integration_test.go          # Integration tests
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## Configuration Details

### Audio Extensions

AudioHoover processes files based on their extensions. By default, it supports:

- `.mp3` - MP3 audio
- `.flac` - FLAC lossless audio
- `.wav` - WAV audio
- `.m4a` - MPEG-4 audio
- `.aac` - AAC audio
- `.ogg` - Ogg Vorbis
- `.wma` - Windows Media Audio

You can customize this list using the `-extensions` flag or `AUDIO_EXTENSIONS` environment variable.

### Duplicate Detection

AudioHoover detects duplicates based on filename only. If a file with the same name already exists in the destination folder, it will be skipped. This prevents overwriting existing files.

### Error Handling

- Non-fatal errors (e.g., individual file copy failures) are logged but don't stop processing
- Statistics report shows the number of errors encountered
- Exit code 1 is returned if any errors occurred
- Fatal errors (e.g., invalid configuration) stop execution immediately

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for new features
- Run `make check` before committing
- Follow Go best practices and idioms
- Add appropriate documentation

## License

The MIT License (MIT)

Copyright (c) 2015 Chris Kibble

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.