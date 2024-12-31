package transform

import (
	"os"
	"testing"

	"github.com/jaxxk/anki-cards-generator/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSaveDeck(t *testing.T) {
	mockDeck := Deck{Title: "Test Deck", Cards: []Flashcards{{Front: "Q1", Back: "A1"}}}

	// Call SaveDeck
	jsonPath, err := SaveDeck(mockDeck)
	assert.NoError(t, err, "SaveDeck should not return an error")

	if _, err := os.Stat(jsonPath); err != nil {
		assert.NoError(t, err, "json file does not exist in path: %v", jsonPath)
	}

	processingPath, err := utils.CreateProcessingDir()
	if err != nil {
		assert.NoError(t, err, "could not get the processing dir path")
	}
	// Clean up
	os.RemoveAll(processingPath)
}
