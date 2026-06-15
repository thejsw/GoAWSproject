package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"go-openai-server/models"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/shared/constant"
)

type QuizRepository interface {
	SaveQuizBundle(ctx context.Context, bundle models.QuizBundle) error
}

type UUIDGenerator func() (string, error)

func GeneratePart5Quiz(
	ctx context.Context,
	topic string,
	repo QuizRepository,
) (*models.QuizBundle, error) {
	key := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(
		option.WithAPIKey(key),
	)

	resp, err := client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: "gpt-4o-mini",
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
					Type: constant.JSONSchema("json_schema"),
					JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
						Name:        "toeic_part5_quiz",
						Description: openai.String("Return a TOEIC Part 5 quiz object only."),
						Strict:      openai.Bool(true),
						Schema: map[string]any{
							"type":                 "object",
							"additionalProperties": false,
							"required":             []string{"type", "week_num", "questions"},
							"properties": map[string]any{
								"type": map[string]any{
									"type": "string",
									"enum": []string{"weekly", "daily"},
								},
								"week_num": map[string]any{
									"type": "integer",
								},
								"questions": map[string]any{
									"type":     "array",
									"minItems": 1,
									"items": map[string]any{
										"type":                 "object",
										"additionalProperties": false,
										"required":             []string{"question_text", "choices", "explanations"},
										"properties": map[string]any{
											"question_text": map[string]any{
												"type":      "string",
												"minLength": 1,
											},
											"choices": map[string]any{
												"type":     "array",
												"minItems": 2,
												"items": map[string]any{
													"type":                 "object",
													"additionalProperties": false,
													"required":             []string{"key", "text", "correct"},
													"properties": map[string]any{
														"key": map[string]any{
															"type": "string",
															"enum": []string{"A", "B", "C", "D", "E"},
														},
														"text": map[string]any{
															"type":      "string",
															"minLength": 1,
														},
														"correct": map[string]any{
															"type": "boolean",
														},
													},
												},
											},
											"explanations": map[string]any{
												"type":     "array",
												"minItems": 1,
												"items": map[string]any{
													"type":                 "object",
													"additionalProperties": false,
													"required":             []string{"language", "text"},
													"properties": map[string]any{
														"language": map[string]any{
															"type": "string",
															"enum": []string{"ko", "ja", "en"},
														},
														"text": map[string]any{
															"type":      "string",
															"minLength": 1,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(
					`You are generating a TOEIC Part 5 quiz payload.

Output rules:
- Return exactly one JSON object.
- Do not wrap the JSON in Markdown, code fences, or commentary.
- Do not add keys that are not in the schema.
- Every question must contain at least 2 choices.
- Exactly one choice per question must have "correct": true.
- Every question must contain at least one explanation.
- Use only these languages for explanations: ko, ja, en.
- If you cannot satisfy the schema, output an empty JSON object is not allowed; instead, fix the content before answering.

Content rules:
- "type" should usually be "weekly" unless the prompt explicitly requires another supported value.
- "week_num" must be a six-digit integer in YYYYWW format.
- "choice.key" must be one of A, B, C, D, E.
- "question_text", "choice.text", and explanation text must be plain text only.
`,
				),
				openai.UserMessage(
					`Create a TOEIC Part 5 mock quiz for the topic below.

Topic:
` + topic + `

Return JSON only.`,
				),
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("openai chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}

	content := resp.Choices[0].Message.Content

	if os.Getenv("DEBUG_OPENAI_RESPONSE") == "true" {
		log.Printf("OpenAI raw response: %s", content)
	}

	quizResponse, err := ParseQuizResponse(content)
	if err != nil {
		return nil, err
	}

	bundle, err := BuildQuizBundle(quizResponse, generateUUID)
	if err != nil {
		return nil, err
	}

	if repo == nil {
		return nil, fmt.Errorf("quiz repository is not initialized")
	}

	if err := repo.SaveQuizBundle(ctx, bundle); err != nil {
		return nil, err
	}

	return &bundle, nil
}

func ParseQuizResponse(content string) (*models.OpenAIQuizResponse, error) {
	normalized, err := normalizeJSONContent(content)
	if err != nil {
		return nil, err
	}

	var response models.OpenAIQuizResponse
	if err := json.Unmarshal([]byte(normalized), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func BuildQuizBundle(
	response *models.OpenAIQuizResponse,
	nextID UUIDGenerator,
) (models.QuizBundle, error) {
	if response == nil {
		return models.QuizBundle{}, fmt.Errorf("quiz response is nil")
	}

	if nextID == nil {
		nextID = generateUUID
	}

	if strings.TrimSpace(response.Type) == "" {
		return models.QuizBundle{}, fmt.Errorf("quiz type is empty")
	}

	if response.WeekNum == nil || *response.WeekNum <= 0 {
		return models.QuizBundle{}, fmt.Errorf("week_num must be greater than zero")
	}

	if len(response.Questions) == 0 {
		return models.QuizBundle{}, fmt.Errorf("questions are empty")
	}

	quizID, err := nextID()
	if err != nil {
		return models.QuizBundle{}, fmt.Errorf("generate quiz id failed: %w", err)
	}

	quiz := models.Quiz{
		ID:     quizID,
		Type:   response.Type,
		DayNum: nil,
	}
	quiz.WeekNum = response.WeekNum

	bundle := models.QuizBundle{
		Quiz: quiz,
	}

	for qIndex, question := range response.Questions {
		if strings.TrimSpace(question.QuestionText) == "" {
			return models.QuizBundle{}, fmt.Errorf("question_text is empty at question %d", qIndex+1)
		}

		if len(question.Choices) == 0 {
			return models.QuizBundle{}, fmt.Errorf("choices are empty at question %d", qIndex+1)
		}

		if len(question.Explanations) == 0 {
			return models.QuizBundle{}, fmt.Errorf("explanations are empty at question %d", qIndex+1)
		}

		questionID, err := nextID()
		if err != nil {
			return models.QuizBundle{}, fmt.Errorf("generate question id failed: %w", err)
		}

		bundle.Questions = append(bundle.Questions, models.QuizQuestion{
			ID:           questionID,
			QuizID:       quizID,
			QuestionText: question.QuestionText,
			OrderIndex:   qIndex + 1,
		})

		correctCount := 0
		for cIndex, choice := range question.Choices {
			if strings.TrimSpace(choice.Key) == "" {
				return models.QuizBundle{}, fmt.Errorf("choice key is empty at question %d choice %d", qIndex+1, cIndex+1)
			}

			if strings.TrimSpace(choice.Text) == "" {
				return models.QuizBundle{}, fmt.Errorf("choice text is empty at question %d choice %d", qIndex+1, cIndex+1)
			}

			if choice.Correct {
				correctCount++
			}

			choiceID, err := nextID()
			if err != nil {
				return models.QuizBundle{}, fmt.Errorf("generate choice id failed: %w", err)
			}

			bundle.Choices = append(bundle.Choices, models.QuizChoice{
				ID:         choiceID,
				QuestionID: questionID,
				ChoiceKey:  choice.Key,
				ChoiceText: choice.Text,
				IsCorrect:  choice.Correct,
				OrderIndex: cIndex + 1,
			})
		}

		if correctCount != 1 {
			return models.QuizBundle{}, fmt.Errorf("expected exactly one correct choice at question %d, got %d", qIndex+1, correctCount)
		}

		for eIndex, explanation := range question.Explanations {
			if strings.TrimSpace(explanation.Language) == "" {
				return models.QuizBundle{}, fmt.Errorf("explanation language is empty at question %d explanation %d", qIndex+1, eIndex+1)
			}

			if strings.TrimSpace(explanation.Text) == "" {
				return models.QuizBundle{}, fmt.Errorf("explanation text is empty at question %d explanation %d", qIndex+1, eIndex+1)
			}

			explanationID, err := nextID()
			if err != nil {
				return models.QuizBundle{}, fmt.Errorf("generate explanation id failed: %w", err)
			}

			bundle.Explanations = append(bundle.Explanations, models.QuizExplanation{
				ID:          explanationID,
				QuestionID:  questionID,
				Language:    explanation.Language,
				Explanation: explanation.Text,
				OrderIndex:  eIndex + 1,
			})
		}
	}

	return bundle, nil
}

func generateUUID() (string, error) {
	var b [16]byte

	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	), nil
}
