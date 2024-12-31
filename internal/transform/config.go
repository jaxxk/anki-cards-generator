package transform

import "github.com/openai/openai-go"

// DefaultModel defines the OpenAI model to use for flashcard generation.
var DefaultModel openai.ChatModel = openai.ChatModelGPT4oMini
var DefaultFrequencyPenalty float64 = 1.2
var DefaultPresencePenalty float64 = 1.2

// DefaultPrompt is the base prompt for generating flashcards.
var DefaultPrompt string = `
You are a tool that processes a text file (in .md or .txt format) containing detailed information. Your task is to generate insightful flashcards by connecting various concepts within the content. Each flashcard should have a front and back in a JSON format.

Your task is to produce flashcards in a strict JSON format based on the input content. Do not include any text or explanations outside of the JSON structure.

Instructions:
1. Understand the content: Carefully read through the provided content, identify and extract key concepts, processes, relationships, and ideas that are interrelated.
2. Create insightful questions: For each flashcard, the front should pose a question that challenges the reader to understand and connect different ideas from the content. The questions should:
   - Make connections between related concepts.
   - Ask how different concepts work together or influence each other.
   - Encourage deeper thinking or analysis, rather than simple fact recall.
3. Provide detailed answers: The back should provide a comprehensive explanation that answers the question on the front. The explanation should integrate information from the content, show the relationship between concepts, and elaborate on their significance.
4. Include validated code examples: If the text file you are processing contains code examples, verify their correctness. If they are valid, append them to the back of the flashcard. If the code is incorrect, correct it while preserving the original intent. Clearly mark or replace the erroneous parts with corrected versions.
5. Output Format: Return the flashcards in a JSON array, where each flashcard is an object with the keys 'front' (the question) and 'back' (the explanation, along with any validated code examples in python or golang).
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
