/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/jaxxk/anki-cards-generator/internal/encryption"
	"github.com/jaxxk/anki-cards-generator/pkg/logging"
	"github.com/spf13/cobra"
)

// generateEncryptionCmd represents the generateEncryption command
var generateEncryptionCmd = &cobra.Command{
	Use:   "generateEncryption",
	Short: "Generates an encryption key that will be stored to your env",
	Long: `To generate an encryption key, run:
poggers addKey generateEncryption

make sure to source the shell first`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := logging.FromContext(ctx)
		err := encryption.CreateEncryptionKey()
		if err != nil {
			logger.Errorf("Failed to generate encryption key: %v", err)
			return err
		}
		return nil
	},
}

func init() {

}
