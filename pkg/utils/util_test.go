package utils

import (
	"os"
	"testing"

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
