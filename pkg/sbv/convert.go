// Package sbv provides functionality to convert SBV (YouTube SubViewer) subtitle files to SRT (SubRip) format.
package sbv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Subtitle represents a single subtitle entry with timing and text content.
type Subtitle struct {
	StartTime time.Duration
	EndTime   time.Duration
	Text      string
}

// Converter defines the interface for converting SBV files to SRT format.
type Converter interface {
	// ParseFromFile reads and parses an SBV file from the given file path.
	// Returns a slice of Subtitle entries or an error if parsing fails.
	ParseFromFile(filename string) ([]Subtitle, error)

	// ParseFromReader reads and parses SBV content from an io.Reader.
	// Returns a slice of Subtitle entries or an error if parsing fails.
	ParseFromReader(reader io.Reader) ([]Subtitle, error)

	// ConvertToSRT converts parsed subtitles to SRT format string.
	// Takes a slice of Subtitle entries and returns the SRT formatted string.
	ConvertToSRT(subtitles []Subtitle) string

	// WriteToFile converts subtitles and writes them directly to an SRT file.
	// Takes subtitles and output filename, returns error if write fails.
	WriteToFile(subtitles []Subtitle, filename string) error

	// WriteToWriter converts subtitles and writes them to an io.Writer.
	// Takes subtitles and writer, returns error if write fails.
	WriteToWriter(subtitles []Subtitle, writer io.Writer) error
}

// DefaultConverter is the standard implementation of the Converter interface.
type DefaultConverter struct {
}

// NewConverter creates a new instance of DefaultConverter.
func NewConverter() *DefaultConverter {
	return &DefaultConverter{}
}

// ConvertToSRT converts parsed subtitles to SRT format string.
func (c *DefaultConverter) ConvertToSRT(subtitles []Subtitle) string {
	var result strings.Builder
	result.Grow(len(subtitles) * 100) // Pre-allocate approximate capacity

	for i, subtitle := range subtitles {
		// SRT sequence number (1-based)
		result.WriteString(strconv.Itoa(i + 1))
		result.WriteByte('\n')

		// SRT timestamp format: HH:MM:SS,mmm --> HH:MM:SS,mmm
		startTime := c.formatSRTTime(subtitle.StartTime)
		endTime := c.formatSRTTime(subtitle.EndTime)
		result.WriteString(startTime)
		result.WriteString(" --> ")
		result.WriteString(endTime)
		result.WriteByte('\n')

		// Subtitle text
		result.WriteString(subtitle.Text)
		result.WriteString("\n\n")
	}

	return result.String()
}

// ParseFromFile reads and parses an SBV file from the given file path.
func (c *DefaultConverter) ParseFromFile(filename string) ([]Subtitle, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error but don't override the main error
			fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", closeErr)
		}
	}()

	return c.ParseFromReader(file)
}

// ParseFromReader reads and parses SBV content from an io.Reader.
func (c *DefaultConverter) ParseFromReader(reader io.Reader) ([]Subtitle, error) {
	var subtitles []Subtitle
	scanner := bufio.NewScanner(reader)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check if this line contains timestamps
		if c.isTimestampLine(line) {
			subtitle, nextIndex, err := c.parseSubtitleBlock(lines, i)
			if err != nil {
				return nil, fmt.Errorf("failed to parse subtitle block: %w", err)
			}
			subtitles = append(subtitles, subtitle)
			i = nextIndex - 1 // -1 because the loop will increment
		}
	}

	return subtitles, nil
}

// WriteToFile converts subtitles and writes them directly to an SRT file.
func (c *DefaultConverter) WriteToFile(subtitles []Subtitle, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error but don't override the main error
			fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", closeErr)
		}
	}()

	return c.WriteToWriter(subtitles, file)
}

// WriteToWriter converts subtitles and writes them to an io.Writer.
func (c *DefaultConverter) WriteToWriter(subtitles []Subtitle, writer io.Writer) error {
	srtContent := c.ConvertToSRT(subtitles)
	_, err := writer.Write([]byte(srtContent))
	if err != nil {
		return fmt.Errorf("failed to write SRT content: %w", err)
	}
	return nil
}

// formatSRTTime formats a time.Duration to SRT timestamp format (HH:MM:SS,mmm).
func (c *DefaultConverter) formatSRTTime(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	milliseconds := int(duration.Nanoseconds()/1_000_000) % 1000

	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, milliseconds)
}

// isTimestampLine checks if a line contains SBV timestamp format.
func (c *DefaultConverter) isTimestampLine(line string) bool {
	// Simple check: contains comma and colon (timestamp indicators)
	return strings.ContainsRune(line, ',') && strings.ContainsRune(line, ':')
}

// parseSubtitleBlock parses a single subtitle block starting with a timestamp line.
func (c *DefaultConverter) parseSubtitleBlock(lines []string, startIndex int) (Subtitle, int, error) {
	if startIndex >= len(lines) {
		return Subtitle{}, startIndex, fmt.Errorf("start index out of bounds")
	}

	timestampLine := lines[startIndex]
	startTime, endTime, err := c.parseTimestamps(timestampLine)
	if err != nil {
		return Subtitle{}, startIndex, fmt.Errorf("failed to parse timestamps: %w", err)
	}

	// Read subtitle text (can be multiple lines)
	var textLines []string
	currentIndex := startIndex + 1

	for currentIndex < len(lines) && lines[currentIndex] != "" {
		textLines = append(textLines, lines[currentIndex])
		currentIndex++
	}

	text := strings.Join(textLines, "\n")

	return Subtitle{
		StartTime: startTime,
		EndTime:   endTime,
		Text:      text,
	}, currentIndex, nil
}

// parseTimestamps parses SBV timestamp format "H:MM:SS.mmm,H:MM:SS.mmm"
func (c *DefaultConverter) parseTimestamps(timestampLine string) (time.Duration, time.Duration, error) {
	parts := strings.Split(timestampLine, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid timestamp format: %s", timestampLine)
	}

	startTime, err := c.parseTime(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse start time: %w", err)
	}

	endTime, err := c.parseTime(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse end time: %w", err)
	}

	return startTime, endTime, nil
}

// parseTime parses a time string in format "H:MM:SS.mmm"
func (c *DefaultConverter) parseTime(timeStr string) (time.Duration, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hours: %s", parts[0])
	}
	if hours < 0 || hours > 23 {
		return 0, fmt.Errorf("hours out of range (0-23): %d", hours)
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %s", parts[1])
	}
	if minutes < 0 || minutes > 59 {
		return 0, fmt.Errorf("minutes out of range (0-59): %d", minutes)
	}

	// Handle seconds and milliseconds
	secondsParts := strings.Split(parts[2], ".")
	if len(secondsParts) != 2 {
		return 0, fmt.Errorf("invalid seconds format: %s", parts[2])
	}

	seconds, err := strconv.Atoi(secondsParts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds: %s", secondsParts[0])
	}
	if seconds < 0 || seconds > 59 {
		return 0, fmt.Errorf("seconds out of range (0-59): %d", seconds)
	}

	milliseconds, err := strconv.Atoi(secondsParts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid milliseconds: %s", secondsParts[1])
	}
	if milliseconds < 0 || milliseconds > 999 {
		return 0, fmt.Errorf("milliseconds out of range (0-999): %d", milliseconds)
	}

	totalDuration := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second +
		time.Duration(milliseconds)*time.Millisecond

	return totalDuration, nil
}
