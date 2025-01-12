package transform

import (
	"context"

	"github.com/jaxxk/anki-cards-generator/internal/encryption"
	"github.com/jaxxk/anki-cards-generator/pkg/logging"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"go.uber.org/zap"
)

func newClient(logger *zap.SugaredLogger) (*openai.Client, error) {
	key, err := encryption.GetAPIKey(logger)
	if err != nil {
		return nil, err
	}
	client := openai.NewClient(
		option.WithAPIKey(key),
	)
	return client, nil
}

// NewChatCompletion creates a new chat completion request to OpenAI using the provided context and input data.
// ctx: the request context for handling timeouts and cancellations.
// promptData: the input string appended to the default prompt to the OpenAI API.
func NewChatCompletion(ctx context.Context, promptData string) (*openai.ChatCompletion, error) {
	logger := logging.FromContext(ctx)
	client, err := newClient(logger)
	if err != nil {
		return nil, err
	}
	configs := DefaultChatCompletionConfigs(promptData)
	logger.Infof("Prompt: \n %v \n", configs.Messages.Value[0])
	chatCompletion, err := client.Chat.Completions.New(ctx, configs)
	if err != nil {
		logger.Errorf("configs: %v", configs)
		logger.Errorf("error: %v", err)
		return nil, err
	}

	return chatCompletion, nil
}
