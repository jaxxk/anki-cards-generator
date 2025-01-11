package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jaxxk/anki-cards-generator/internal/create"
	"github.com/jaxxk/anki-cards-generator/internal/transform"
	"github.com/jaxxk/anki-cards-generator/pkg/logging"
	"github.com/jaxxk/anki-cards-generator/pkg/utils"
	"github.com/spf13/cobra"
)

var FilePath string
var Title string

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

		// ensures anki is running before processing the notes
		if ok, err := create.EnsureAnkiConnect(); err != nil || !ok {
			return errors.New("cannot connect to Anki Connect")
		}

		// Validate and resolve file path
		FilePath, err := utils.ValidateAndResolvePath(FilePath, logger)
		if err != nil {
			return fmt.Errorf("validation error for file path: %w", err)
		}

		// Read input text from file
		inputText, err := utils.ReadFromFile(FilePath, logger)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", FilePath, err)
		}

		// Generate chat completion using the transform package
		result := transform.NewChatCompletion(ctx, inputText)
		if result == nil || len(result.Choices) == 0 {
			logger.Error("Failed to generate flashcards or received empty response")
			return fmt.Errorf("failed to generate flashcards or received empty response")
		}

		// Unmarshal JSON response into flashcards deck
		rawOutput := result.Choices[0].Message.Content
		newDeck := transform.Deck{}
		err = json.Unmarshal([]byte(rawOutput), &newDeck)
		if err != nil {
			logger.Errorf("Failed to parse flashcards JSON: %v", err)
			return fmt.Errorf("invalid JSON response from transform package")
		}

		// Validate deck content
		if len(newDeck.Cards) == 0 {
			logger.Error("Generated flashcards deck is empty")
			return fmt.Errorf("generated flashcards deck is empty")
		}

		// Update Title
		if len(Title) > 0 {
			newDeck.UpdateTitle(Title)
		}

		// Save Deck to Processing Dir For retry
		jsonPath, err := transform.SaveDeck(newDeck)
		if err != nil {
			logger.Error("Failed to save deck to %v", jsonPath)
			return fmt.Errorf("failed to save deck to %v", jsonPath)
		}

		logger.Infof("Successfully Created %v deck JSON", newDeck.Title)

		err = create.SendToAnki(newDeck, logger)
		if err != nil {
			return err
		}
		logger.Infof("Successfully Created %v deck in anki", newDeck.Title)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Add file flag
	generateCmd.Flags().StringVarP(&FilePath, "file", "f", "", "Path to .md/.txt file (required)")
	generateCmd.MarkFlagRequired("file")
	// Add title flag
	generateCmd.Flags().StringVarP(&Title, "title", "t", "", "Title for the generated deck of flashcards (optional) will be automatically generated")
}
