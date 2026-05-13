package main

type interpreter struct {
	statements []Stmt
	current    int
}

func createInterpreter(statements []Stmt) *interpreter {
	return &interpreter{
		statements: statements,
		current:    0,
	}
}

func (i *interpreter) interpret() {
	// todo
}
