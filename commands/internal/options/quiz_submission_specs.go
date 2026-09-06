package options

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseQuestionScores parses "<question-id>=<score>" entries into a map.
func ParseQuestionScores(specs []string) (map[int64]float64, error) {
	scores := make(map[int64]float64, len(specs))
	for _, spec := range specs {
		id, value, err := splitQuestionSpec("question-score", spec)
		if err != nil {
			return nil, err
		}
		score, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid question-score %q: score must be a number", spec)
		}
		scores[id] = score
	}
	return scores, nil
}

// ParseQuestionComments parses "<question-id>=<text>" entries into a map.
func ParseQuestionComments(specs []string) (map[int64]string, error) {
	comments := make(map[int64]string, len(specs))
	for _, spec := range specs {
		id, value, err := splitQuestionSpec("question-comment", spec)
		if err != nil {
			return nil, err
		}
		comments[id] = value
	}
	return comments, nil
}

func splitQuestionSpec(flag, spec string) (int64, string, error) {
	parts := strings.SplitN(spec, "=", 2)
	if len(parts) != 2 || parts[0] == "" {
		return 0, "", fmt.Errorf("invalid %s %q: expected <question-id>=<value>", flag, spec)
	}
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id <= 0 {
		return 0, "", fmt.Errorf("invalid %s %q: question ID must be a positive integer", flag, spec)
	}
	return id, parts[1], nil
}
