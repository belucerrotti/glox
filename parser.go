package main

import "fmt"

type parser struct {
	tokens  []token
	current int
}

func (p *parser) currentToken() token {
	return p.tokens[p.current]
}

func (p *parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.currentToken().tokenType == EOF
}

func createParser(t []token) *parser {
	return &parser{
		tokens:  t,
		current: 0,
	}
}

func (p *parser) parse() ([]Stmt, error) {
	p.current = 0
	var statements []Stmt

	for !p.isAtEnd() {
		stmt, err := p.getStatement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

func (p *parser) matches(tokenType TokenType) bool {
	if !p.isAtEnd() && p.tokens[p.current].tokenType == tokenType {
		return true
	}
	return false
}

func (p *parser) getStatement() (Stmt, error) {
	if p.matches(PRINT) {
		return p.printStatement()
	}
	if p.matches(VAR) {
		return p.varDeclaration()
	}
	if p.matches(FUN) {
		return p.funDeclaration()
	}
	if p.matches(RETURN) {
		return p.returnStatement()
	}
	// if p.matches(IF) {
	// 	return p.ifStatement()
	// }
	// if p.matches(WHILE) {
	// 	return p.whileStatement()
	// }
	// if p.matches(FOR) {
	// 	return p.forStatement()
	// }
	// if p.matches(LEFT_BRACE) {
	// 	return p.blockStatement()
	// }

	return p.expressionStatement()
}

func (p *parser) returnStatement() (Stmt, error) {
	if p.isAtEnd() {
		return nil, fmt.Errorf("line %d: expected return value but found '%s'", p.currentToken().line, p.currentToken().value)
	}
	return p.expressionStatement()
}

func (p *parser) funDeclaration() (Stmt, error) {
	println("entré a fun")
	funName := p.currentToken()
	if !p.matches(IDENTIFIER) {
		return nil, fmt.Errorf("line %d: expected IDENTIFIER but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	println("la funName es " + funName.value)

	if !p.matches(LEFT_PAREN) {
		return nil, fmt.Errorf("line %d: expected '(' but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	parameters := []token{}

	for !p.isAtEnd() && p.matches(IDENTIFIER) {
		parameters = append(parameters, p.tokens[p.current-1])

		if p.matches(COMMA) {
			if !p.matches(IDENTIFIER) {
				return nil, fmt.Errorf("line %d: expected IDENTIFIER but found '%s'", p.currentToken().line, p.currentToken().value)
			}
			continue
		}
	}

	println("parameters:")
	for i := 0; i < len(parameters); i++ {
		println(parameters[i].value)
	}

	if !p.matches(RIGHT_PAREN) {
		return nil, fmt.Errorf("line %d: expected ')' but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	body := []Stmt{}

	if !p.matches(LEFT_BRACE) {
		return nil, fmt.Errorf("line %d: expected '{' but found '%s'", p.currentToken().line, p.currentToken().value)
	} else {
		for !p.isAtEnd() && !p.matches(RIGHT_BRACE) {
			stmt, err := p.getStatement()
			if err != nil {
				return nil, err
			}
			body = append(body, stmt)
		}

		if !p.matches(RIGHT_BRACE) && p.isAtEnd() {
			return nil, fmt.Errorf("line %d: expected '}' but found '%s'", p.currentToken().line, p.currentToken().value)
		}

	}

	return &FunDecl{name: funName, parameters: parameters, body: body}, nil
}

func (p *parser) varDeclaration() (Stmt, error) {
	p.current++ // consumo el 'var'

	name := p.currentToken()
	if !p.matches(IDENTIFIER) {
		return nil, fmt.Errorf("line %d: expected IDENTIFIER but found '%s'", p.currentToken().line, p.currentToken().value)
	} else {
		if p.matches(EQUALS) {
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			return &VarDecl{name: name, value: expr}, nil
		}
	}

	if !p.matches(SEMICOLON) {
		return nil, fmt.Errorf("line %d: expected ';' but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	return &VarDecl{name: name}, nil
}

func (p *parser) printStatement() (Stmt, error) {
	p.current++ // consumo el 'print'
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if !p.matches(SEMICOLON) {
		return nil, fmt.Errorf("line %d: expected ';' but found '%s'", p.currentToken().line, p.currentToken().value)
	} else {
	}

	return &PrintStmt{expr: expr}, nil
}

func (p *parser) expressionStatement() (Stmt, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if !p.matches(SEMICOLON) {
		return nil, fmt.Errorf("line %d: expected ';' but found '%s'", p.currentToken().line, p.currentToken().value)
	} else {
	}

	return &ExpressionStmt{expr: expr}, nil
}

// expression → equality
func (p *parser) parseExpression() (Expr, error) {
	return p.parseEquality()
}

// equality → comparison ( ( "!=" | "==" ) comparison )*
func (p *parser) parseEquality() (Expr, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	for p.matches(EQUALS_EQUALS) || p.matches(NOTEQUAL) {
		operator := p.tokens[p.current]
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{left: left, operator: operator, right: right}
	}
	return left, nil
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *parser) parseComparison() (Expr, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	if p.matches(GREATER_THAN) || p.matches(GREATER_EQUALS) || p.matches(LESS_THAN) || p.matches(LESS_EQUALS) {
		operator := p.tokens[p.current]
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		return &BinaryExpr{left: left, operator: operator, right: right}, nil
	}

	return left, nil
}

// term → factor ( ( "+" | "-" ) factor )*
func (p *parser) parseTerm() (Expr, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	if p.matches(PLUS) || p.matches(MINUS) {
		operator := p.tokens[p.current]
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		return &BinaryExpr{left: left, operator: operator, right: right}, nil
	}

	return left, nil
}

// factor → unary ( ( "*" | "/" ) unary )*
func (p *parser) parseFactor() (Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	if p.matches(STAR) || p.matches(DIVIDE) {
		operator := p.tokens[p.current]
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &BinaryExpr{left: left, operator: operator, right: right}, nil
	}

	return left, nil
}

// unary → ( "!" | "-" ) unary | primary
func (p *parser) parseUnary() (Expr, error) {
	if p.matches(BANG) || p.matches(MINUS) {
		operator := p.tokens[p.current]
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{operator: operator, right: right}, nil
	}

	return p.parsePrimary()
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
func (p *parser) parsePrimary() (Expr, error) {
	if p.matches(NUMBER) || p.matches(STRING) || p.matches(TRUE) || p.matches(FALSE) || p.matches(NIL) {
		return &LiteralExpr{value: p.tokens[p.current-1]}, nil
	}

	// todo: grouping "(expression)"
	return nil, nil
}
