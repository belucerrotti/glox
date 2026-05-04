package main

type parser struct {
	tokens []token
}

func createParser() *parser {
	return &parser{
		tokens: []token{},
	}
}

func (p *parser) parse(tokens []token) {
	p.tokens = tokens
}
