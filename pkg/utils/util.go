package utils

import (
	"encoding/json"
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
