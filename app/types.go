package main

const numeric = "0123456789"
const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"

type tokenType int

const (
	LiteralString    tokenType = iota
	CharacterClass   tokenType = iota
	RepeatOnceOrMore tokenType = iota
	RepeatZeroOrMore tokenType = iota
	Optional         tokenType = iota
	Wildcard         tokenType = iota
	Anchor           tokenType = iota
	Group            tokenType = iota
	NonExistent      tokenType = iota

	// Other types?
)

type regexToken struct {
	typeOfToken   tokenType
	char          rune
	runes         []rune
	body          map[rune]bool
	negated       bool
	anchoredStart bool
	anchoredEnd   bool
	alternatives  optionStack
}

type regexTokens []regexToken
type optionStack []regexTokens
