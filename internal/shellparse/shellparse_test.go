package shellparse

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "simple words",
			input: "courses list",
			want:  []string{"courses", "list"},
		},
		{
			name:  "double quoted argument",
			input: `assignments create --name "Homework 1"`,
			want:  []string{"assignments", "create", "--name", "Homework 1"},
		},
		{
			name:  "single quoted argument",
			input: `config set --value 'hello world'`,
			want:  []string{"config", "set", "--value", "hello world"},
		},
		{
			name:  "mixed quotes",
			input: `search --name "some thing" --desc 'another thing'`,
			want:  []string{"search", "--name", "some thing", "--desc", "another thing"},
		},
		{
			name:  "multiple spaces",
			input: "courses   list   --verbose",
			want:  []string{"courses", "list", "--verbose"},
		},
		{
			name:  "leading and trailing spaces",
			input: "  courses list  ",
			want:  []string{"courses", "list"},
		},
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "only spaces",
			input: "   ",
			want:  nil,
		},
		{
			name:  "empty quoted string produces argument",
			input: `courses ""`,
			want:  []string{"courses"},
		},
		{
			name:  "adjacent quoted and unquoted",
			input: `--name="hello world"`,
			want:  []string{"--name=hello world"},
		},
		{
			name:  "single word",
			input: "help",
			want:  []string{"help"},
		},
		{
			name:  "quote inside other quote type",
			input: `say "it's fine"`,
			want:  []string{"say", "it's fine"},
		},
		{
			name:  "double quote inside single quotes",
			input: `say 'he said "hello"'`,
			want:  []string{"say", `he said "hello"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Split(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
