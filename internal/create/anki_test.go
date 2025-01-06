package create

import (
	"testing"
)

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
		// {
		// 	name:           "valid deck",
		// 	input:          "existingDeck", // Deck exists
		// 	expectedError:  false,          // GetDeck should return no error
		// 	expectedOutput: true,           // Deck exists
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := GetDeck(tt.input)

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
