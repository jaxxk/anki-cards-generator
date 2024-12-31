package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestReadFromFileToString(t *testing.T) {
	// Step 1: Create a temporary file
	tempFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Ensure the file is deleted after the test

	expected := "hello\ngo\n"
	if _, err := tempFile.WriteString(expected); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	logger := zaptest.NewLogger(t).Sugar()

	result, err := ReadFromFileToString(tempFile.Name(), logger)
	if err != nil {
		t.Fatalf("ReadFromFileToString returned an error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}

func TestReadFromFileToString_FileNotFound(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	// Attempt to read a nonexistent file
	_, err := ReadFromFileToString("nonexistentfile.txt", logger)
	if err == nil {
		t.Errorf("Expected an error for nonexistent file but got none")
	}
}

// Updated Test for CreateProcessingDir
func TestCreateProcessingDir(t *testing.T) {
	t.Run("Create new processing directory", func(t *testing.T) {
		// Get the path where the directory should be created
		userHome, err := os.UserHomeDir()
		assert.NoError(t, err)

		processingDirPath := userHome + string(os.PathSeparator) + PROCESSING_DIR

		// Ensure the directory does not exist before the test
		if _, err := os.Stat(processingDirPath); err == nil {
			os.RemoveAll(processingDirPath)
		}

		// Call CreateProcessingDir
		path, err := CreateProcessingDir()
		assert.NoError(t, err)

		// Verify the returned path matches the expected path
		assert.Equal(t, processingDirPath, path, "Returned path should match expected processing directory path")

		// Check if the directory now exists
		exists, err := directoryExists(path)
		assert.NoError(t, err)
		assert.True(t, exists, "Processing directory should exist after creation")

		// Clean up
		os.RemoveAll(path)
	})

	t.Run("Directory already exists", func(t *testing.T) {
		// Get the path where the directory should exist
		userHome, err := os.UserHomeDir()
		assert.NoError(t, err)

		processingDirPath := userHome + string(os.PathSeparator) + PROCESSING_DIR

		// Create the directory manually
		err = os.Mkdir(processingDirPath, 0755)
		assert.NoError(t, err)

		// Call CreateProcessingDir and ensure no error occurs
		path, err := CreateProcessingDir()
		assert.NoError(t, err)

		// Verify the returned path matches the expected path
		assert.Equal(t, processingDirPath, path, "Returned path should match expected processing directory path")

		// Check the directory still exists
		exists, err := directoryExists(path)
		assert.NoError(t, err)
		assert.True(t, exists, "Processing directory should still exist")

		// Clean up
		os.RemoveAll(path)
	})
}
