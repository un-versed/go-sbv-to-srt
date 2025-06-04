# SBV to SRT Converter
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![coverage](https://img.shields.io/badge/coverage-78.2%25-brightgreen)

A command-line tool written in Go to convert SBV (SubViewer) subtitle files to SRT (SubRip) format.

## Features

- ✅ Convert SBV files to SRT format
- ✅ Automatic output file naming (when no output path is specified)
- ✅ Comprehensive input validation and error handling
- ✅ Cross-platform support (Linux, Windows, macOS)
- ✅ CLI with Cobra framework for robust argument parsing
- ✅ Built with Go 1.24.3 and modern Go practices
- ✅ Comprehensive test coverage with unit and integration tests

## Installation

### From Source

```bash
git clone https://github.com/un-versed/go-sbv-to-srt.git
cd go-sbv-to-srt
go build -o go-sbv-to-srt
```

### Using Go Install

```bash
go install github.com/un-versed/go-sbv-to-srt@latest
```

## Dependencies

This project uses minimal external dependencies:
- **[Cobra](https://github.com/spf13/cobra)** v1.9.1 - CLI framework for robust command-line interface
- **Go standard library** - For file I/O, string processing, and time handling

No additional runtime dependencies are required.

## Usage

### Basic Usage

```bash
# Convert SBV file (output will be created in same directory with .srt extension)
go-sbv-to-srt -i input.sbv

# Specify custom output path
go-sbv-to-srt -i input.sbv -o output.srt

# Using long flags
go-sbv-to-srt --input video.sbv --output subtitles.srt
```

### Command Line Options

- `-i, --input`: Input SBV file path (required)
- `-o, --output`: Output SRT file path (optional)
- `-h, --help`: Show help information
- `version`: Show version information
- `completion`: Generate shell completion scripts

### Examples

```bash
# Convert subtitle.sbv to subtitle.srt in the same directory
go-sbv-to-srt -i subtitle.sbv

# Convert with custom output location
go-sbv-to-srt -i ./videos/movie.sbv -o ./subtitles/movie.srt

# Show help
go-sbv-to-srt --help

# Show version information
go-sbv-to-srt version

# Generate bash completion
go-sbv-to-srt completion bash > /etc/bash_completion.d/go-sbv-to-srt

# Generate zsh completion
go-sbv-to-srt completion zsh > "${fpath[1]}/_go-sbv-to-srt"
```

## File Format Support

### Input Format (SBV)
SBV (SubViewer) is a subtitle format commonly used by YouTube and other video platforms. It uses timestamps in the format `HH:MM:SS.mmm,HH:MM:SS.mmm` followed by subtitle text.

**Example SBV format:**
```
0:00:01.000,0:00:04.000
This is a sample subtitle
that spans multiple lines.

0:00:05.500,0:00:08.200
This is another subtitle entry.
```

### Output Format (SRT)
SRT (SubRip) is a widely supported subtitle format that can be used with most media players and video editing software.

**Example SRT format:**
```
1
00:00:01,000 --> 00:00:04,000
This is a sample subtitle
that spans multiple lines.

2
00:00:05,500 --> 00:00:08,200
This is another subtitle entry.
```

### Conversion Features

- **Accurate timestamp conversion**: Handles SBV's colon-separated format to SRT's arrow-separated format
- **Multi-line text preservation**: Maintains original line breaks and formatting
- **Sequential numbering**: Automatically generates SRT sequence numbers
- **Error handling**: Validates timestamp formats and provides helpful error messages

## Development

### Prerequisites

- Go 1.24.3 or later
- Make (optional, for using Makefile)
- Git (for version information)

### Quick Start

```bash
# Clone and build
git clone https://github.com/un-versed/go-sbv-to-srt.git
cd go-sbv-to-srt
make build

# Or use the build script (includes version info)
./build.sh

# Or use Go directly
go build -o go-sbv-to-srt
```

### Available Make Targets

- `make build` - Build the application
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report
- `make test-race` - Run tests with race detection
- `make clean` - Clean build artifacts
- `make deps` - Install dependencies
- `make fmt` - Format code
- `make vet` - Vet code for issues
- `make lint` - Lint code (requires golangci-lint)
- `make build-all` - Build for multiple platforms (Linux, Windows, macOS)
- `make run-sample` - Run with sample test data
- `make dev` - Complete development workflow (deps, fmt, vet, test, build)
- `make help` - Show all available targets

## CI/CD

The project includes comprehensive GitHub Actions workflows:

- **CI Workflow** (`.github/workflows/ci.yml`) - Runs on every push and PR:
  - Tests on Go 1.24.x
  - Cross-platform testing (Linux, Windows, macOS)
  - Code quality checks (formatting, vetting, linting)
  - Test coverage reporting

- **Build Workflow** (`.github/workflows/build.yml`) - Automated builds:
  - Multi-platform binary generation
  - Artifact uploads for releases

- **Release Workflow** (`.github/workflows/release.yml`) - Automated releases:
  - Tagged release creation
  - Cross-platform binary distribution
  - Release notes generation

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed development guidelines, including:

- Project structure and architecture
- Development workflow and best practices
- Testing guidelines and coverage requirements
- Code style and formatting standards
- Available Make targets for development tasks

### Quick Contributing Steps

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes and add tests
4. Run the full test suite (`make test`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- [x] CLI scaffolding with proper argument parsing (Cobra framework)
- [x] Input/output validation with comprehensive error handling  
- [x] SBV file parsing with robust timestamp handling
- [x] SRT file generation with proper formatting
- [x] Comprehensive test coverage (unit and integration tests)
- [x] Cross-platform build support and automation
- [x] Development tooling and Make targets
- [x] GitHub Actions CI/CD pipeline (build, test, release workflows)
- [x] Shell completion support (bash, zsh, fish, PowerShell)
- [x] Version command and build information
- [ ] Batch conversion support (convert multiple files at once)
