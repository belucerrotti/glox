package main

type parser struct {
	tokens []token
}

func createParser() *parser {
	return &parser{
		tokens: []token{},
	}
}

func (p *parser) parse(tokens []token) ([]Expr, error) {
	p.tokens = tokens

	if len(p.tokens) == 0 {
		return nil, nil
	}

	var expressions []Expr

	return expressions, nil
}
