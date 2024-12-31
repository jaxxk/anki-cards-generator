package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

var PROCESSING_DIR string = ".anki-cards-generator"

func ResolvePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %v", err)
	}
	return absPath, nil
}

func ReadFromFileToString(file string, logger *zap.SugaredLogger) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		if logger != nil {
			logger.Errorf("failed to read file: %v", err)
		}
		return "", err
	}
	return string(data), nil
}

// WriteJSONToFile writes a given data object to a JSON file in the specified directory.
// Returns the full path to the created JSON file.
func WriteJSONToFile(data interface{}, folder, filename string) (string, error) {
	filePath := filepath.Join(folder, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %s, error: %v", filePath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return "", fmt.Errorf("failed to write JSON to file: %s, error: %v", filePath, err)
	}

	return filePath, nil
}

// ValidateAndResolvePath validates the provided file path and resolves it to an absolute path.
func ValidateAndResolvePath(path string, logger *zap.SugaredLogger) (string, error) {
	if path == "" {
		logger.Error("File path is empty")
		return "", errors.New("file path cannot be empty")
	}

	resolvedPath, err := ResolvePath(path)
	if err != nil {
		logger.Errorf("Failed to resolve file path: %v", err)
		return "", fmt.Errorf("failed to resolve file path: %v", err)
	}

	return resolvedPath, nil
}

// ReadFromFile reads the content of a file and returns it as a string.
func ReadFromFile(path string, logger *zap.SugaredLogger) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Errorf("File does not exist: %s", path)
		return "", fmt.Errorf("file does not exist: %s", path)
	}

	logger.Infof("Reading file: %s", path)
	inputText, err := ReadFromFileToString(path, logger)
	if err != nil {
		logger.Errorf("Failed to read from file: %s, error: %v", path, err)
		return "", fmt.Errorf("failed to read from file: %s, error: %v", path, err)
	}

	return inputText, nil
}

func directoryExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // The path does not exist
		}
		return false, err // Some other error occurred
	}
	return true, nil // Check if it's a directory
}

// CreateProcessingDir creates a directory under user Home called PROCESSING_DIR if it doesn't exist already.
// Returns path to the processing directory
func CreateProcessingDir() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	processingDirPath := userHome + string(os.PathSeparator) + PROCESSING_DIR
	dirExists, err := directoryExists(processingDirPath)
	if err != nil {
		return "", err
	}
	if !dirExists {
		os.Mkdir(processingDirPath, 0755)
	}
	return processingDirPath, nil
}

// Generates a random file name with a prefix and extension
func GenerateRandomFileName(prefix, extension string) (string, error) {
	// Generate 16 random bytes
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert to hex string and construct file name
	randomHex := hex.EncodeToString(randomBytes)
	fileName := fmt.Sprintf("%s-%s%s", prefix, randomHex, extension)

	return fileName, nil
}
