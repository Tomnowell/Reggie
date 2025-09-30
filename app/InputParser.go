package main

import (
	"fmt"
)

func parseExpression(expression string) (optionStack, error) {
	var expressionList = regexTokens{}
	var optionStack = optionStack{}
	var newToken = regexToken{}
	newToken.typeOfToken = NonExistent
	var err error
	escaped := false
	var literalBuffer []rune
	expressionRunes := []rune(expression)

	flush := func() {
		if len(literalBuffer) > 0 {
			expressionList = append(expressionList, regexToken{
				typeOfToken: LiteralString,
				runes:       literalBuffer,
			})
			literalBuffer = nil
		}
	}

	for i:=0; i < len(expressionRunes); i++ {
		ch := expressionRunes[i]
		// We're wrapping the switch in an anon func to use defer each time to append the token
		// only once per loop...I hope that makes sense: sorry!
		switch expression[i] {
		case '\\':
			if escaped {
				literalBuffer = append(literalBuffer, ch)
				escaped = false
			} else {
				escaped = true
			}
		case 'd', 'D', 'w', 'W':
			if escaped {
				// 'Escaped' actually means proceeded by backslash so in this case
				// it's a command. There must be a better solution: sorry, future me!
				flush()
				newToken = makeSlashToken(ch)
				expressionList = append(expressionList, newToken)
				newToken = regexToken{}
				escaped = false
			} else {
				// Just a normal d D or w W
				literalBuffer = append(literalBuffer, ch)
			}
		case '[':
			if escaped {
				literalBuffer = append(literalBuffer, ch)
				escaped = false
			} else {
				flush()
				offset := 0
				newToken, offset, err = parseSquareParen(expression, i)
				if err != nil {
					return optionStack, err
				}
				newToken.typeOfToken = CharacterClass
				expressionList = append(expressionList, newToken)
				newToken = regexToken{}
				i += offset
			}
		case ']':
			if escaped {
				literalBuffer = append(literalBuffer, ch)
				escaped = false
			} else {
				// TODO Imbalanced square brackets
				return optionStack, err
			}

		case '+':
			if len(literalBuffer) > 0 {
				last := literalBuffer[len(literalBuffer)-1]
				literalBuffer = literalBuffer[:len(literalBuffer)-1]
				flush()
				expressionList = append(expressionList, regexToken{
					typeOfToken: RepeatOnceOrMore,
					char:        last,
				})
				literalBuffer = nil
			} else if len(expressionList) > 0 {
				last := &expressionList[len(expressionList)-1]
				if last.typeOfToken == LiteralString && len(last.runes) == 1 {
					expressionList[len(expressionList)-1] = regexToken{
						typeOfToken: RepeatOnceOrMore,
						char:        last.runes[0],
					}
				} else if last.typeOfToken == CharacterClass {
					expressionList[len(expressionList)-1] = regexToken{
						typeOfToken: RepeatOnceOrMore,
						body: last.body,
						negated: last.negated,
					}
				} else {
					return nil, fmt.Errorf("unsupported + usage at position %d", i)
				}
			} else {
				return nil, fmt.Errorf("unsupported + at start")
			}

		case '?':
			if len(literalBuffer) > 0 {
				last := literalBuffer[len(literalBuffer)-1]
				literalBuffer = literalBuffer[:len(literalBuffer)-1]
				flush()
				expressionList = append(expressionList, regexToken{
					typeOfToken: Optional,
					char:        last,
				})
				literalBuffer = nil
			} else if len(expressionList) > 0 {
				last := &expressionList[len(expressionList)-1]
				if last.typeOfToken == LiteralString && len(last.runes) == 1 {
					expressionList[len(expressionList)-1] = regexToken{
						typeOfToken: Optional,
						char:        last.runes[0],
					}
				} else {
					return nil, fmt.Errorf("unsupported ? usage at position %d", i)
				}
			} else {
				return nil, fmt.Errorf("unsupported ? at start")
			}

		case '*':
			if escaped {
				literalBuffer = append(literalBuffer, ch)
				escaped = false
			} else {
				if len(literalBuffer) > 0 {
					last := literalBuffer[len(literalBuffer)-1]
					literalBuffer = literalBuffer[:len(literalBuffer)-1]
					flush()
					expressionList = append(expressionList, regexToken{
						typeOfToken: RepeatZeroOrMore,
						char:        last})
				} else if len(expressionList) > 0 {
					last := &expressionList[len(expressionList)-1]
					if last.typeOfToken == LiteralString && len(last.runes) == 1 {
						expressionList[len(expressionList)-1] = regexToken{
							typeOfToken: RepeatZeroOrMore,
							char:        last.runes[0],
						}
					} else if last.typeOfToken == CharacterClass {
						expressionList[len(expressionList)-1] = regexToken{
							typeOfToken: RepeatZeroOrMore,
							body: last.body,
							negated: last.negated,
						}
					} else {
						return nil, fmt.Errorf("* qualifier not supported at position %d", i)
					}

				} else {
					return nil, fmt.Errorf("* qualifier not supported at start")
				}
			}
		case '^':
			if escaped {
				literalBuffer = append(literalBuffer, ch)
				escaped = false
			} else {
				if i == 0 {
					expressionList = append(expressionList, regexToken{
						typeOfToken:   Anchor,
						anchoredStart: true,
					})
				} else {
					literalBuffer = append(literalBuffer, ch)
				}
			}
		case '$':
			if escaped {
				literalBuffer = append(literalBuffer, ch)
				escaped = false
			} else {
				if i == len(expression)-1 {
					if len(literalBuffer) > 0 {
						expressionList = append(expressionList, regexToken{
							typeOfToken: LiteralString,
							runes:       literalBuffer,
						})
					}
					expressionList = append(expressionList, regexToken{
						typeOfToken: Anchor,
						anchoredEnd: true,
					})
					// Should be the last character of the expression...might as well return

				} else {
					literalBuffer = append(literalBuffer, ch)
				}
			}
		case '|':
			if escaped {
				literalBuffer = append(literalBuffer, ch)
				escaped = false
			} else {
				if len(literalBuffer) > 0 {
					expressionList = append(expressionList, regexToken{
						typeOfToken: LiteralString,
						runes:       literalBuffer,
					})
					literalBuffer = nil
				}
				optionStack = append(optionStack, expressionList)
				expressionList = regexTokens{}
				continue
			}
		case '(':
			if escaped {
				literalBuffer = append(literalBuffer, ch)
				escaped = false
			} else {
				flush()
				alternatives, offset := parseParen(expression, i)
				expressionList = append(expressionList, regexToken{
					typeOfToken:  Group,
					alternatives: alternatives,
				})
				i += offset
			}

		default:
			// If in doubt bung it in the string to search literally
			literalBuffer = append(literalBuffer, ch)
		}

	}
	if len(literalBuffer) > 0 {
		expressionList = append(expressionList, regexToken{
			typeOfToken: LiteralString,
			runes:       literalBuffer,
		})
	}
	optionStack = append(optionStack, expressionList)
	return optionStack, nil
}
