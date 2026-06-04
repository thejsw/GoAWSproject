package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"go-openai-server/models"
	"go-openai-server/repository"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func Generate(
	ctx context.Context,
	topic string,
	repo *repository.WordRepository,
) (
	[]models.Word,
	error,
) {

	key :=
		os.Getenv(
			"OPENAI_API_KEY",
		)

	client :=
		openai.NewClient(
			option.WithAPIKey(
				key,
			),
		)

	resp, err :=
		client.Chat.Completions.New(
			ctx,
			openai.ChatCompletionNewParams{
				Model: "gpt-4o",
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.UserMessage(topic),
				},
			},
		)

	if err != nil {
		return nil, err
	}

	if os.Getenv("DEBUG_OPENAI_RESPONSE") == "true" {
		log.Printf("OpenAI raw response: %s", resp.Choices[0].Message.Content)
	}

	words, err := ParseWordsResponse(
		resp.Choices[0].Message.Content,
	)
	if err != nil {
		return nil, err
	}

	if err := repo.InsertMany(ctx, words); err != nil {
		return nil, err
	}

	return words, nil
}

func ParseWordsResponse(content string) ([]models.Word, error) {
	var response models.GenerateResponse

	normalized, err := normalizeJSONContent(content)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(normalized), &response); err != nil {
		return nil, err
	}

	return response.Words, nil
}

func normalizeJSONContent(content string) (string, error) {
	trimmed := strings.TrimSpace(content)

	if strings.HasPrefix(trimmed, "```") {
		firstNewline := strings.Index(trimmed, "\n")
		if firstNewline == -1 {
			return "", fmt.Errorf("invalid fenced JSON content")
		}

		trimmed = strings.TrimSpace(trimmed[firstNewline+1:])

		if strings.HasPrefix(trimmed, "json") {
			trimmed = strings.TrimSpace(trimmed[len("json"):])
		}

		if strings.HasSuffix(trimmed, "```") {
			trimmed = strings.TrimSpace(trimmed[:len(trimmed)-3])
		}
	}

	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start == -1 || end == -1 || end < start {
		return "", fmt.Errorf("no json object found in response")
	}

	return trimmed[start : end+1], nil
}
