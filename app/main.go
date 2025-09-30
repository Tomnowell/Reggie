package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	expression := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}
	text := string(line)

	ok, err := matchLine(expression, text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	fmt.Println("Result: ", ok)
	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchTokens(tokens regexTokens, text []rune, position int) bool {
	// Base case:
	if len(tokens) == 0 {
		return true
	}

	token := tokens[0]
	remaining := tokens[1:]

	switch token.typeOfToken {
	case LiteralString:
		if position+len(token.runes) > len(text) {
			return false
		}
		for i, ch := range token.runes {
			if ch != '.' && text[position+i] != ch {
				return false
			}
		}
		return matchTokens(remaining, text, position+len(token.runes))

	case Wildcard:
		if position >= len(text) {
			return false
		}
		return matchTokens(remaining, text, position+1)

	case CharacterClass:
		if position >= len(text) {
			return false
		}
		present := token.body[text[position]]
		if token.negated {
			present = !present
		}

		if !present {
			return false
		}
		return matchTokens(remaining, text, position+1)
	case Optional:
		if matchTokens(remaining, text, position) {
			return true
		}
		if position < len(text) && (token.char == '.' || token.char == text[position]) {
			return matchTokens(remaining, text, position+1)
		}
		return false
	case RepeatOnceOrMore:
		// Need at least one character available
		if position >= len(text) {
			return false
		}
		max := len(text) - position
		for parsed := 1; parsed <= max; parsed++ {
			present := true
			for i := 0; i < parsed; i++ {
				if !matches(text[position+i], token) {
					present = false
					break
				}
			}
			if !present {
				break
			}
			if matchTokens(remaining, text, position+parsed) {
				return true
			}
		}
		return false

	case RepeatZeroOrMore:
		max := len(text) - position
		for parsed := 0; parsed <= max; parsed++ {
			// parsed==0 means consume nothing; otherwise validate run
			if parsed > 0 && token.char != '.' {
				if !matches(text[position+parsed-1], token) {
					break
				}
			}
			if matchTokens(remaining, text, position+parsed) {
				return true
			}
		}
		return false

	case Anchor:
		if token.anchoredStart {
			if position != 0 {
				return false
			}
			return matchTokens(remaining, text, position)
		}
		if token.anchoredEnd {
			if position != len(text) {
				return false
			} else {
				return true
			}
		}
	case Group:
		remaining := tokens[1:]
		for _, alt := range token.alternatives {
			if matchTokens(append(alt, remaining...), text, position) {
				return true
			}
		}
		return false

	default:
		return false
	}
	return false
}

func matchLine(expression, text string) (bool, error) {
	if len(expression) < 1 {
		return false, fmt.Errorf("unsupported pattern: %q", expression)
	}

	optionStack, err := parseExpression(expression)
	if err != nil {
		return false, err
	}

	runeText := []rune(text)
	// Try each option in the stack
	for _, expressionList := range optionStack {
		if len(expressionList) == 0 {
			continue
		}

		// anchored start → must begin at 0
		if expressionList[0].anchoredStart {
			return matchTokens(expressionList[1:], runeText, 0), nil
		}

		// otherwise, slide across the string
		for start := 0; start < len(runeText); start++ {
			if matchTokens(expressionList, runeText[start:], 0) {
				return true, nil
			}
		}

	}
	return false, nil
}
