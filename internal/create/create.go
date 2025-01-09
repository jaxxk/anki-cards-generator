package create

import (
	"fmt"

	"github.com/jaxxk/anki-cards-generator/internal/transform"
	"go.uber.org/zap"
)

func SendToAnki(deck transform.Deck, logger *zap.SugaredLogger) error {
	// Ensure the deck exists
	if err := existsDeck(deck.Title, logger); err != nil {
		return fmt.Errorf("failed to ensure deck exists: %w", err)
	}

	// Prepare to batch cards into `FlashcardBatchSize`
	batch := []Note{}
	for i, card := range deck.Cards {
		// Create a new Note from the card
		note := NewNote(card.Front, card.Back, deck.Title)
		batch = append(batch, note)

		// Send the batch if it reaches the `FlashcardBatchSize`
		if len(batch) == FlashcardBatchSize || i == len(deck.Cards)-1 {
			if err := sendBatchToAnki(batch, logger); err != nil {
				return fmt.Errorf("failed to send batch to Anki: %w", err)
			}
			batch = []Note{} // Reset the batch after sending
		}
	}

	return nil
}

func sendBatchToAnki(batch []Note, logger *zap.SugaredLogger) error {
	// Prepare the request body
	params := Notes{ListOfNotes: batch}
	reqBody := NewAnkiRequestBody("addNotes", params)

	// Send the request
	_, err := processRequest(reqBody)
	if err != nil {
		logger.Errorf("Failed to send batch to Anki: %v", err)
		return err
	}

	logger.Infof("Successfully sent %d cards to Anki", len(batch))
	return nil
}

func existsDeck(title string, logger *zap.SugaredLogger) error {
	// Check if the deck exists
	deckExists, err := GetDeck(title, logger)
	if err != nil {
		return fmt.Errorf("failed to check if deck exists: %w", err)
	}

	// If the deck does not exist, create it
	if !deckExists {
		if _, err := CreateDeck(title, logger); err != nil {
			return fmt.Errorf("failed to create deck '%s': %w", title, err)
		}
	}
	return nil
}
