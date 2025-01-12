package encryption

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jaxxk/anki-cards-generator/pkg/utils"
	"go.uber.org/zap"
)

var KEY_FILE = "KEY_FILE.json"
var ENC_KEY_FILE = "key.enc"

func CreateEncryptionKey() error {
	key, err := GenerateRandomKey()
	if err != nil {
		return err
	}
	err = SaveEncryptionKeyToEnv(ENC_KEY, key)
	if err != nil {
		return err
	}
	return nil
}

func SaveAPIKey(key string, logger *zap.SugaredLogger) error {
	// Create a new Key struct
	newKey := Key{
		Key: key,
	}

	// Create a processing directory
	processingDirPath, err := utils.CreateProcessingDir()
	if err != nil {
		return fmt.Errorf("failed to create processing directory: %w", err)
	}

	// Write the key to a JSON file
	keyFilePath, err := utils.WriteJSONToFile(newKey, processingDirPath, KEY_FILE)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	encryptionKey, err := GetEncKey()
	if err != nil {
		return err
	}
	outputEncryptionFile := processingDirPath + string(os.PathSeparator) + ENC_KEY_FILE
	// Encrypt the key file
	err = EncryptFile(keyFilePath, outputEncryptionFile, encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt key file: %w", err)
	}

	// Optionally remove the plaintext file for security
	err = os.Remove(keyFilePath)
	if err != nil {
		return fmt.Errorf("failed to remove plaintext key file: %w", err)
	}

	logger.Infof("Encrypted key file successfully saved in: %s\n", filepath.Join(processingDirPath, KEY_FILE))
	return nil
}

func GetAPIKey(logger *zap.SugaredLogger) (string, error) {
	processingDirPath, err := utils.CreateProcessingDir()
	if err != nil {
		return "", err
	}
	encryptedFilePath := processingDirPath + string(os.PathSeparator) + ENC_KEY_FILE
	if _, err := os.Stat(encryptedFilePath); os.IsNotExist(err) {
		logger.Errorf("File does not exist: %s", encryptedFilePath)
		return "", fmt.Errorf("file does not exist: %s", encryptedFilePath)
	}
	encryptionKey, err := GetEncKey()
	if err != nil {
		return "", err
	}
	apiKey, err := DecryptFile(encryptedFilePath, encryptionKey)
	if err != nil {
		logger.Errorf("Failed to decrypt api key at :%s", encryptedFilePath)
		return "", fmt.Errorf("failed to decrypt api key at :%s", encryptedFilePath)
	}

	return apiKey.Key, nil
}
