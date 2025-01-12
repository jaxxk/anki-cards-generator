/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addKeyCmd represents the addKey command
var addKeyCmd = &cobra.Command{
	Use:   "addKey",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addKey called")
	},
}

func init() {
	rootCmd.AddCommand(addKeyCmd)

	addKeyCmd.AddCommand(generateEncryptionCmd)
}
