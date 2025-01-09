package create

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// var ANKI_ENDPOINT = "http://localhost:8765"
var ANKI_ENDPOINT = "http://host.docker.internal:8765"

type AnkiConnectGenericResponse struct {
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
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

func processRequest(reqBody interface{}) (interface{}, error) {
	var genericResp AnkiConnectGenericResponse

	// Serialize the request body to JSON
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request body: %w", err)
	}

	// Make the HTTP POST request to AnkiConnect
	resp, err := http.Post(ANKI_ENDPOINT, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the raw response
	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal into the generic response structure
	if err := json.Unmarshal(rawResp, &genericResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Handle errors from AnkiConnect
	if genericResp.Error != "" {
		var decodedError string
		if err := json.Unmarshal([]byte(genericResp.Error), &decodedError); err == nil {
			return nil, fmt.Errorf("anki error: %s", decodedError)
		}
		return nil, fmt.Errorf("anki error: %s", genericResp.Error)
	}

	// Decode `Result` dynamically based on its type
	var result interface{}
	if err := decodeResult(genericResp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to decode result: %w", err)
	}

	return result, nil
}

// Helper function to decode `Result` into its appropriate type
func decodeResult(rawMessage json.RawMessage, result *interface{}) error {
	// Attempt to decode as []string
	var resultArray []string
	if err := json.Unmarshal(rawMessage, &resultArray); err == nil {
		*result = resultArray
		return nil
	}

	// Attempt to decode as string
	var resultString string
	if err := json.Unmarshal(rawMessage, &resultString); err == nil {
		*result = resultString
		return nil
	}

	// Attempt to decode as json.Number
	var resultNumber json.Number
	if err := json.Unmarshal(rawMessage, &resultNumber); err == nil {
		*result = resultNumber
		return nil
	}
	// Attempt to decode as []int64
	var resultArrayInt []int64
	if err := json.Unmarshal(rawMessage, &resultArrayInt); err == nil {
		*result = resultArrayInt
		return nil
	}
	// If no known type matches, return an error
	return fmt.Errorf("unknown result format: %s", string(rawMessage))
}

// GetDeck checks if a given deck exists in Anki
func GetDeck(deckName string, logger *zap.SugaredLogger) (bool, error) {
	reqBody := map[string]interface{}{
		"action":  "deckNames",
		"version": 6,
	}

	resp, err := processRequest(reqBody)
	if err != nil {
		logger.Errorf("Failed to process request: %v", err)
		return false, err
	}

	// Check if the result is a []string
	if deckNames, ok := resp.([]string); ok {
		for _, name := range deckNames {
			if name == deckName {
				return true, nil
			}
		}
		return false, nil
	}

	return false, fmt.Errorf("unexpected result format for deckNames")
}

// CreateDeck creates a new deck in Anki
func CreateDeck(deckName string, logger *zap.SugaredLogger) (bool, error) {
	reqBody := NewAnkiRequestBody("createDeck", map[string]string{
		"deck": deckName,
	})

	resp, err := processRequest(reqBody)
	if err != nil {
		logger.Errorf("Failed to process request: %v", err)
		return false, err
	}

	if resp == nil {
		return false, errors.New("reponse for createdeck is nil")
	}
	return true, nil
}

func deleteDeck(deckName string, logger *zap.SugaredLogger) (bool, error) {
	reqBody := NewAnkiRequestBody("deleteDecks", map[string]interface{}{
		"decks":    []string{deckName},
		"cardsToo": true,
	})
	_, err := processRequest(reqBody)
	if err != nil {
		logger.Errorf("Failed to process request: %v", err)
		return false, err
	}
	return true, nil
}
