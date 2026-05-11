package main

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

func createParser() *parser {
	return &parser{
		tokens:  []token{},
		current: 0,
	}
}

func (p *parser) parse(tokens []token) ([]Stmt, error) {
	p.tokens = tokens
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

func (p *parser) getStatement() (Stmt, error) {
	var statement Stmt

	return statement, nil

}
