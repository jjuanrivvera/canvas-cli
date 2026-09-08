package options

import "testing"

func TestParseQuestionScores(t *testing.T) {
	got, err := ParseQuestionScores([]string{"11=2.5", "12=0", "13=-1"})
	if err != nil {
		t.Fatalf("ParseQuestionScores: %v", err)
	}
	if got[11] != 2.5 || got[12] != 0 || got[13] != -1 {
		t.Errorf("ParseQuestionScores = %v", got)
	}
	if _, err := ParseQuestionScores([]string{"0=1"}); err == nil {
		t.Error("expected error for question ID 0")
	}
	if _, err := ParseQuestionScores([]string{"=1"}); err == nil {
		t.Error("expected error for empty question ID")
	}
}

func TestParseQuestionComments(t *testing.T) {
	got, err := ParseQuestionComments([]string{"11=a=b", "12="})
	if err != nil {
		t.Fatalf("ParseQuestionComments: %v", err)
	}
	if got[11] != "a=b" {
		t.Errorf("comment with '=' in value = %q, want %q", got[11], "a=b")
	}
	if v, ok := got[12]; !ok || v != "" {
		t.Errorf("empty comment should be accepted, got %v", got)
	}
}
