package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var ANKI_ENDPOINT = "http://localhost:8765"

// EnsureAnkiConnect checks if AnkiConnect is running
func EnsureAnkiConnect() (bool, error) {
	resp, err := http.Get(ANKI_ENDPOINT)
	if err != nil {
		return false, fmt.Errorf("failed to connect to AnkiConnect: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return true, nil
}

// GetDeck checks if a given deck exists in Anki
func GetDeck(deckName string) (bool, error) {
	// Prepare the request body
	reqBody := map[string]interface{}{
		"action":  "deckNames",
		"version": 6,
	}

	// Serialize the request body to JSON
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return false, fmt.Errorf("failed to serialize request body: %w", err)
	}

	// Make the HTTP POST request to AnkiConnect
	resp, err := http.Post(ANKI_ENDPOINT, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return false, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response body
	var response struct {
		Result []string `json:"result"`
		Error  string   `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors in the response
	if response.Error != "" {
		return false, fmt.Errorf("anki error: %s", response.Error)
	}

	// Check if the deck exists in the result
	for _, name := range response.Result {
		if name == deckName {
			return true, nil
		}
	}
	return false, nil
}

// func CreateDeck()(bool, error) {

// }
