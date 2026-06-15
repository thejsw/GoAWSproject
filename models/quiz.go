package models

type OpenAIQuizResponse struct {
	Type      string               `json:"type"`
	WeekNum   *int                 `json:"week_num"`
	Questions []OpenAIQuizQuestion `json:"questions"`
}

type OpenAIQuizQuestion struct {
	QuestionText string                  `json:"question_text"`
	Choices      []OpenAIQuizChoice      `json:"choices"`
	Explanations []OpenAIQuizExplanation `json:"explanations"`
}

type OpenAIQuizChoice struct {
	Key     string `json:"key"`
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
}

type OpenAIQuizExplanation struct {
	Language string `json:"language"`
	Text     string `json:"text"`
}

type Quiz struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	DayNum  *int   `json:"day_num,omitempty"`
	WeekNum *int   `json:"week_num,omitempty"`
}

type QuizQuestion struct {
	ID           string `json:"id"`
	QuizID       string `json:"quiz_id"`
	QuestionText string `json:"question_text"`
	OrderIndex   int    `json:"order_index"`
}

type QuizChoice struct {
	ID         string `json:"id"`
	QuestionID string `json:"question_id"`
	ChoiceKey  string `json:"choice_key"`
	ChoiceText string `json:"choice_text"`
	IsCorrect  bool   `json:"is_correct"`
	OrderIndex int    `json:"order_index"`
}

type QuizExplanation struct {
	ID          string `json:"id"`
	QuestionID  string `json:"question_id"`
	Language    string `json:"language"`
	Explanation string `json:"explanation"`
	OrderIndex  int    `json:"order_index"`
}

type QuizBundle struct {
	Quiz         Quiz              `json:"quiz"`
	Questions    []QuizQuestion    `json:"questions"`
	Choices      []QuizChoice      `json:"choices"`
	Explanations []QuizExplanation `json:"explanations"`
}
