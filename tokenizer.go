package main

import "errors"

var ErrNoClosingQuote = errors.New("missing closing quote")

func tokenizeInput(input string) ([]string, error) {
	var (
		start        int      = 0
		end          int      = 0
		tokens       []string = make([]string, 0)
		shouldEscape bool     = false
	)

	// om nom nom
	consumeToken := func(start, end int) {
		tokens = append(tokens, input[start:end])
	}

	for end < len(input) {
		if shouldEscape {
			end++
			shouldEscape = false
			continue
		}

		switch input[end] {
		// spaces separate tokens
		case ' ':
			if start != end {
				consumeToken(start, end)
			}
			start = end + 1
		// escape certain characters
		case '\\':
			shouldEscape = true
		// consume string token
		case '\'':
			fallthrough
		case '"':
			quoteChar := input[end]
			end++
			for end < len(input) && input[end] != quoteChar {
				end++
			}
			if end == len(input) {
				return []string{}, ErrNoClosingQuote
			}
			if input[start] != quoteChar {
				// this case is for something like cmd var="askdjfhkdfhj"
				consumeToken(start, end+1)
			} else {
				// do not consume leading quote
				consumeToken(start+1, end)
			}
			start = end + 1
		}
		end++
	}

	// consume last token if it exists
	if start != end {
		consumeToken(start, end)
	}

	return tokens, nil
}
