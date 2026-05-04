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
	EQUALS_EQUALS
	NOTEQUAL
	BANG
	BANG_EQUALS
	LESS_THAN
	LESS_EQUALS
	GREATER_THAN
	GREATER_EQUALS
	MOD

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

	EQUALS:         "EQUALS",
	EQUALS_EQUALS:  "EQUALS_EQUALS",
	NOTEQUAL:       "NOTEQUAL",
	BANG:           "BANG",
	BANG_EQUALS:    "BANG_EQUALS",
	LESS_THAN:      "LESS_THAN",
	LESS_EQUALS:    "LESS_EQUALS",
	GREATER_THAN:   "GREATER_THAN",
	GREATER_EQUALS: "GREATER_EQUALS",
	MOD:            "MOD",

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
		if content[i] == ' ' || content[i] == '\t' || content[i] == '\r' || content[i] == '\n' {
			continue
		}
		t, newIndex, err := s.scanToken(content, i)
		if err != nil {
			return nil, err
		}
		i = newIndex
		tokens = append(tokens, t)
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
		// chequeo si no es un comentario
		if index+1 < len(content) && content[index+1] == '/' {
			start := index + 2
			end := start
			for end < len(content) && content[end] != '\n' {
				end++
			}
			t = token{tokenType: COMMENT, name: tokenNames[COMMENT], value: string(content[start:end])}
			index = end
		} else {
			t = token{tokenType: DIVIDE, name: tokenNames[DIVIDE], value: string(content[index])}
		}
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
	case ',':
		t = token{tokenType: COMMA, name: tokenNames[COMMA], value: string(content[index])}
	case '.':
		t = token{tokenType: DOT, name: tokenNames[DOT], value: string(content[index])}
	case ';':
		t = token{tokenType: SEMICOLON, name: tokenNames[SEMICOLON], value: string(content[index])}
	case ':':
		t = token{tokenType: COLON, name: tokenNames[COLON], value: string(content[index])}
	case '=':
		if index+1 < len(content) && content[index+1] == '=' {
			t = token{tokenType: EQUALS_EQUALS, name: tokenNames[EQUALS_EQUALS], value: "=="}
			index++
		} else {
			t = token{tokenType: EQUALS, name: tokenNames[EQUALS], value: string(content[index])}
		}
	case '!':
		if index+1 < len(content) && content[index+1] == '=' {
			t = token{tokenType: BANG_EQUALS, name: tokenNames[BANG_EQUALS], value: "!="}
			index++
		} else {
			t = token{tokenType: BANG, name: tokenNames[BANG], value: string(content[index])}
		}
	case '<':
		if index+1 < len(content) && content[index+1] == '=' {
			t = token{tokenType: LESS_EQUALS, name: tokenNames[LESS_EQUALS], value: "<="}
			index++
		} else {
			t = token{tokenType: LESS_THAN, name: tokenNames[LESS_THAN], value: string(content[index])}
		}
	case '>':
		if index+1 < len(content) && content[index+1] == '=' {
			t = token{tokenType: GREATER_EQUALS, name: tokenNames[GREATER_EQUALS], value: ">="}
			index++
		} else {
			t = token{tokenType: GREATER_THAN, name: tokenNames[GREATER_THAN], value: string(content[index])}
		}
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
		if isDigit(content[index]) {
			err := scanDigit(content, &index, &t)
			if err != nil {
				return token{}, index, err
			}
		} else if isAlpha(content[index]) {
			scanAlpha(content, &index, &t)
		} else {
			return token{}, index, fmt.Errorf("token/caracter no reconocido: '%c'", content[index])
		}
	}

	return t, index, nil
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
		t.valueInt = t.valueInt*10 + int(content[*index]) - int('0')
		(*index)++
	}
	t.tokenType = NUMBER
	t.name = tokenNames[NUMBER]
	(*index)--
	return nil
}

func scanAlpha(content []byte, index *int, t *token) {
	for *index < len(content) && isAlpha(content[*index]) {
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
