package sbv

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewConverter(t *testing.T) {
	converter := NewConverter()
	if converter == nil {
		t.Fatal("NewConverter() returned nil")
	}
}

func TestParseTimestamps(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name      string
		input     string
		wantStart time.Duration
		wantEnd   time.Duration
		wantErr   bool
	}{
		{
			name:      "valid timestamp",
			input:     "0:00:01.000,0:00:04.000",
			wantStart: 1 * time.Second,
			wantEnd:   4 * time.Second,
			wantErr:   false,
		},
		{
			name:      "timestamp with hours",
			input:     "1:30:15.500,1:30:20.750",
			wantStart: 1*time.Hour + 30*time.Minute + 15*time.Second + 500*time.Millisecond,
			wantEnd:   1*time.Hour + 30*time.Minute + 20*time.Second + 750*time.Millisecond,
			wantErr:   false,
		},
		{
			name:    "invalid format - no comma",
			input:   "0:00:01.000 0:00:04.000",
			wantErr: true,
		},
		{
			name:    "invalid format - multiple commas",
			input:   "0:00:01.000,0:00:04.000,0:00:06.000",
			wantErr: true,
		},
		{
			name:    "invalid format - empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid start time format",
			input:   "invalid:00:01.000,0:00:04.000",
			wantErr: true,
		},
		{
			name:    "invalid end time format",
			input:   "0:00:01.000,invalid:00:04.000",
			wantErr: true,
		},
		{
			name:    "start time out of range",
			input:   "0:60:01.000,0:00:04.000",
			wantErr: true,
		},
		{
			name:    "end time out of range",
			input:   "0:00:01.000,0:00:70.000",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := converter.parseTimestamps(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseTimestamps() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("parseTimestamps() unexpected error: %v", err)
				return
			}

			if start != tt.wantStart {
				t.Errorf("parseTimestamps() start time = %v, want %v", start, tt.wantStart)
			}

			if end != tt.wantEnd {
				t.Errorf("parseTimestamps() end time = %v, want %v", end, tt.wantEnd)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name        string
		input       string
		want        time.Duration
		wantErr     bool
		description string
	}{
		{
			name:        "valid time with zero values",
			input:       "0:00:00.000",
			want:        0,
			wantErr:     false,
			description: "should parse zero time correctly",
		},
		{
			name:        "valid time with seconds",
			input:       "0:00:30.500",
			want:        30*time.Second + 500*time.Millisecond,
			wantErr:     false,
			description: "should parse seconds and milliseconds",
		},
		{
			name:        "valid time with hours and minutes",
			input:       "2:15:45.250",
			want:        2*time.Hour + 15*time.Minute + 45*time.Second + 250*time.Millisecond,
			wantErr:     false,
			description: "should parse hours, minutes, seconds, and milliseconds",
		},
		{
			name:        "invalid format - too few parts",
			input:       "00:30.500",
			wantErr:     true,
			description: "should reject format with only 2 colon-separated parts",
		},
		{
			name:        "invalid format - too many parts",
			input:       "1:2:3:4.500",
			wantErr:     true,
			description: "should reject format with more than 3 colon-separated parts",
		},
		{
			name:        "invalid hours - non-numeric",
			input:       "abc:00:30.500",
			wantErr:     true,
			description: "should reject non-numeric hours",
		},
		{
			name:        "invalid hours - out of range high",
			input:       "25:00:30.500",
			wantErr:     true,
			description: "should reject hours greater than 23",
		},
		{
			name:        "invalid hours - negative",
			input:       "-1:00:30.500",
			wantErr:     true,
			description: "should reject negative hours",
		},
		{
			name:        "invalid minutes - non-numeric",
			input:       "1:xy:30.500",
			wantErr:     true,
			description: "should reject non-numeric minutes",
		},
		{
			name:        "invalid minutes - out of range high",
			input:       "1:60:30.500",
			wantErr:     true,
			description: "should reject minutes greater than 59",
		},
		{
			name:        "invalid minutes - negative",
			input:       "1:-5:30.500",
			wantErr:     true,
			description: "should reject negative minutes",
		},
		{
			name:        "invalid seconds format - no decimal point",
			input:       "1:30:45500",
			wantErr:     true,
			description: "should reject seconds without decimal point",
		},
		{
			name:        "invalid seconds format - multiple decimal points",
			input:       "1:30:45.5.00",
			wantErr:     true,
			description: "should reject seconds with multiple decimal points",
		},
		{
			name:        "invalid seconds - non-numeric",
			input:       "1:30:ab.500",
			wantErr:     true,
			description: "should reject non-numeric seconds",
		},
		{
			name:        "invalid seconds - out of range high",
			input:       "1:30:60.500",
			wantErr:     true,
			description: "should reject seconds greater than 59",
		},
		{
			name:        "invalid seconds - negative",
			input:       "1:30:-5.500",
			wantErr:     true,
			description: "should reject negative seconds",
		},
		{
			name:        "invalid milliseconds - non-numeric",
			input:       "1:30:45.xyz",
			wantErr:     true,
			description: "should reject non-numeric milliseconds",
		},
		{
			name:        "invalid milliseconds - out of range high",
			input:       "1:30:45.1000",
			wantErr:     true,
			description: "should reject milliseconds greater than 999",
		},
		{
			name:        "invalid milliseconds - negative",
			input:       "1:30:45.-100",
			wantErr:     true,
			description: "should reject negative milliseconds",
		},
		{
			name:        "empty string",
			input:       "",
			wantErr:     true,
			description: "should reject empty string",
		},
		{
			name:        "only colons",
			input:       "::",
			wantErr:     true,
			description: "should reject string with only colons",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.parseTime(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseTime() expected error for %s, got nil", tt.description)
				}
				return
			}

			if err != nil {
				t.Errorf("parseTime() unexpected error for %s: %v", tt.description, err)
				return
			}

			if result != tt.want {
				t.Errorf("parseTime() %s: got %v, want %v", tt.description, result, tt.want)
			}
		})
	}
}

func TestParseFromReader(t *testing.T) {
	converter := NewConverter()

	sbvContent := `0:00:01.000,0:00:04.000
This is a sample SBV subtitle file
used for testing purposes.

0:00:05.500,0:00:08.200
SBV files use this timestamp format
and are commonly used by YouTube.

0:00:10.000,0:00:12.500
This will be converted to SRT format.`

	reader := strings.NewReader(sbvContent)
	subtitles, err := converter.ParseFromReader(reader)

	if err != nil {
		t.Fatalf("ParseFromReader() error: %v", err)
	}

	expectedCount := 3
	if len(subtitles) != expectedCount {
		t.Errorf("ParseFromReader() got %d subtitles, want %d", len(subtitles), expectedCount)
	}

	// Test first subtitle
	if subtitles[0].StartTime != 1*time.Second {
		t.Errorf("First subtitle start time = %v, want %v", subtitles[0].StartTime, 1*time.Second)
	}

	if subtitles[0].EndTime != 4*time.Second {
		t.Errorf("First subtitle end time = %v, want %v", subtitles[0].EndTime, 4*time.Second)
	}

	expectedText := "This is a sample SBV subtitle file\nused for testing purposes."
	if subtitles[0].Text != expectedText {
		t.Errorf("First subtitle text = %q, want %q", subtitles[0].Text, expectedText)
	}
}

func TestConvertToSRT(t *testing.T) {
	converter := NewConverter()

	subtitles := []Subtitle{
		{
			StartTime: 1 * time.Second,
			EndTime:   4 * time.Second,
			Text:      "First subtitle",
		},
		{
			StartTime: 5*time.Second + 500*time.Millisecond,
			EndTime:   8*time.Second + 200*time.Millisecond,
			Text:      "Second subtitle\nwith multiple lines",
		},
	}

	result := converter.ConvertToSRT(subtitles)

	expected := `1
00:00:01,000 --> 00:00:04,000
First subtitle

2
00:00:05,500 --> 00:00:08,200
Second subtitle
with multiple lines

`

	if result != expected {
		t.Errorf("ConvertToSRT() = %q, want %q", result, expected)
	}
}

func TestFormatSRTTime(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: "00:00:00,000",
		},
		{
			name:     "1 second",
			duration: 1 * time.Second,
			expected: "00:00:01,000",
		},
		{
			name:     "1 hour 30 minutes 15.5 seconds",
			duration: 1*time.Hour + 30*time.Minute + 15*time.Second + 500*time.Millisecond,
			expected: "01:30:15,500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.formatSRTTime(tt.duration)
			if result != tt.expected {
				t.Errorf("formatSRTTime() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFullConversion(t *testing.T) {
	converter := NewConverter()

	// Test with the sample SBV content
	sbvContent := `0:00:01.000,0:00:04.000
This is a sample SBV subtitle file
used for testing purposes.

0:00:05.500,0:00:08.200
SBV files use this timestamp format
and are commonly used by YouTube.

0:00:10.000,0:00:12.500
This will be converted to SRT format.`

	reader := strings.NewReader(sbvContent)
	subtitles, err := converter.ParseFromReader(reader)
	if err != nil {
		t.Fatalf("ParseFromReader() error: %v", err)
	}

	srtOutput := converter.ConvertToSRT(subtitles)

	// Verify the output contains expected elements
	expectedElements := []string{
		"1\n00:00:01,000 --> 00:00:04,000",
		"2\n00:00:05,500 --> 00:00:08,200",
		"3\n00:00:10,000 --> 00:00:12,500",
		"This is a sample SBV subtitle file\nused for testing purposes.",
		"SBV files use this timestamp format\nand are commonly used by YouTube.",
		"This will be converted to SRT format.",
	}

	for _, element := range expectedElements {
		if !strings.Contains(srtOutput, element) {
			t.Errorf("SRT output missing expected element: %q", element)
		}
	}
}

func TestParseFromFile(t *testing.T) {
	converter := NewConverter()

	// Create temporary test files
	validSBV := `0:00:01.000,0:00:04.000
First subtitle line
Second subtitle line

0:00:05.500,0:00:08.200
Another subtitle`

	invalidSBV := `invalid timestamp format
This is not a valid SBV file`

	tests := []struct {
		name        string
		content     string
		setupFile   bool
		wantCount   int
		wantErr     bool
		description string
	}{
		{
			name:        "valid SBV file",
			content:     validSBV,
			setupFile:   true,
			wantCount:   2,
			wantErr:     false,
			description: "should parse valid SBV file with multiple subtitle blocks",
		},
		{
			name:        "invalid SBV content",
			content:     invalidSBV,
			setupFile:   true,
			wantCount:   0,
			wantErr:     false, // No valid timestamps found, but not an error
			description: "should handle invalid SBV content gracefully",
		},
		{
			name:        "non-existent file",
			content:     "",
			setupFile:   false,
			wantCount:   0,
			wantErr:     true,
			description: "should return error when file does not exist",
		},
		{
			name:        "empty file",
			content:     "",
			setupFile:   true,
			wantCount:   0,
			wantErr:     false,
			description: "should handle empty file without error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filename string
			if tt.setupFile {
				// Create temporary file
				tmpFile, err := os.CreateTemp("", "test_*.sbv")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				filename = tmpFile.Name()
				defer os.Remove(filename)

				if tt.content != "" {
					if _, err := tmpFile.WriteString(tt.content); err != nil {
						tmpFile.Close()
						t.Fatalf("Failed to write test content: %v", err)
					}
				}
				tmpFile.Close()
			} else {
				filename = "non_existent_file.sbv"
			}

			subtitles, err := converter.ParseFromFile(filename)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseFromFile() expected error for %s, got nil", tt.description)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseFromFile() unexpected error for %s: %v", tt.description, err)
				return
			}

			if len(subtitles) != tt.wantCount {
				t.Errorf("ParseFromFile() %s: got %d subtitles, want %d", tt.description, len(subtitles), tt.wantCount)
			}
		})
	}
}

func TestWriteToFile(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name        string
		subtitles   []Subtitle
		wantErr     bool
		description string
	}{
		{
			name: "single subtitle",
			subtitles: []Subtitle{
				{
					StartTime: 1 * time.Second,
					EndTime:   4 * time.Second,
					Text:      "Single subtitle",
				},
			},
			wantErr:     false,
			description: "should write single subtitle to SRT file",
		},
		{
			name: "multiple subtitles",
			subtitles: []Subtitle{
				{
					StartTime: 1 * time.Second,
					EndTime:   4 * time.Second,
					Text:      "First subtitle",
				},
				{
					StartTime: 5 * time.Second,
					EndTime:   8 * time.Second,
					Text:      "Second subtitle\nwith multiple lines",
				},
			},
			wantErr:     false,
			description: "should write multiple subtitles with multiline text to SRT file",
		},
		{
			name:        "empty subtitles",
			subtitles:   []Subtitle{},
			wantErr:     false,
			description: "should handle empty subtitle slice without error",
		},
		{
			name: "subtitle with complex timing",
			subtitles: []Subtitle{
				{
					StartTime: 1*time.Hour + 30*time.Minute + 15*time.Second + 500*time.Millisecond,
					EndTime:   1*time.Hour + 30*time.Minute + 20*time.Second + 750*time.Millisecond,
					Text:      "Complex timing subtitle",
				},
			},
			wantErr:     false,
			description: "should write subtitle with hours, minutes, seconds, and milliseconds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile, err := os.CreateTemp("", "test_output_*.srt")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			filename := tmpFile.Name()
			tmpFile.Close()
			defer os.Remove(filename)

			err = converter.WriteToFile(tt.subtitles, filename)

			if tt.wantErr {
				if err == nil {
					t.Errorf("WriteToFile() expected error for %s, got nil", tt.description)
				}
				return
			}

			if err != nil {
				t.Errorf("WriteToFile() unexpected error for %s: %v", tt.description, err)
				return
			}

			// Verify file was created and has content
			content, err := os.ReadFile(filename)
			if err != nil {
				t.Errorf("Failed to read output file for %s: %v", tt.description, err)
				return
			}

			if len(tt.subtitles) > 0 && len(content) == 0 {
				t.Errorf("WriteToFile() %s: output file is empty but subtitles were provided", tt.description)
			}

			// Verify SRT format structure for non-empty subtitle lists
			if len(tt.subtitles) > 0 {
				contentStr := string(content)
				if !strings.Contains(contentStr, "1\n") {
					t.Errorf("WriteToFile() %s: output missing sequence number", tt.description)
				}
				if !strings.Contains(contentStr, " --> ") {
					t.Errorf("WriteToFile() %s: output missing SRT timestamp separator", tt.description)
				}
			}
		})
	}
}
