package transform

import (
	"context"
	"os"

	"github.com/jaxxk/anki-cards-generator/pkg/logging"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func newClient() *openai.Client {
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)
	return client
}

// NewChatCompletion creates a new chat completion request to OpenAI using the provided context and input data.
// ctx: the request context for handling timeouts and cancellations.
// promptData: the input string appended to the default prompt to the OpenAI API.
func NewChatCompletion(ctx context.Context, promptData string) *openai.Completion {
	client := newClient()
	configs := DefaultCompletionConfigs(promptData)

	logger := logging.FromContext(ctx)
	completion, err := client.Completions.New(ctx, configs)
	if err != nil {
		logger.Errorf("error: %v", err)
		return nil
	}

	return completion
}
