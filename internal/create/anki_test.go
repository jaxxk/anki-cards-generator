package create

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestEnsureAnkiConnect(t *testing.T) {
	ok, err := EnsureAnkiConnect()
	assert.NoError(t, err, "expected no error, but got one")
	if !ok {
		t.Fatal("expected true but got false")
	}
}

func TestCreateDeck(t *testing.T) {
	ok, err := CreateDeck("test", zap.NewExample().Sugar())
	assert.NoError(t, err, "expected no error, but got one")
	if !ok {
		t.Fatal("true but got false")
	}
}

func TestGetDeck(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedError  bool
		expectedOutput bool
	}{
		{
			name:           "invalid deck",
			input:          "dne", // Deck does not exist
			expectedError:  false, // GetDeck should handle this gracefully without returning an error
			expectedOutput: false, // Since the deck does not exist
		},
		{
			name:           "valid deck",
			input:          "test", // Deck exists
			expectedError:  false,  // GetDeck should return no error
			expectedOutput: true,   // Deck exists
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := GetDeck(tt.input, zap.NewExample().Sugar())

			// Check for errors
			if (err != nil) != tt.expectedError {
				t.Errorf("GetDeck(%q) error = %v, expectedError = %v", tt.input, err, tt.expectedError)
			}

			// Check for output
			if output != tt.expectedOutput {
				t.Errorf("GetDeck(%q) = %v, expectedOutput = %v", tt.input, output, tt.expectedOutput)
			}
		})
	}
}
