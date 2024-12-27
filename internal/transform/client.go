package transform

import (
	"context"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func newClient() *openai.Client {
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)
	return client
}

func NewChatCompletion() *openai.Completion {
	client := newClient()
	configs := DefaultCompletionConfigs()

	completion, err := client.Completions.New(context.TODO(), configs)
	if err != nil {
		return nil
	}

	return completion
}
