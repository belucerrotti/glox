package main

import "fmt"

type TokenType int

const (
	PLUS TokenType = iota
	MINUS
	DIVIDE
	PERCENT

	LEFT_BRACKET
	RIGHT_BRACKET

	LEFT_PAREN
	RIGHT_PAREN

	LEFT_BRACE
	RIGHT_BRACE

	EQUALS
	NOTEQUAL
	BANG
	LESS_THAN
	GREATER_THAN
	MOD

	COMMA
	DOT
	SEMICOLON
	COLON
	QUESTION
	AMPERSAND
	STAR

	NUMBER
	STRING
)

var tokenNames = map[TokenType]string{
	PLUS:    "PLUS",
	MINUS:   "MINUS",
	DIVIDE:  "DIVIDE",
	PERCENT: "PERCENT",

	LEFT_BRACKET:  "LEFT_BRACKET",
	RIGHT_BRACKET: "RIGHT_BRACKET",

	LEFT_PAREN:  "LEFT_PAREN",
	RIGHT_PAREN: "RIGHT_PAREN",

	LEFT_BRACE:  "LEFT_BRACE",
	RIGHT_BRACE: "RIGHT_BRACE",

	EQUALS:       "EQUALS",
	NOTEQUAL:     "NOTEQUAL",
	BANG:         "BANG",
	LESS_THAN:    "LESS_THAN",
	GREATER_THAN: "GREATER_THAN",
	MOD:          "MOD",

	COMMA:     "COMMA",
	DOT:       "DOT",
	SEMICOLON: "SEMICOLON",
	COLON:     "COLON",
	QUESTION:  "QUESTION",
	AMPERSAND: "AMPERSAND",
	STAR:      "STAR",

	NUMBER: "NUMBER",
	STRING: "STRING",
}

type token struct {
	tokenType TokenType
	name      string
	value     string
	valueInt  int
}

type scanner struct {
	tokens []token
}

func createScanner() *scanner {
	return &scanner{
		tokens: []token{},
	}
}

func (s *scanner) scan(content []byte) ([]token, error) {
	var tokens = []token{}
	for i := 0; i < len(content); i++ {
		if content[i] != ' ' {
			t, newIndex, err := s.scanToken(content, i)
			if err != nil {
				return nil, err
			}
			i = newIndex
			tokens = append(tokens, t)
		}
	}
	return tokens, nil
}

func (s *scanner) scanToken(content []byte, index int) (token, int, error) {
	var t = token{}
	switch content[index] {
	case '+':
		t = token{tokenType: PLUS, name: tokenNames[PLUS], value: string(content[index])}
	case '-':
		t = token{tokenType: MINUS, name: tokenNames[MINUS], value: string(content[index])}
	case '/':
		t = token{tokenType: DIVIDE, name: tokenNames[DIVIDE], value: string(content[index])}
	case '%':
		t = token{tokenType: PERCENT, name: tokenNames[PERCENT], value: string(content[index])}
	case '?':
		t = token{tokenType: QUESTION, name: tokenNames[QUESTION], value: string(content[index])}
	case '&':
		t = token{tokenType: AMPERSAND, name: tokenNames[AMPERSAND], value: string(content[index])}
	case '*':
		t = token{tokenType: STAR, name: tokenNames[STAR], value: string(content[index])}
	case '(':
		t = token{tokenType: LEFT_PAREN, name: tokenNames[LEFT_PAREN], value: string(content[index])}
	case ')':
		t = token{tokenType: RIGHT_PAREN, name: tokenNames[RIGHT_PAREN], value: string(content[index])}
	case '{':
		t = token{tokenType: LEFT_BRACE, name: tokenNames[LEFT_BRACE], value: string(content[index])}
	case '}':
		t = token{tokenType: RIGHT_BRACE, name: tokenNames[RIGHT_BRACE], value: string(content[index])}
	case '[':
		t = token{tokenType: LEFT_BRACKET, name: tokenNames[LEFT_BRACKET], value: string(content[index])}
	case ']':
		t = token{tokenType: RIGHT_BRACKET, name: tokenNames[RIGHT_BRACKET], value: string(content[index])}
	case '=':
		t = token{tokenType: EQUALS, name: tokenNames[EQUALS], value: string(content[index])}
	case '!':
		t = token{tokenType: BANG, name: tokenNames[BANG], value: string(content[index])}
	case '<':
		t = token{tokenType: LESS_THAN, name: tokenNames[LESS_THAN], value: string(content[index])}
	case '>':
		t = token{tokenType: GREATER_THAN, name: tokenNames[GREATER_THAN], value: string(content[index])}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		t = token{tokenType: NUMBER, name: tokenNames[NUMBER], value: string(content[index]), valueInt: int(content[index]) - int('0')}
	case ',':
		t = token{tokenType: COMMA, name: tokenNames[COMMA], value: string(content[index])}
	case '.':
		t = token{tokenType: DOT, name: tokenNames[DOT], value: string(content[index])}
	case ';':
		t = token{tokenType: SEMICOLON, name: tokenNames[SEMICOLON], value: string(content[index])}
	case ':':
		t = token{tokenType: COLON, name: tokenNames[COLON], value: string(content[index])}
	case '"':
		start := index + 1
		end := start
		for end < len(content) && content[end] != '"' {
			end++
		}
		if end >= len(content) {
			return token{}, index, fmt.Errorf("string sin cerrar: falta '\"' de cierre")
		}
		t = token{tokenType: STRING, name: tokenNames[STRING], value: string(content[start:end])}
		index = end

	default:
		t = token{value: string(content[index]), valueInt: int(content[index]) - int('0')}
	}

	return t, index, nil
}
