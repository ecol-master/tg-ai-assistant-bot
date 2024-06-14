package openai

import (
	"context"
	"planpilot/internal/config"
	"planpilot/internal/logger"

	"github.com/sashabaranov/go-openai"
)

func MakeOpenAICall(userInput string) string {
	client := openai.NewClient(config.New().OPENAI_TOKEN)
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4TurboPreview,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userInput,
			},
		},
	},
	)

	if err != nil {
		logger.Error("ChatCompletion error: ", err)
		return "Error"
	}

	return resp.Choices[0].Message.Content
}
