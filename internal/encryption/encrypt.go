package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"runtime"
)

var ENC_KEY string = "encryption_key"

type Key struct {
	Key string `json:"key"`
}

// EncryptFile encrypts the file at inputPath and writes the result to outputPath.
func EncryptFile(inputPath, outputPath string, key []byte) error {
	// Open the input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Read the input file contents
	plainText, err := io.ReadAll(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Create a new AES cipher using the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Generate a random nonce
	nonce := make([]byte, 12) // GCM nonce size is 12 bytes
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Create a GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	// Encrypt the plaintext
	cipherText := aesGCM.Seal(nil, nonce, plainText, nil)

	// Open the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Write the nonce followed by the ciphertext to the output file
	if _, err := outputFile.Write(nonce); err != nil {
		return fmt.Errorf("failed to write nonce to output file: %w", err)
	}
	if _, err := outputFile.Write(cipherText); err != nil {
		return fmt.Errorf("failed to write ciphertext to output file: %w", err)
	}

	return nil
}

func DecryptFile(inputPath string, key []byte) (*Key, error) {
	// Open the input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Read the file contents
	data, err := io.ReadAll(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	// Separate nonce and ciphertext
	nonce, cipherText := data[:12], data[12:]

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt the ciphertext
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	// Parse the decrypted plaintext into the Key struct
	var keyStruct Key
	if err := json.Unmarshal(plainText, &keyStruct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal decrypted data into struct: %w", err)
	}

	return &keyStruct, nil
}

// GenerateRandomKey generates a random 32-byte key for AES-256.
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256 requires 32 bytes
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// SaveEncryptionKeyToEnv saves a []byte to an environment variable and persists it for future sessions.
func SaveEncryptionKeyToEnv(key string, data []byte) error {
	// Convert []byte to base64 string
	encoded := base64.StdEncoding.EncodeToString(data)

	// Save the environment variable for the current session
	if err := os.Setenv(key, encoded); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}

	// Persist the environment variable for future sessions
	switch runtime.GOOS {
	case "windows":
		// Use `setx` to persist environment variable on Windows
		cmd := exec.Command("setx", key, encoded)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to persist environment variable on Windows: %w", err)
		}
	case "linux", "darwin":
		// Determine the shell configuration file
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}

		// Check for common shell configuration files
		shellConfigPath := usr.HomeDir + "/.bashrc"
		if shell := os.Getenv("SHELL"); shell != "" && shell == "/bin/zsh" {
			shellConfigPath = usr.HomeDir + "/.zshrc"
		}

		// Append the export statement to the shell configuration file
		exportCmd := fmt.Sprintf("export %s=%s\n", key, encoded)
		file, err := os.OpenFile(shellConfigPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("failed to open shell configuration file: %w", err)
		}
		defer file.Close()

		if _, err := file.WriteString(exportCmd); err != nil {
			return fmt.Errorf("failed to write to shell configuration file: %w", err)
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return nil
}

// GetBytesFromEnv retrieves a []byte from an environment variable
func GetBytesFromEnv(key string) ([]byte, error) {
	// Retrieve the base64 string from the environment variable
	encoded := os.Getenv(key)
	if encoded == "" {
		return nil, fmt.Errorf("environment variable %s not found. Run poggers addkey generate-encryption", key)
	}

	// Decode the base64 string back into []byte
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode environment variable: %w", err)
	}
	return data, nil
}

func GetEncKey() ([]byte, error) {
	encKey, err := GetBytesFromEnv(ENC_KEY)
	if err != nil {
		return nil, err
	}
	return encKey, nil
}
