package service

import "testing"

func TestParseWordsResponse(t *testing.T) {
	words, err := ParseWordsResponse(`{"words":[{"english":"ability","korean":"\ub2a5\ub825"},{"english":"accept","korean":"\uc218\ub77d\ud558\ub2e4"}]}`)
	if err != nil {
		t.Fatalf("ParseWordsResponse returned error: %v", err)
	}

	if len(words) != 2 {
		t.Fatalf("expected 2 words, got %d", len(words))
	}

	if words[0].English != "ability" || words[0].Korean != "\ub2a5\ub825" {
		t.Fatalf("unexpected first word: %+v", words[0])
	}
}

func TestParseWordsResponseInvalidJSON(t *testing.T) {
	_, err := ParseWordsResponse(`{"words":`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseWordsResponseWithCodeFence(t *testing.T) {
	words, err := ParseWordsResponse("```json\n{\"words\":[{\"english\":\"ability\",\"korean\":\"\\ub2a5\\ub825\"}]}\n```")
	if err != nil {
		t.Fatalf("ParseWordsResponse returned error: %v", err)
	}

	if len(words) != 1 {
		t.Fatalf("expected 1 word, got %d", len(words))
	}

	if words[0].English != "ability" {
		t.Fatalf("unexpected word: %+v", words[0])
	}
}
