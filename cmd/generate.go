package cmd

import (
	"fmt"
	"os"

	"github.com/jaxxk/anki-cards-generator/internal/transform"
	"github.com/jaxxk/anki-cards-generator/pkg/logging"
	"github.com/jaxxk/anki-cards-generator/pkg/utils"
	"github.com/spf13/cobra"
)

var FilePath string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates Anki flashcards from input .md/.txt file",
	Long: `The "generate" command processes the specified .md or .txt file to 
	generate insightful Anki flashcards based on its content and saves the result in a temporary JSON file.

	Example Usage:
	poggers generate -f /Users/jaxk/notes/notes.md
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := logging.FromContext(ctx)

		// Validate file path
		if FilePath == "" {
			logger.Error("File path is empty")
			return fmt.Errorf("file path cannot be empty")
		}

		// Resolve absolute file path
		resolvedPath, err := utils.ResolvePath(FilePath)
		if err != nil {
			logger.Errorf("Failed to resolve file path: %v", err)
			return err
		}
		FilePath = resolvedPath

		if _, err := os.Stat(FilePath); os.IsNotExist(err) {
			logger.Errorf("File does not exist: %s", FilePath)
			return fmt.Errorf("file does not exist: %s", FilePath)
		}

		logger.Infof("Reading file: %s", FilePath)
		// Read input text from the file
		inputText, err := utils.ReadFromFileToString(FilePath, logger)
		if err != nil {
			logger.Errorf("Failed to read from file: %s, error: %v", FilePath, err)
			return fmt.Errorf("failed to read from file: %s, error: %v", FilePath, err)
		}

		logger.Info("File read successfully. Starting flashcard generation.")
		// Generate chat completion using the transform package
		result := transform.NewChatCompletion(ctx, inputText)
		if result == nil || len(result.Choices) == 0 {
			logger.Error("Failed to generate flashcards or received empty response")
			return fmt.Errorf("failed to generate flashcards or received empty response")
		}

		// Write raw result to a file for debugging
		rawOutput := result.Choices[0].Text
		debugFile := "test.txt"
		if err := os.WriteFile(debugFile, []byte(rawOutput), 0644); err != nil {
			logger.Errorf("Failed to write raw output to file: %s, error: %v", debugFile, err)
			return fmt.Errorf("failed to write raw output to file: %s, error: %v", debugFile, err)
		}

		logger.Infof("Raw output successfully written to: %s", debugFile)
		fmt.Println(FilePath) // For valid cases, print the file path
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Add file flag
	generateCmd.Flags().StringVarP(&FilePath, "file", "f", "", "Path to .md/.txt file (required)")
	generateCmd.MarkFlagRequired("file")
}
