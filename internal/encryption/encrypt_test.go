package encryption

import (
	"encoding/base64"
	"os"
	"testing"
)

func TestSaveBytesToEnv(t *testing.T) {
	key := "TEST_ENV_VAR"
	data := []byte("test-data")
	expectedValue := base64.StdEncoding.EncodeToString(data)

	// Save the original state of the environment variable
	originalValue, exists := os.LookupEnv(key)

	// Ensure cleanup after the test using defer
	defer func() {
		if exists {
			// Restore the original value if it was set
			_ = os.Setenv(key, originalValue)
		} else {
			// Unset the environment variable if it didn't exist
			_ = os.Unsetenv(key)
		}
	}()

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
