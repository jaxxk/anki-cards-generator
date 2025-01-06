package create

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

var ANKI_ENDPOINT = "http://localhost:8765"

type AnkiConnect_Response struct {
	Result []string `json:"result"`
	Error  string   `json:"error"`
}

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

func CreateDeck(deckName string, logger *zap.SugaredLogger) (bool, error) {
	params := map[string]string{
		"deck": deckName,
	}
	// Prepare the request body
	reqBody := map[string]interface{}{
		"action":  "createDeck",
		"version": 6,
		"params":  params,
	}

	response, err := processRequest(reqBody)
	if err != nil {
		logger.Errorf("Failed to process request: %v", err)
		return false, err
	}
	// Non-empty response
	if response.Error != "null" && response.Result != nil {
		return true, nil
	}
	return false, errors.New("failed to create deck, empty reponse from anki connect")
}

// GetDeck checks if a given deck exists in Anki
func GetDeck(deckName string, logger *zap.SugaredLogger) (bool, error) {
	// Prepare the request body
	reqBody := map[string]interface{}{
		"action":  "deckNames",
		"version": 6,
	}
	response, err := processRequest(reqBody)
	if err != nil {
		logger.Errorf("Failed to process request: %v", err)
		return false, err
	}
	// Check if the deck exists in the result
	for _, name := range response.Result {
		if name == deckName {
			return true, nil
		}
	}
	return false, nil
}

func processRequest(reqBody interface{}) (AnkiConnect_Response, error) {
	response := AnkiConnect_Response{}
	// Serialize the request body to JSON
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return response, fmt.Errorf("failed to serialize request body: %w", err)
	}

	// Make the HTTP POST request to AnkiConnect
	resp, err := http.Post(ANKI_ENDPOINT, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return response, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response body
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors in the response
	if response.Error != "" {
		return response, fmt.Errorf("anki error: %s", response.Error)
	}
	return response, nil
}
