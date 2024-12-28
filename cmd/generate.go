/*
Copyright © 2024 Jack Wei jackwei2018@outlook.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var FilePath string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates anki flashcards based off input md/txt file",
	Long: `generate expects the flag --file or -f followed by the full path of the
	file. Example: 
	
	poggers generate -f /Users/jaxk/notes/notes.md or poggers generate --file ~/notes/notes.md
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// ctx := cmd.Context()
		// logger := logging.FromContext(ctx)
		if FilePath == "" {
			return fmt.Errorf("file path cannot be empty")
		}

		if _, err := os.Stat(FilePath); os.IsNotExist(err) {
			return fmt.Errorf("error: file %s does not exist", FilePath)
		}

		fmt.Println(FilePath) // For valid cases, print the file path
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	generateCmd.Flags().StringVarP(&FilePath, "file", "f", "", "path to md/txt file (required)")
	generateCmd.MarkFlagRequired("file")
}
