package main

func parseSquareParen(expression string, position int) (regexToken, int, error) {
	expressionPart := expression[position:]
	parsing := false
	charSequence := ""
	offset := 0
	for i, ch := range expressionPart {
		if ch == ']' {
			parsing = false
			offset = i + 1
			break
		}
		if parsing {
			charSequence = charSequence + string(ch)
		}
		if ch == '[' {
			parsing = true
		}
	}
	negated := charSequence[0] == '^'
	body := charSequence
	if negated {
		body = charSequence[1:]
	}
	newToken := regexToken{
		typeOfToken: CharacterClass,
		body:        makeClassSet(body),
		negated:     negated,
	}
	return newToken, offset, nil
}

func parseParen(expression string, position int) (optionStack, int) {
	charSequence := ""
	parsing := false
	expressionPart := expression[position:]
	offset := 0
	for i, ch := range expressionPart {
		if ch == ')' {
			parsing = false
			offset = i
			break
		}
		if parsing {
			charSequence = charSequence + string(ch)
		}
		if ch == '(' {
			parsing = true
		}

	}

	optionStack, _ := parseExpression(charSequence)
	return optionStack, offset
}
