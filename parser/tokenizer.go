package parser

import (
	"strings"
	"unicode"
)

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
}

// TokenType represents the type of token
type TokenType int

const (
	TokenKeyword TokenType = iota
	TokenIdentifier
	TokenNumber
	TokenString
	TokenOperator
	TokenComma
	TokenLeftParen
	TokenRightParen
	TokenEOF
)

// Tokenize breaks an input string into tokens
func Tokenize(input string) []Token {
	var tokens []Token
	input = strings.TrimSpace(input)
	i := 0

	for i < len(input) {
		// Skip whitespace
		if unicode.IsSpace(rune(input[i])) {
			i++
			continue
		}

		// Handle strings (single or double quotes)
		if input[i] == '\'' || input[i] == '"' {
			quote := input[i]
			i++
			start := i
			for i < len(input) && input[i] != quote {
				i++
			}
			tokens = append(tokens, Token{
				Type:  TokenString,
				Value: input[start:i],
			})
			i++ // Skip closing quote
			continue
		}

		// Handle operators and special characters
		if input[i] == '=' || input[i] == '!' || input[i] == '>' || input[i] == '<' {
			start := i
			i++
			// Handle != >= <=
			if i < len(input) && input[i] == '=' {
				i++
			}
			tokens = append(tokens, Token{
				Type:  TokenOperator,
				Value: input[start:i],
			})
			continue
		}

		if input[i] == ',' {
			tokens = append(tokens, Token{Type: TokenComma, Value: ","})
			i++
			continue
		}

		if input[i] == '*' {
			tokens = append(tokens, Token{Type: TokenIdentifier, Value: "*"})
			i++
			continue
		}

		if input[i] == '(' {
			tokens = append(tokens, Token{Type: TokenLeftParen, Value: "("})
			i++
			continue
		}

		if input[i] == ')' {
			tokens = append(tokens, Token{Type: TokenRightParen, Value: ")"})
			i++
			continue
		}

		// Handle numbers
		if unicode.IsDigit(rune(input[i])) {
			start := i
			for i < len(input) && unicode.IsDigit(rune(input[i])) {
				i++
			}
			tokens = append(tokens, Token{
				Type:  TokenNumber,
				Value: input[start:i],
			})
			continue
		}

		// Handle identifiers and keywords
		if unicode.IsLetter(rune(input[i])) || input[i] == '_' {
			start := i
			for i < len(input) && (unicode.IsLetter(rune(input[i])) || unicode.IsDigit(rune(input[i])) || input[i] == '_' || input[i] == '.') {
				i++
			}
			value := input[start:i]
			tokenType := TokenIdentifier

			// Check if it's a keyword
			upperValue := strings.ToUpper(value)
			if isKeyword(upperValue) {
				tokenType = TokenKeyword
			}

			tokens = append(tokens, Token{
				Type:  tokenType,
				Value: value,
			})
			continue
		}

		// Unknown character, skip it
		i++
	}

	tokens = append(tokens, Token{Type: TokenEOF, Value: ""})
	return tokens
}

// isKeyword checks if a string is a SQL keyword
func isKeyword(s string) bool {
	keywords := map[string]bool{
		"CREATE": true, "TABLE": true, "INSERT": true, "INTO": true,
		"VALUES": true, "SELECT": true, "FROM": true, "WHERE": true,
		"UPDATE": true, "SET": true, "DELETE": true, "INNER": true,
		"JOIN": true, "ON": true, "AND": true, "OR": true,
		"PRIMARY": true, "KEY": true, "UNIQUE": true, "NOT": true,
		"NULL": true, "INT": true, "STRING": true, "BOOL": true,
	}
	return keywords[s]
}
