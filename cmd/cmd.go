package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/un-versed/go-sbv-to-srt/pkg/sbv"
)

var (
	inputFile  string
	outputFile string
	version    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-sbv-to-srt",
	Short: "Convert SBV subtitle files to SRT format",
	Long: `A CLI tool to convert SBV (SubViewer) subtitle files to SRT (SubRip) format.
		SBV files are commonly used by YouTube and other platforms, while SRT is a more
		widely supported subtitle format that can be used across various media players
		and video editing software.

		Examples:
		go-sbv-to-srt -i input.sbv
		go-sbv-to-srt -i input.sbv -o output.srt
		go-sbv-to-srt --input video.sbv --output subtitles.srt`,
	RunE: convertSbvToSrt,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo sets the version information
func SetVersionInfo(v string) {
	version = v
}

// init initializes the root command and its flags
// It also sets up the version command as a subcommand.
func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input SBV file path (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output SRT file path (optional - defaults to input filename with .srt extension)")
	if err := rootCmd.MarkFlagRequired("input"); err != nil {
		panic(fmt.Sprintf("Failed to mark flag as required: %v", err))
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("go-sbv-to-srt version %s\n", version)
		},
	}
	rootCmd.AddCommand(versionCmd)
}

func convertSbvToSrt(cmd *cobra.Command, args []string) error {
	if err := validateInputFile(inputFile); err != nil {
		return fmt.Errorf("input validation failed: %w", err)
	}

	// Determine output file path
	outputPath, err := determineOutputPath(inputFile, outputFile)
	if err != nil {
		return fmt.Errorf("output path determination failed: %w", err)
	}

	fmt.Printf("Converting SBV file: %s\n", inputFile)
	fmt.Printf("Output SRT file: %s\n", outputPath)

	// Create converter instance
	converter := sbv.NewConverter()

	// Parse the SBV file
	subtitles, err := converter.ParseFromFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse SBV file: %w", err)
	}

	fmt.Printf("Parsed %d subtitle entries\n", len(subtitles))

	// Convert and write to SRT file
	err = converter.WriteToFile(subtitles, outputPath)
	if err != nil {
		return fmt.Errorf("failed to write SRT file: %w", err)
	}

	fmt.Printf("Successfully converted %d subtitles to SRT format\n", len(subtitles))
	fmt.Printf("Output saved to: %s\n", outputPath)

	return nil
}

func validateInputFile(input string) error {
	if input == "" {
		return fmt.Errorf("input file path cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", input)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(input))
	if ext != ".sbv" {
		return fmt.Errorf("input file must have .sbv extension, got: %s", ext)
	}

	file, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("cannot read input file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error but don't override the main error
			fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", closeErr)
		}
	}()

	return nil
}

func determineOutputPath(input, output string) (string, error) {
	if output != "" {
		outputDir := filepath.Dir(output)
		if outputDir != "." {
			// Check if output directory exists
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				return "", fmt.Errorf("output directory does not exist: %s", outputDir)
			}
		}

		// Ensure output has .srt extension
		if !strings.HasSuffix(strings.ToLower(output), ".srt") {
			return "", fmt.Errorf("output file must have .srt extension")
		}

		return output, nil
	}

	// Generate output filename from input
	inputBase := strings.TrimSuffix(input, filepath.Ext(input))
	outputPath := inputBase + ".srt"

	return outputPath, nil
}
