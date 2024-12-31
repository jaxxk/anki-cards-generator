package transform

import (
	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
)

// Flashcards represents a single flashcard with a front and back.
type Flashcards struct {
	Front string `json:"front" jsonschema_description:"The front side of the flashcard"`
	Back  string `json:"back" jsonschema_description:"The back side of the flashcard"`
}

// Deck represents a collection of flashcards.
type Deck struct {
	Title string       `json:"Title" jsonschema_description:"The title of the deck"`
	Cards []Flashcards `json:"cards" jsonschema_description:"A deck consisting of flashcards"`
}

func (deck *Deck) UpdateTitle(title string) {
	deck.Title = title
}

// generateSchema generates a JSON schema for the provided type.
func generateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

// CreateResponseSchema creates a JSON schema parameter for the OpenAI API response format.
func CreateResponseSchema() openai.ResponseFormatJSONSchemaJSONSchemaParam {
	deckSchema := generateSchema[Deck]()
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("deck"),
		Description: openai.F("A deck consisting of flashcards with questions and answers"),
		Schema:      openai.F(deckSchema),
		Strict:      openai.Bool(true),
	}
	return schemaParam
}
