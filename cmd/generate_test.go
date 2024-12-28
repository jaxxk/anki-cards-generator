package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCmd(t *testing.T) {
	// Path to testdata directory
	testdataDir := "testdata"
	wd, _ := os.Getwd()
	sampleDataPath := filepath.Join(wd, testdataDir)
	// Define test cases
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "Valid File Path",
			args:           []string{"--file", filepath.Join(sampleDataPath, "sample1.md")},
			expectedOutput: filepath.Join(sampleDataPath, "sample1.md"),
			expectError:    false,
		},
		{
			name:           "Short Flag Valid File Path",
			args:           []string{"-f", filepath.Join(sampleDataPath, "sample2.md")},
			expectedOutput: filepath.Join(sampleDataPath, "sample2.md"),
			expectError:    false,
		},
		{
			name:           "Empty File Path",
			args:           []string{"--file", ""},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name:           "Nonexistent File",
			args:           []string{"--file", filepath.Join(sampleDataPath, "nonexistent.md")},
			expectedOutput: "",
			expectError:    true,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := new(bytes.Buffer)

			// Initialize the command
			cmd := &cobra.Command{
				Use:   generateCmd.Use,
				Short: generateCmd.Short,
				RunE:  generateCmd.RunE,
			}
			cmd.Flags().StringVarP(&FilePath, "file", "f", "", "path to md/txt file (required)")
			cmd.MarkFlagRequired("file")
			cmd.SetArgs(tt.args)
			cmd.SetOut(output)
			cmd.SetErr(output)

			// Execute the command
			err := cmd.Execute()
			if tt.expectError {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "expected no error")
			}
		})
	}
}
