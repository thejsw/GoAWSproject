package service

import (
	"context"
	"fmt"
	"os"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func Generate(topic string) (string, error) {

	key :=
		os.Getenv(
			"OPENAI_API_KEY",
		)

	fmt.Println("KEY:", key)

	client :=
		openai.NewClient(
			option.WithAPIKey(
				key,
			),
		)

	resp, err :=
		client.Chat.Completions.New(
			context.Background(),
			openai.ChatCompletionNewParams{
				Model: "gpt-4o",
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.UserMessage(topic),
				},
			},
		)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
