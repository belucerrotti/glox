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
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	return statements, nil
}

func (p *parser) matches(tokenType TokenType) bool {
	if !p.isAtEnd() && p.tokens[p.current].tokenType == tokenType {
		p.current++
		return true
	}
	return false
}

func (p *parser) getStatement() (Stmt, error) {
	if p.matches(COMMENT) {
		return nil, nil
	}
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
	if p.matches(IF) {
		return p.ifStatement()
	}
	if p.matches(WHILE) {
		return p.whileStatement()
	}
	if p.matches(FOR) {
		return p.forStatement()
	}
	if p.matches(LEFT_BRACE) {
		return p.blockStatement()
	}

	return p.expressionStatement()
}

func (p *parser) blockStatement() (Stmt, error) {
	statements := []Stmt{}
	for !p.matches(RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.getStatement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	if p.isAtEnd() && p.tokens[p.current-1].tokenType != RIGHT_BRACE {
		return nil, fmt.Errorf("line %d: expected '}'", p.currentToken().line)
	}
	return &BlockStmt{statements: statements}, nil
}

func (p *parser) forStatement() (Stmt, error) {
	if !p.matches(LEFT_PAREN) {
		return nil, fmt.Errorf("line %d: expected '(' after 'for'", p.currentToken().line)
	}

	initializer, err := p.getStatement()
	if err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if !p.matches(SEMICOLON) {
		return nil, fmt.Errorf("line %d: expected ';' after 'for' condition", p.currentToken().line)
	}

	increment, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if !p.matches(RIGHT_PAREN) {
		return nil, fmt.Errorf("line %d: expected ')' after 'for' increment", p.currentToken().line)
	}

	body, err := p.getStatement()
	if err != nil {
		return nil, err
	}

	return &ForStmt{initializer: initializer, condition: condition, increment: increment, body: body}, nil
}

func (p *parser) whileStatement() (Stmt, error) {
	if !p.matches(LEFT_PAREN) {
		return nil, fmt.Errorf("line %d: expected '(' after 'while'", p.currentToken().line)
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if !p.matches(RIGHT_PAREN) {
		return nil, fmt.Errorf("line %d: expected ')' but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	body, err := p.getStatement()
	if err != nil {
		return nil, err
	}

	return &WhileStmt{condition: condition, body: body}, nil
}

func (p *parser) ifStatement() (Stmt, error) {
	if !p.matches(LEFT_PAREN) {
		return nil, fmt.Errorf("line %d: expected '(' but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if !p.matches(RIGHT_PAREN) {
		return nil, fmt.Errorf("line %d: expected ')' but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	thenBranch, err := p.getStatement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt
	if p.matches(ELSE) {
		eb, err := p.getStatement()
		if err != nil {
			return nil, err
		}
		elseBranch = eb
	}

	return &IfStmt{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}, nil

}

func (p *parser) returnStatement() (Stmt, error) {
	if p.matches(SEMICOLON) {
		return &ReturnStmt{}, nil
	}

	expr, err := p.parseExpression()

	if err != nil {
		return nil, err
	}

	if !p.matches(SEMICOLON) {
		return nil, fmt.Errorf("line %d: expected ';' but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	return &ReturnStmt{value: expr}, nil
}

func (p *parser) funDeclaration() (Stmt, error) {
	funName := p.currentToken()

	if !p.matches(IDENTIFIER) {
		return nil, fmt.Errorf("line %d: expected IDENTIFIER but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	if !p.matches(LEFT_PAREN) {
		return nil, fmt.Errorf("line %d: expected '(' but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	parameters := []token{}

	for !p.isAtEnd() && p.matches(IDENTIFIER) {
		parameters = append(parameters, p.tokens[p.current-1])
		if !p.matches(COMMA) {
			break
		}
	}

	if !p.matches(RIGHT_PAREN) {
		return nil, fmt.Errorf("line %d: expected ')' but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	if !p.matches(LEFT_BRACE) {
		return nil, fmt.Errorf("line %d: expected '{' but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	body := []Stmt{}
	for !p.isAtEnd() && !p.matches(RIGHT_BRACE) {
		stmt, err := p.getStatement()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}

	return &FunDecl{name: funName, parameters: parameters, body: body}, nil
}

func (p *parser) varDeclaration() (Stmt, error) {
	name := p.currentToken()

	if !p.matches(IDENTIFIER) {
		return nil, fmt.Errorf("line %d: expected IDENTIFIER but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	var value Expr
	if p.matches(EQUALS) {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		value = expr
	}

	if !p.matches(SEMICOLON) {
		return nil, fmt.Errorf("line %d: expected ';' but found '%s'", p.currentToken().line, p.currentToken().value)
	}

	return &VarDecl{name: name, value: value}, nil
}

func (p *parser) printStatement() (Stmt, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if !p.matches(SEMICOLON) {
		return nil, fmt.Errorf("line %d: expected ';' but found '%s'", p.currentToken().line, p.currentToken().value)
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
	}

	return &ExpressionStmt{expr: expr}, nil
}

// expression → assignment
func (p *parser) parseExpression() (Expr, error) {
	return p.parseAssignment()
}

// assignment → IDENTIFIER "=" assignment | equality
func (p *parser) parseAssignment() (Expr, error) {
	expr, err := p.parseLogicalOr()
	if err != nil {
		return nil, err
	}

	if p.matches(EQUALS) {
		value, err := p.parseAssignment()
		if err != nil {
			return nil, err
		}

		if variable, ok := expr.(*VariableExpr); ok {
			return &AssignExpr{name: variable.name, value: value}, nil
		}

		return nil, fmt.Errorf("line %d: invalid assignment target", p.currentToken().line)
	}

	return expr, nil
}

// logic_or → logic_and ( "or" logic_and )*
func (p *parser) parseLogicalOr() (Expr, error) {
	left, err := p.parseLogicalAnd()
	if err != nil {
		return nil, err
	}
	for p.matches(OR) {
		operator := p.tokens[p.current-1]
		right, err := p.parseLogicalAnd()
		if err != nil {
			return nil, err
		}
		left = &LogicalExpr{left: left, operator: operator, right: right}
	}
	return left, nil
}

// logic_and → equality ( "and" equality )*
func (p *parser) parseLogicalAnd() (Expr, error) {
	left, err := p.parseEquality()
	if err != nil {
		return nil, err
	}
	for p.matches(AND) {
		operator := p.tokens[p.current-1]
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		left = &LogicalExpr{left: left, operator: operator, right: right}
	}
	return left, nil
}

// equality → comparison ( ( "!=" | "==" ) comparison )*
func (p *parser) parseEquality() (Expr, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.matches(EQUALS_EQUALS) || p.matches(NOTEQUAL) {
		operator := p.tokens[p.current-1]
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

	for p.matches(GREATER_THAN) || p.matches(GREATER_EQUALS) || p.matches(LESS_THAN) || p.matches(LESS_EQUALS) {
		operator := p.tokens[p.current-1]
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{left: left, operator: operator, right: right}
	}

	return left, nil
}

// term → factor ( ( "+" | "-" ) factor )*
func (p *parser) parseTerm() (Expr, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.matches(PLUS) || p.matches(MINUS) {
		operator := p.tokens[p.current-1]
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{left: left, operator: operator, right: right}
	}

	return left, nil
}

// factor → unary ( ( "*" | "/" ) unary )*
func (p *parser) parseFactor() (Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.matches(STAR) || p.matches(DIVIDE) {
		operator := p.tokens[p.current-1]
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{left: left, operator: operator, right: right}
	}

	return left, nil
}

// unary → ( "!" | "-" ) unary | primary
func (p *parser) parseUnary() (Expr, error) {
	if p.matches(BANG) || p.matches(MINUS) {
		operator := p.tokens[p.current-1]
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{operator: operator, right: right}, nil
	}

	return p.parseCall()
}

// call → primary ( "(" arguments? ")" )*
func (p *parser) parseCall() (Expr, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for p.matches(LEFT_PAREN) {
		paren := p.tokens[p.current-1]
		arguments := []Expr{}

		for !p.isAtEnd() && !p.matches(RIGHT_PAREN) {
			arg, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
			if !p.matches(COMMA) {
				if !p.matches(RIGHT_PAREN) {
					return nil, fmt.Errorf("line %d: expected ')' but found '%s'", p.currentToken().line, p.currentToken().value)
				}
				break
			}
		}

		expr = &CallExpr{callee: expr, paren: paren, arguments: arguments}
	}

	return expr, nil
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
func (p *parser) parsePrimary() (Expr, error) {
	if p.matches(NUMBER) || p.matches(STRING) || p.matches(TRUE) || p.matches(FALSE) || p.matches(NIL) {
		return &LiteralExpr{value: p.tokens[p.current-1]}, nil
	}

	if p.matches(IDENTIFIER) {
		return &VariableExpr{name: p.tokens[p.current-1]}, nil
	}

	if p.matches(LEFT_PAREN) {
		e, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if !p.matches(RIGHT_PAREN) {
			return nil, fmt.Errorf("line %d: expected ')' but found '%s'", p.currentToken().line, p.currentToken().value)
		}

		return &GroupingExpr{expression: e}, nil
	}

	return nil, fmt.Errorf("line %d: expected expression but found '%s'", p.currentToken().line, p.currentToken().value)
}
