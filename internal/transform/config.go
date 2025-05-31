package transform

import "github.com/openai/openai-go"

// DefaultModel defines the OpenAI model to use for flashcard generation.
var DefaultModel openai.ChatModel = openai.ChatModelGPT4oMini
var DefaultFrequencyPenalty float64 = 1.2
var DefaultPresencePenalty float64 = 1.2

// DefaultPrompt is the base prompt for generating flashcards.
var DefaultPrompt string = `
You are a specialized flashcard generator. Your task is to process a .md or .txt file containing structured technical notes, including math equations and code snippets, and produce a series of high-quality flashcards in strict JSON format.

Each flashcard must include:

1. "front": A question that either:
   - Tests conceptual understanding by requiring reasoning or comparison.
   - Tests factual recall such as definitions, formulas, or terms.

2. "back": A detailed explanation that:
   - Integrates relevant insights from the notes.
   - Preserves math equations and code formatting.

Formatting Rules:
- Use **standard Markdown** to represent formatting:
  - Inline math: wrap with single dollar signs, e.g., $E[R] = \\mu + \\beta F$
  - Block math: wrap with double backslashes and square brackets, e.g., \\[ E[R] = \\mu + \\beta F \\]
  - Code blocks: use triple backticks and specify the language, e.g., 
    ` + "```python\n" + `    def f(x):\n        return x**2\n    ` + "```" + `

Output Requirements:
- Return only a JSON array of flashcards.
- Do not include any text, explanations, or formatting outside the JSON array.
- Use escaped backslashes (\\\\) where necessary in math expressions to ensure LaTeX renders properly.

Example Output:
[
  {
    "front": "What is the formula for the Sharpe Ratio?",
    "back": "The Sharpe Ratio is defined as:\\n\\n$Sharpe = \\\\frac{E[R - R_f]}{\\\\sigma}$\\n\\nWhere:\\n- $R$: Portfolio return\\n- $R_f$: Risk-free rate\\n- $\\\\sigma$: Standard deviation of excess return"
  },
  {
    "front": "How do you define a simple function in Python?",
    "back": "Use the \def keyword:\\n\\n` + "```python\\n" + `def square(x):\\n    return x**2\\n` + "```" + `"
  }
]`

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
