package encryption

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaxxk/anki-cards-generator/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSaveBytesToEnv(t *testing.T) {
	key := "TEST_ENV_VAR"
	data := []byte("test-data")
	expectedValue := base64.StdEncoding.EncodeToString(data)

	// Run the function
	err := SaveEncryptionKeyToEnv(key, data)
	if err != nil {
		t.Fatalf("SaveBytesToEnv failed: %v", err)
	}

	// Verify the environment variable is set for the current process
	actualValue := os.Getenv(key)
	if actualValue != expectedValue {
		t.Errorf("Environment variable value mismatch: got %v, want %v", actualValue, expectedValue)
	}

	// Unset the environment variable to verify cleanup
	if err := os.Unsetenv(key); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}

	// Verify that the environment variable is unset
	if _, exists := os.LookupEnv(key); exists {
		t.Errorf("Environment variable %s was not properly unset", key)
	}
}

// Test function for SaveAPIKey with filesystem
func TestSaveAPIKey_FileSystem(t *testing.T) {
	key := "test-api-key"
	err := CreateEncryptionKey()
	assert.NoError(t, err)
	// Call SaveAPIKey
	err = SaveAPIKey(key, zap.NewExample().Sugar())
	assert.NoError(t, err)
	processing_dir, err := utils.CreateProcessingDir()
	assert.NoError(t, err)

	// Verify that the encrypted file exists
	encryptedFile := filepath.Join(processing_dir, ENC_KEY_FILE)
	_, err = os.Stat(encryptedFile)
	assert.NoError(t, err, "encrypted file should exist")
	assert.FileExists(t, encryptedFile)
	// clean up
	os.Remove(encryptedFile)
	// Unset the environment variable to verify cleanup
	if err := os.Unsetenv(ENC_KEY); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}

	// Verify that the environment variable is unset
	if _, exists := os.LookupEnv(ENC_KEY); exists {
		t.Errorf("Environment variable %s was not properly unset", key)
	}
}

func TestGetAPIKey(t *testing.T) {
	key := "test-api-key"
	err := CreateEncryptionKey()
	assert.NoError(t, err)
	// Call SaveAPIKey
	err = SaveAPIKey(key, zap.NewExample().Sugar())
	assert.NoError(t, err)
	// Call GetAPIKey
	actualKey, err := GetAPIKey(zap.NewExample().Sugar())
	assert.NoError(t, err)
	// validate
	assert.Equal(t, actualKey, key)

	processing_dir, err := utils.CreateProcessingDir()
	assert.NoError(t, err)
	// Verify that the encrypted file exists
	encryptedFile := filepath.Join(processing_dir, ENC_KEY_FILE)
	_, err = os.Stat(encryptedFile)
	assert.NoError(t, err, "encrypted file should exist")
	assert.FileExists(t, encryptedFile)
	// clean up
	os.Remove(encryptedFile)
	// Unset the environment variable to verify cleanup
	if err := os.Unsetenv(ENC_KEY); err != nil {
		t.Errorf("Failed to unset environment variable: %v", err)
	}

	// Verify that the environment variable is unset
	if _, exists := os.LookupEnv(ENC_KEY); exists {
		t.Errorf("Environment variable %s was not properly unset", key)
	}
}
