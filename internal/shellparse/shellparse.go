package shellparse

// Split splits a shell-like input string into arguments, respecting
// double-quoted and single-quoted strings. Multiple consecutive spaces
// are treated as a single delimiter. Quotes are consumed (not included
// in the output).
func Split(input string) []string {
	var args []string
	var current []byte
	inQuote := false
	quoteChar := byte(0)

	for i := 0; i < len(input); i++ {
		c := input[i]

		switch {
		case (c == '"' || c == '\'') && !inQuote:
			inQuote = true
			quoteChar = c
		case inQuote && c == quoteChar:
			inQuote = false
			quoteChar = 0
		case c == ' ' && !inQuote:
			if len(current) > 0 {
				args = append(args, string(current))
				current = current[:0]
			}
		default:
			current = append(current, c)
		}
	}

	if len(current) > 0 {
		args = append(args, string(current))
	}

	return args
}
