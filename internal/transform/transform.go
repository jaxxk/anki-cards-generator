package transform

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jaxxk/anki-cards-generator/pkg/logging"
	"github.com/jaxxk/anki-cards-generator/pkg/utils"
)

const FILE_SIZE_LIMIT int64 = 500000000

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

// Returns a channel of Deck, will block until all of the file has been scanned. Handles the closing of the two channels.
func streamDocument(ctx context.Context, docPath string) (<-chan Deck, <-chan error) {
	decksCh := make(chan Deck)
	errCh := make(chan error, 1) // buffer of 1 so send won't block if no one reads immediately

	go func() {
		// Ensure we close channels when we're done
		defer close(decksCh)
		defer close(errCh)

		// Check file size
		fileInfo, err := os.Stat(docPath)
		if err != nil {
			errCh <- fmt.Errorf("failed to stat file: %w", err)
			return
		}
		if fileInfo.Size() > FILE_SIZE_LIMIT {
			errCh <- fmt.Errorf("file too large to process (%d bytes), limit %d",
				fileInfo.Size(), FILE_SIZE_LIMIT)
			return
		}

		// Open file
		file, err := os.Open(docPath)
		if err != nil {
			errCh <- fmt.Errorf("failed to open file: %w", err)
			return
		}
		defer file.Close()

		// Set up scanner
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)

		maxWords := 800
		minWords := 500
		var words []string

		//  Read and accumulate words
		for scanner.Scan() {
			// Respect context cancellation
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
				// proceed
			}

			word := scanner.Text()

			// If we’ve reached the triple dash (after minWords) or hit maxWords, we create a new deck.
			if (len(words) > minWords && word == "---") || len(words) >= maxWords {
				accumulatedText := strings.Join(words, " ")
				deck, err := createDeck(ctx, accumulatedText)
				if err != nil {
					errCh <- fmt.Errorf("failed to create deck: %w", err)
					return
				}

				// Stream the deck to the channel
				decksCh <- deck

				// Reset the words slice
				words = nil

				// If the current word *is* the marker ("---"), skip adding it to the next chunk
				if word == "---" {
					continue
				}
			}

			words = append(words, word)
		}

		if err := scanner.Err(); err != nil {
			errCh <- fmt.Errorf("error while reading file: %w", err)
			return
		}

		if len(words) > 0 {
			accumulatedText := strings.Join(words, " ")
			deck, err := createDeck(ctx, accumulatedText)
			if err != nil {
				errCh <- fmt.Errorf("failed to create deck from leftover words: %w", err)
				return
			}
			decksCh <- deck
		}
	}()

	return decksCh, errCh
}

// createDeck calls (NewChatCompletion) and parses the JSON response into a Deck.
func createDeck(ctx context.Context, text string) (Deck, error) {
	logger := logging.FromContext(ctx)
	result, err := NewChatCompletion(ctx, text)
	if err != nil {
		return Deck{}, fmt.Errorf("failed to create a new chat completion: %w", err)
	}
	if result == nil || len(result.Choices) == 0 {
		logger.Error("Failed to generate flashcards or received empty response")
		return Deck{}, fmt.Errorf("failed to generate flashcards or received empty response")
	}

	rawOutput := result.Choices[0].Message.Content
	newDeck := Deck{}
	err = json.Unmarshal([]byte(rawOutput), &newDeck)
	if err != nil {
		logger.Errorf("Failed to parse flashcards JSON: %v", err)
		return Deck{}, fmt.Errorf("invalid JSON response from transform package")
	}
	return newDeck, nil
}

// Reads a deck from the deck channel and appends it to one final deck. Will block until the deck channel is closed.
// Depends on streamDocument
func joinDeck(decksCh <-chan Deck) (Deck, error) {
	joinedDeck := Deck{
		Title: "",
		Cards: []Flashcards{},
	}

	// Keep track of whether we’ve set a title yet.
	firstDeck := true

	for deck := range decksCh {
		// If this is the first deck we see, adopt its title.
		if firstDeck {
			joinedDeck.Title = deck.Title
			firstDeck = false
		}
		// Merge flashcards
		joinedDeck.Cards = append(joinedDeck.Cards, deck.Cards...)
	}
	return joinedDeck, nil
}

func TransformNote(ctx context.Context, docPath string) (Deck, error) {
	deckChan, errChan := streamDocument(ctx, docPath)
	deck, err := joinDeck(deckChan)
	if err != nil {
		return Deck{}, err
	}
	if err = <-errChan; err != nil {
		return Deck{}, err
	}
	return deck, err
}
