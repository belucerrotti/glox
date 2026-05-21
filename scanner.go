package main

import (
	"fmt"
	"strconv"
)

type TokenType int

const (
	PLUS TokenType = iota
	MINUS
	DIVIDE

	LEFT_BRACKET
	RIGHT_BRACKET

	LEFT_PAREN
	RIGHT_PAREN

	LEFT_BRACE
	RIGHT_BRACE

	EQUALS
	EQUALS_EQUALS
	NOTEQUAL
	BANG
	LESS_THAN
	LESS_EQUALS
	GREATER_THAN
	GREATER_EQUALS

	COMMA
	DOT
	SEMICOLON
	COLON
	QUESTION
	AMPERSAND
	STAR
	COMMENT

	NUMBER
	STRING
	IDENTIFIER

	EOF

	AND
	CLASS
	ELSE
	FALSE
	FOR
	FUN
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
)

var tokenNames = map[TokenType]string{
	PLUS:   "PLUS",
	MINUS:  "MINUS",
	DIVIDE: "DIVIDE",

	LEFT_BRACKET:  "LEFT_BRACKET",
	RIGHT_BRACKET: "RIGHT_BRACKET",

	LEFT_PAREN:  "LEFT_PAREN",
	RIGHT_PAREN: "RIGHT_PAREN",

	LEFT_BRACE:  "LEFT_BRACE",
	RIGHT_BRACE: "RIGHT_BRACE",

	EQUALS:         "EQUALS",
	EQUALS_EQUALS:  "EQUALS_EQUALS",
	NOTEQUAL:       "NOTEQUAL",
	BANG:           "BANG",
	LESS_THAN:      "LESS_THAN",
	LESS_EQUALS:    "LESS_EQUALS",
	GREATER_THAN:   "GREATER_THAN",
	GREATER_EQUALS: "GREATER_EQUALS",

	COMMA:     "COMMA",
	DOT:       "DOT",
	SEMICOLON: "SEMICOLON",
	COLON:     "COLON",
	QUESTION:  "QUESTION",
	AMPERSAND: "AMPERSAND",
	STAR:      "STAR",
	COMMENT:   "COMMENT",

	NUMBER:     "NUMBER",
	STRING:     "STRING",
	IDENTIFIER: "IDENTIFIER",

	EOF: "EOF",

	AND:    "AND",
	CLASS:  "CLASS",
	ELSE:   "ELSE",
	FALSE:  "FALSE",
	FOR:    "FOR",
	FUN:    "FUN",
	IF:     "IF",
	NIL:    "NIL",
	OR:     "OR",
	PRINT:  "PRINT",
	RETURN: "RETURN",
	SUPER:  "SUPER",
	THIS:   "THIS",
	TRUE:   "TRUE",
	VAR:    "VAR",
	WHILE:  "WHILE",
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type token struct {
	tokenType  TokenType
	name       string
	value      string
	valueFloat float64
	line       int
	callable   *loxFunction // solo tiene valor cuando el token representa una función
}

type scanner struct {
	tokens       []token
	content      []byte
	currentIndex int
}

func createScanner() *scanner {
	return &scanner{
		tokens:       []token{},
		content:      []byte{},
		currentIndex: 0,
	}
}

func (s *scanner) scan(content []byte) ([]token, error) {
	s.content = content
	currentLine := 1
	for i := 0; i < len(content); i++ {
		if content[i] == ' ' || content[i] == '\t' || content[i] == '\r' || content[i] == '\n' {
			if content[i] == '\n' {
				currentLine++
			}
			continue
		}
		s.currentIndex = i
		newIndex, err := s.scanTokens(currentLine)
		if err != nil {
			return nil, err
		}
		i = newIndex
	}
	s.tokens = append(s.tokens, token{tokenType: EOF, name: tokenNames[EOF], value: "", line: currentLine})
	return s.tokens, nil
}

func (s *scanner) scanTokens(currentLine int) (int, error) {
	content := s.content

	var t = token{}
	switch content[s.currentIndex] {
	case '+':
		t = token{tokenType: PLUS, name: tokenNames[PLUS], value: string(content[s.currentIndex])}
	case '-':
		t = token{tokenType: MINUS, name: tokenNames[MINUS], value: string(content[s.currentIndex])}
	case '/':
		// chequeo si no es un comentario
		if s.currentIndex+1 < len(content) && content[s.currentIndex+1] == '/' {
			start := s.currentIndex + 2
			end := start
			for end < len(content) && content[end] != '\n' {
				end++
			}
			t = token{tokenType: COMMENT, name: tokenNames[COMMENT], value: string(content[start:end])}
			s.currentIndex = end - 1
		} else {
			t = token{tokenType: DIVIDE, name: tokenNames[DIVIDE], value: string(content[s.currentIndex])}
		}
	case '?':
		t = token{tokenType: QUESTION, name: tokenNames[QUESTION], value: string(content[s.currentIndex])}
	case '&':
		t = token{tokenType: AMPERSAND, name: tokenNames[AMPERSAND], value: string(content[s.currentIndex])}
	case '*':
		t = token{tokenType: STAR, name: tokenNames[STAR], value: string(content[s.currentIndex])}
	case '(':
		t = token{tokenType: LEFT_PAREN, name: tokenNames[LEFT_PAREN], value: string(content[s.currentIndex])}
	case ')':
		t = token{tokenType: RIGHT_PAREN, name: tokenNames[RIGHT_PAREN], value: string(content[s.currentIndex])}
	case '{':
		t = token{tokenType: LEFT_BRACE, name: tokenNames[LEFT_BRACE], value: string(content[s.currentIndex])}
	case '}':
		t = token{tokenType: RIGHT_BRACE, name: tokenNames[RIGHT_BRACE], value: string(content[s.currentIndex])}
	case '[':
		t = token{tokenType: LEFT_BRACKET, name: tokenNames[LEFT_BRACKET], value: string(content[s.currentIndex])}
	case ']':
		t = token{tokenType: RIGHT_BRACKET, name: tokenNames[RIGHT_BRACKET], value: string(content[s.currentIndex])}
	case ',':
		t = token{tokenType: COMMA, name: tokenNames[COMMA], value: string(content[s.currentIndex])}
	case '.':
		t = token{tokenType: DOT, name: tokenNames[DOT], value: string(content[s.currentIndex])}
	case ';':
		t = token{tokenType: SEMICOLON, name: tokenNames[SEMICOLON], value: string(content[s.currentIndex])}
	case ':':
		t = token{tokenType: COLON, name: tokenNames[COLON], value: string(content[s.currentIndex])}
	case '=':
		if s.currentIndex+1 < len(content) && content[s.currentIndex+1] == '=' {
			t = token{tokenType: EQUALS_EQUALS, name: tokenNames[EQUALS_EQUALS], value: "=="}
			s.currentIndex++
		} else {
			t = token{tokenType: EQUALS, name: tokenNames[EQUALS], value: string(content[s.currentIndex])}
		}
	case '!':
		if s.currentIndex+1 < len(content) && content[s.currentIndex+1] == '=' {
			t = token{tokenType: NOTEQUAL, name: tokenNames[NOTEQUAL], value: "!="}
			s.currentIndex++
		} else {
			t = token{tokenType: BANG, name: tokenNames[BANG], value: string(content[s.currentIndex])}
		}
	case '<':
		if s.currentIndex+1 < len(content) && content[s.currentIndex+1] == '=' {
			t = token{tokenType: LESS_EQUALS, name: tokenNames[LESS_EQUALS], value: "<="}
			s.currentIndex++
		} else {
			t = token{tokenType: LESS_THAN, name: tokenNames[LESS_THAN], value: string(content[s.currentIndex])}
		}
	case '>':
		if s.currentIndex+1 < len(content) && content[s.currentIndex+1] == '=' {
			t = token{tokenType: GREATER_EQUALS, name: tokenNames[GREATER_EQUALS], value: ">="}
			s.currentIndex++
		} else {
			t = token{tokenType: GREATER_THAN, name: tokenNames[GREATER_THAN], value: string(content[s.currentIndex])}
		}
	case '"':
		start := s.currentIndex + 1
		end := start
		for end < len(content) && content[end] != '"' {
			end++
		}
		if end >= len(content) {
			return s.currentIndex, fmt.Errorf("string sin cerrar: falta '\"' de cierre")
		}
		t = token{tokenType: STRING, name: tokenNames[STRING], value: string(content[start:end])}
		s.currentIndex = end
	default:
		if isDigit(content[s.currentIndex]) {
			err := scanDigit(content, &s.currentIndex, &t)
			if err != nil {
				return s.currentIndex, err
			}
		} else if isAlpha(content[s.currentIndex]) {
			scanAlpha(content, &s.currentIndex, &t)
		} else {
			return s.currentIndex, fmt.Errorf("token/caracter no reconocido: '%c'", content[s.currentIndex])
		}
	}

	t.line = currentLine
	s.tokens = append(s.tokens, t)
	return s.currentIndex, nil
}

func scanDigit(content []byte, index *int, t *token) error {
	isDecimal := false
	for *index < len(content) && ((content[*index] >= '0' && content[*index] <= '9') || content[*index] == '.') {
		if content[*index] == '.' {
			if isDecimal {
				return fmt.Errorf("número con más de un punto decimal")
			}
			isDecimal = true
			t.value += "."
			(*index)++
			continue
		}
		t.value += string(content[*index])
		t.valueFloat = t.valueFloat*10 + float64(content[*index]) - float64('0')
		(*index)++
	}
	(*index)--
	f, err := strconv.ParseFloat(t.value, 64)
	if err != nil {
		return fmt.Errorf("número inválido %s", t.value)
	}
	t.valueFloat = f
	t.tokenType = NUMBER
	t.name = tokenNames[NUMBER]
	return nil
}

func scanAlpha(content []byte, index *int, t *token) {
	for *index < len(content) && isAlphaNumeric(content[*index]) {
		t.value += string(content[*index])
		(*index)++
	}
	(*index)--
	t.tokenType = IDENTIFIER
	t.name = tokenNames[IDENTIFIER]
	if tokenType, ok := keywords[t.value]; ok {
		t.tokenType = tokenType
		t.name = tokenNames[tokenType]
	}
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
