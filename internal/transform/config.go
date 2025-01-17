package transform

import "github.com/openai/openai-go"

// DefaultModel defines the OpenAI model to use for flashcard generation.
var DefaultModel openai.ChatModel = openai.ChatModelGPT4oMini
var DefaultFrequencyPenalty float64 = 1.2
var DefaultPresencePenalty float64 = 1.2

// DefaultPrompt is the base prompt for generating flashcards.
var DefaultPrompt string = `
You are a specialized flashcard generator. Your task is to process a .md or .txt file containing detailed information and produce a series of flashcards in strict JSON format. Each flashcard must include:

1. "front": A question that either:
   - Challenges deeper analysis (showing relationships between concepts), or
   - Tests quick recall of fundamental facts.
2. "back": A comprehensive explanation that integrates relevant details from the content. Include validated Python or Go code examples if they add clarity.

Output Requirements:
- Return only a JSON array of flashcards.
- Do not include any text, explanations, or formatting outside the JSON structure.

The final output must look like this:

[
  {
    "front": "Some question here",
    "back": "Some explanation here with optional code snippets"
  },
  {
    "front": "...",
    "back": "..."
  }
]

Do not deviate from this format.
`

// DefaultChatCompletionConfigs constructs the OpenAI ChatCompletionNewParams for the given input text.
// inputText: The content to be processed for generating flashcards.
// Returns: OpenAI ChatCompletionNewParams with the configured parameters.
func DefaultChatCompletionConfigs(inputText string) openai.ChatCompletionNewParams {
	responseSchema := CreateResponseSchema()
	// Construct the parameters
	params := openai.ChatCompletionNewParams{
		Model: openai.F(DefaultModel),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.ChatCompletionDeveloperMessageParam{
				Role: openai.F(openai.ChatCompletionDeveloperMessageParamRoleDeveloper),
				Content: openai.F([]openai.ChatCompletionContentPartTextParam{
					openai.TextPart(DefaultPrompt),
				}),
			},
			openai.UserMessage(inputText),
		}),
		FrequencyPenalty: openai.Float(DefaultFrequencyPenalty),
		// only have 1 chat completion choice
		N:               openai.Int(1),
		PresencePenalty: openai.Float(DefaultPresencePenalty),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(responseSchema),
			},
		),
	}
	return params
}
