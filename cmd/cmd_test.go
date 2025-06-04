package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateInputFile(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test*.sbv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Logf("Warning: failed to remove temp file: %v", err)
		}
	}()
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
			errMsg:  "input file path cannot be empty",
		},
		{
			name:    "non-existent file",
			input:   "nonexistent.sbv",
			wantErr: true,
			errMsg:  "input file does not exist",
		},
		{
			name:    "wrong extension",
			input:   tempFile.Name() + ".txt",
			wantErr: true,
			errMsg:  "input file must have .sbv extension",
		},
		{
			name:    "valid sbv file",
			input:   tempFile.Name(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For the wrong extension test, create a file with .txt extension
			if tt.name == "wrong extension" {
				txtFile, err := os.CreateTemp("", "test*.txt")
				if err != nil {
					t.Fatalf("Failed to create temp txt file: %v", err)
				}
				defer func() {
					if err := os.Remove(txtFile.Name()); err != nil {
						t.Logf("Warning: failed to remove temp txt file: %v", err)
					}
				}()
				if err := txtFile.Close(); err != nil {
					t.Fatalf("Failed to close temp txt file: %v", err)
				}
				tt.input = txtFile.Name()
			}

			err := validateInputFile(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInputFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("validateInputFile() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

func TestDetermineOutputPath(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test_output")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: failed to remove temp directory: %v", err)
		}
	}()

	tests := []struct {
		name    string
		input   string
		output  string
		want    string
		wantErr bool
		errMsg  string
	}{
		{
			name:   "auto-generate output from input",
			input:  "video.sbv",
			output: "",
			want:   "video.srt",
		},
		{
			name:   "auto-generate with path",
			input:  "/path/to/video.sbv",
			output: "",
			want:   "/path/to/video.srt",
		},
		{
			name:   "explicit output file",
			input:  "video.sbv",
			output: "subtitle.srt",
			want:   "subtitle.srt",
		},
		{
			name:   "explicit output with path",
			input:  "video.sbv",
			output: filepath.Join(tempDir, "output.srt"),
			want:   filepath.Join(tempDir, "output.srt"),
		},
		{
			name:    "output without .srt extension",
			input:   "video.sbv",
			output:  "output.txt",
			wantErr: true,
			errMsg:  "output file must have .srt extension",
		},
		{
			name:    "output directory doesn't exist",
			input:   "video.sbv",
			output:  "/nonexistent/dir/output.srt",
			wantErr: true,
			errMsg:  "output directory does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := determineOutputPath(tt.input, tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("determineOutputPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("determineOutputPath() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("determineOutputPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
