# SBV Package Documentation

The `sbv` package provides functionality to convert SBV (YouTube SubViewer) subtitle files to SRT (SubRip) format.

## Features

- Clean interface-based design
- Multiple input methods (files and io.Reader)  
- Multiple output methods (string, files, io.Writer)
- Robust parsing with multi-line subtitle support
- Idiomatic Go error handling

## Usage

```go
package main

import (
    "fmt"
    "github.com/un-versed/go-sbv-to-srt/pkg/sbv"
)

func main() {
    converter := sbv.NewConverter()
    
    // Parse SBV file
    subtitles, err := converter.ParseFromFile("input.sbv")
    if err != nil {
        panic(err)
    }
    
    // Convert and save as SRT
    err = converter.WriteToFile(subtitles, "output.srt")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Converted %d subtitles successfully!\n", len(subtitles))
}
```

## Interface

```go
type Converter interface {
    ParseFromFile(filename string) ([]Subtitle, error)
    ParseFromReader(reader io.Reader) ([]Subtitle, error)
    ConvertToSRT(subtitles []Subtitle) string
    WriteToFile(subtitles []Subtitle, filename string) error
    WriteToWriter(subtitles []Subtitle, writer io.Writer) error
}
```

## Testing

Run tests with: `go test ./pkg/sbv/... -v`
