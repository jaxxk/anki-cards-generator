/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/jaxxk/anki-cards-generator/internal/encryption"
	"github.com/jaxxk/anki-cards-generator/pkg/logging"
	"github.com/spf13/cobra"
)

var Key string

// addKeyCmd represents the addKey command
var addKeyCmd = &cobra.Command{
	Use:   "addKey",
	Short: "Adds your OpenAI API key to the poggers",
	Long: `Need to run poggers addKey generateEncryption to generate the encryption before 
	running poggers addKey -k <your-key>
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := logging.FromContext(ctx)
		_, err := encryption.GetEncKey()
		if err != nil {
			logger.Error("Need to generate encryption key first, run poggers addKey -h")
			return fmt.Errorf("need to generate encryption key first, run poggers addKey -h")
		}
		err = encryption.SaveAPIKey(Key, logger)
		if err != nil {
			logger.Errorf("Failed to save and encrypt api key: %v", err)
			return fmt.Errorf("failed to save and encrypt api key: %v", err)
		}
		return nil
	},
}

func init() {
	addKeyCmd.AddCommand(generateEncryptionCmd)
	// Add key flag
	addKeyCmd.Flags().StringVarP(&Key, "key", "k", "", "OpenAI API Key (required)")
	addKeyCmd.MarkFlagRequired("key")
	rootCmd.AddCommand(addKeyCmd)

}
