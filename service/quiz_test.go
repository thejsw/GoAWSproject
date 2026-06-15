package service

import (
	"testing"

	"go-openai-server/models"
)

func TestParseQuizResponse(t *testing.T) {
	response, err := ParseQuizResponse(`{"type":"weekly","week_num":202623,"questions":[{"question_text":"Choose the best answer.","choices":[{"key":"A","text":"goes","correct":true},{"key":"B","text":"going","correct":false}],"explanations":[{"language":"ko","text":"정답 설명"},{"language":"ja","text":"解説"}]}]}`)
	if err != nil {
		t.Fatalf("ParseQuizResponse returned error: %v", err)
	}

	if response.Type != "weekly" {
		t.Fatalf("expected type weekly, got %s", response.Type)
	}

	if response.WeekNum == nil || *response.WeekNum != 202623 {
		t.Fatalf("expected week_num 202623, got %+v", response.WeekNum)
	}

	if len(response.Questions) != 1 {
		t.Fatalf("expected 1 question, got %d", len(response.Questions))
	}
}

func TestParseQuizResponseWithCodeFence(t *testing.T) {
	response, err := ParseQuizResponse("```json\n{\"type\":\"weekly\",\"week_num\":202623,\"questions\":[{\"question_text\":\"Choose the best answer.\",\"choices\":[{\"key\":\"A\",\"text\":\"goes\",\"correct\":true}],\"explanations\":[{\"language\":\"ko\",\"text\":\"정답 설명\"}]}]}\n```")
	if err != nil {
		t.Fatalf("ParseQuizResponse returned error: %v", err)
	}

	if response.WeekNum == nil || *response.WeekNum != 202623 {
		t.Fatalf("expected week_num 202623, got %+v", response.WeekNum)
	}
}

func TestBuildQuizBundle(t *testing.T) {
	next := func() (string, error) {
		return "fixed-id", nil
	}

	response := &models.OpenAIQuizResponse{
		Type:    "weekly",
		WeekNum: intPtr(202623),
		Questions: []models.OpenAIQuizQuestion{
			{
				QuestionText: "Choose the best answer.",
				Choices: []models.OpenAIQuizChoice{
					{Key: "A", Text: "goes", Correct: true},
					{Key: "B", Text: "going", Correct: false},
				},
				Explanations: []models.OpenAIQuizExplanation{
					{Language: "ko", Text: "정답 설명"},
				},
			},
		},
	}

	bundle, err := BuildQuizBundle(response, next)
	if err != nil {
		t.Fatalf("BuildQuizBundle returned error: %v", err)
	}

	if bundle.Quiz.ID != "fixed-id" {
		t.Fatalf("expected quiz id fixed-id, got %s", bundle.Quiz.ID)
	}

	if len(bundle.Questions) != 1 {
		t.Fatalf("expected 1 question, got %d", len(bundle.Questions))
	}

	if len(bundle.Choices) != 2 {
		t.Fatalf("expected 2 choices, got %d", len(bundle.Choices))
	}

	if len(bundle.Explanations) != 1 {
		t.Fatalf("expected 1 explanation, got %d", len(bundle.Explanations))
	}
}

func intPtr(v int) *int {
	return &v
}
