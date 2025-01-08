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

var ANKI_ENDPOINT = "http://localhost:8765"

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

	// Decode the generic response
	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(rawResp, &genericResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors in the response
	if genericResp.Error != "" {
		return nil, fmt.Errorf("anki error: %s", genericResp.Error)
	}

	// Attempt to decode `Result` as both types
	var resultArray []string
	if err := json.Unmarshal(genericResp.Result, &resultArray); err == nil {
		return resultArray, nil // Successfully decoded as []string
	}

	var resultString string
	if err := json.Unmarshal(genericResp.Result, &resultString); err == nil {
		return resultString, nil // Successfully decoded as string
	}

	var resultNumber json.Number
	if err := json.Unmarshal(genericResp.Result, &resultNumber); err == nil {
		return resultNumber, nil // Successfully decoded as a number
	}

	// Log the raw result for debugging if no type matches
	return nil, fmt.Errorf("unknown result format: %v", genericResp.Result)
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
	reqBody := map[string]interface{}{
		"action":  "createDeck",
		"version": 6,
		"params": map[string]string{
			"deck": deckName,
		},
	}

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
