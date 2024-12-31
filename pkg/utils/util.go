package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

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
