package main

func makeSlashToken(ch rune) regexToken {

	token := regexToken{typeOfToken: CharacterClass}
	switch ch {
	case 'd':
		token.body = makeClassSet(numeric)
		token.negated = false

	case 'D':
		token.body = makeClassSet(numeric)
		token.negated = true

	case 'w':
		token.body = makeClassSet(alphanumeric)
		token.negated = false
	case 'W':
		token.body = makeClassSet(alphanumeric)
		token.negated = true
	}
	return token
}

func makeClassSet(s string) map[rune]bool {
	// Chizu is Japanese for map if anyone is wondering :D
	chizu := make(map[rune]bool)
	for _, ch := range s {
		chizu[ch] = true
	}
	return chizu
}

func matches(r rune, token regexToken) bool {
	if token.body != nil {
		present := token.body[r]
		if token.negated {
			present = !present
		}
		return present
	}

	// Wildcard
	if token.char == '.' {
		return true
	}

	// Literal
	return token.char == r
}
