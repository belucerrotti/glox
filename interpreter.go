package main

type interpreter struct {
	statements  []Stmt
	current     int
	environment *environment
}

type environment struct {
	variables map[string]interface{}
	father    *environment
}

func createEnvironment(father *environment) *environment {
	return &environment{
		variables: map[string]interface{}{},
		father:    father,
	}
}

func createInterpreter(statements []Stmt) *interpreter {
	return &interpreter{
		statements:  statements,
		current:     0,
		environment: createEnvironment(nil),
	}
}

func (i *interpreter) isAtEnd() bool {
	return i.current >= len(i.statements)
}

// ejecuta los statements
func (i *interpreter) execute(statement Stmt) error {

	switch s := statement.(type) {
	case *PrintStmt:
		err := i.executePrintStmt(s.expr)
		if err != nil {
			return err
		}
	}

	return nil
}

// evalua las expressions
func (i *interpreter) evaluate(expression Expr) error {
	// todo
	return nil
}

func (i *interpreter) executePrintStmt(expression Expr) error {
	// todo
	return nil
}

func (i *interpreter) interpret() error {

	for !i.isAtEnd() {
		err := i.execute(i.statements[i.current])
		if err != nil {
			return err
		}
		i.current++
	}

	return nil
}
