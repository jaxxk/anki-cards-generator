package transform

import (
	"testing"

	"github.com/openai/openai-go"
)

func TestGenerateSchema(t *testing.T) {
	schema := generateSchema[Deck]()
	if schema == nil {
		t.Error("Expected schema, got nil")
	}
}

func TestCreateResponseSchema(t *testing.T) {
	responseSchema := CreateResponseSchema()
	if responseSchema.Name != openai.F("deck") {
		t.Errorf("Expected name 'deck', got %v", responseSchema.Name)
	}
}
