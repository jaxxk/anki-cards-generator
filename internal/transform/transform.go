package transform

import (
	"fmt"

	"github.com/jaxxk/anki-cards-generator/pkg/utils"
)

// SaveDeck saves the deck object to a json file. Returns the path to the json file.
func SaveDeck(deck Deck) (string, error) {
	processingPath, err := utils.CreateProcessingDir()
	if err != nil {
		return "", err
	}

	randomFileName, err := utils.GenerateRandomFileName("deck", ".json")
	if err != nil {
		return "", err
	}

	jsonPath, err := utils.WriteJSONToFile(deck, processingPath, randomFileName)
	if err != nil {
		return "", fmt.Errorf("failed to write JSON to file: %w", err)
	}

	return jsonPath, nil
}
