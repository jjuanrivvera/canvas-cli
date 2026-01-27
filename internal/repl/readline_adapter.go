package repl

import (
	"strings"

	"github.com/chzyer/readline"
)

// Verify interface compliance at compile time.
var _ readline.AutoCompleter = (*ReadlineCompleter)(nil)

// ReadlineCompleter adapts a Completer to readline's AutoCompleter interface.
type ReadlineCompleter struct {
	completer *Completer
}

// NewReadlineCompleter wraps a Completer for use with readline.
func NewReadlineCompleter(c *Completer) *ReadlineCompleter {
	return &ReadlineCompleter{completer: c}
}

// Do implements readline.AutoCompleter. It returns completion candidates
// and the length of the word prefix being completed.
func (rc *ReadlineCompleter) Do(line []rune, pos int) ([][]rune, int) {
	// Use only text up to cursor position
	input := string(line[:pos])

	candidates := rc.completer.Complete(input)
	if len(candidates) == 0 {
		return nil, 0
	}

	// Determine the word being completed (last token before cursor)
	lastSpace := strings.LastIndex(input, " ")
	prefix := ""
	if lastSpace >= 0 && lastSpace < len(input)-1 {
		prefix = input[lastSpace+1:]
	} else if lastSpace < 0 {
		prefix = input
	}
	// If input ends with space, prefix is empty (completing new word)

	prefixLen := len([]rune(prefix))

	result := make([][]rune, 0, len(candidates))
	for _, c := range candidates {
		// Return only the suffix after the shared prefix
		if strings.HasPrefix(c, prefix) {
			suffix := c[len(prefix):]
			result = append(result, []rune(suffix+" "))
		} else {
			result = append(result, []rune(c+" "))
		}
	}

	return result, prefixLen
}
