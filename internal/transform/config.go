package transform

import "github.com/openai/openai-go"

// DefaultModel defines the OpenAI model to use for flashcard generation.
var DefaultModel openai.ChatModel = openai.ChatModelGPT4oMini

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

4. Output Format: Return the flashcards in a JSON array, where each flashcard is an object with the keys 'front' (the question) and 'back' (the explanation).

Example Structure:
[
    {
        "front": "How does [concept 1] influence the behavior of [concept 2], and what are the practical implications of this relationship?",
        "back": "[Explanation that ties the two concepts together, explaining how one influences the other, and why it matters in practical terms]"
    },
    {
        "front": "What are the key differences between [concept 1] and [concept 2], and how do these differences impact [related concept or process]?",
        "back": "[Detailed explanation comparing the two concepts and their practical impacts]"
    }
]
`

// DefaultCompletionConfigs constructs the OpenAI CompletionNewParams for the given input text.
// inputText: The content to be processed for generating flashcards.
// Returns: OpenAI CompletionNewParams with the configured parameters.
func DefaultCompletionConfigs(inputText string) openai.CompletionNewParams {
	// Concatenate prompt and input text
	prompt := DefaultPrompt + "\n" + inputText

	// Construct the parameters
	params := openai.CompletionNewParams{
		Model:            openai.F(openai.CompletionNewParamsModel(DefaultModel)),
		Prompt:           openai.F(openai.CompletionNewParamsPromptUnion(openai.CompletionNewParamsPromptArrayOfStrings{prompt})),
		BestOf:           openai.Int(1),
		Temperature:      openai.Float(0.7),
		FrequencyPenalty: openai.Float(0.5),
		PresencePenalty:  openai.Float(0.5),
	}

	return params
}
