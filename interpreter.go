package main

import (
	"fmt"
)

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
	var err error

	switch s := statement.(type) {
	case *PrintStmt:
		err = i.executePrintStmt(s.expr)

	}

	if err != nil {
		return err
	}

	return nil
}

// evalua las expressions
func (i *interpreter) evaluate(expression Expr) (token, error) {
	var err error
	var token token

	switch e := expression.(type) {
	case *LiteralExpr:
		token = e.value
	case *GroupingExpr:
		token, err = i.evaluate(e.expression)
	case *BinaryExpr:
		token, err = i.evaluateBinaryExpression(e.left, e.operator, e.right)
	case *UnaryExpr:
		token, err = i.evaluateUnaryExpression(e.operator, e.right)
	}

	if err != nil {
		return token, err
	}
	return token, nil
}

func (i *interpreter) evaluateUnaryExpression(op token, right Expr) (token, error) {
	value, err := i.evaluate(right)
	if err != nil {
		return token{}, err
	}

	switch op.tokenType {
	case MINUS:
		if value.tokenType != NUMBER {
			return token{}, fmt.Errorf("line %d: '-' operator can only be applied to numbers", op.line)
		}
		result := -value.valueFloat
		return token{tokenType: NUMBER, valueFloat: result, value: fmt.Sprintf("%g", result), line: op.line}, nil
	case BANG:
		// en Lox, false y nil son falsy, todo lo demás es truthy
		isTruthy := value.tokenType != FALSE && value.tokenType != NIL
		if !isTruthy {
			return token{tokenType: TRUE, value: "true", line: op.line}, nil
		}
		return token{tokenType: FALSE, value: "false", line: op.line}, nil
	}

	return token{}, fmt.Errorf("line %d: unknown unary operator '%s'", op.line, op.value)
}

func (i *interpreter) evaluateBinaryExpression(left Expr, op token, right Expr) (token, error) {
	leftFinalValue, err := i.evaluate(left)
	if err != nil {
		return token{}, err
	}

	rightFinalValue, err := i.evaluate(right)
	if err != nil {
		return token{}, err
	}

	sameType := true

	if leftFinalValue.tokenType != rightFinalValue.tokenType {
		sameType = false
	}

	numericOperation := sameType && leftFinalValue.tokenType == NUMBER
	stringOperation := sameType && leftFinalValue.tokenType == STRING

	switch op.tokenType {
	case PLUS:
		if numericOperation {
			result := leftFinalValue.valueFloat + rightFinalValue.valueFloat
			return token{tokenType: NUMBER, valueFloat: result, value: fmt.Sprintf("%g", result), line: op.line}, nil
		}
		if stringOperation {
			result := leftFinalValue.value + rightFinalValue.value
			return token{tokenType: STRING, value: result, line: op.line}, nil
		}
		return token{}, fmt.Errorf("line %d, cannot apply '+' to %s and %s, only between numbers or strings", op.line, leftFinalValue.name, rightFinalValue.name)
	case MINUS:
		if numericOperation {
			result := leftFinalValue.valueFloat - rightFinalValue.valueFloat
			return token{tokenType: NUMBER, valueFloat: result, value: fmt.Sprintf("%g", result), line: op.line}, nil
		}
		if stringOperation {
			return token{}, fmt.Errorf("line %d: cannot substract strings", leftFinalValue.line)
		}
		return token{}, fmt.Errorf("line %d, cannot apply '-' to %s and %s, only between numbers", op.line, leftFinalValue.name, rightFinalValue.name)
	case STAR:
		if numericOperation {
			result := leftFinalValue.valueFloat * rightFinalValue.valueFloat
			return token{tokenType: NUMBER, valueFloat: result, value: fmt.Sprintf("%g", result), line: op.line}, nil
		}
		if stringOperation {
			return token{}, fmt.Errorf("line %d: cannot multiply strings", leftFinalValue.line)
		}
		return token{}, fmt.Errorf("line %d, cannot apply '*' to %s and %s, only between numbers", op.line, leftFinalValue.name, rightFinalValue.name)
	case DIVIDE:
		if numericOperation {
			if rightFinalValue.valueFloat == 0 {
				return token{}, fmt.Errorf("line %d: division by zero", op.line)
			}
			result := leftFinalValue.valueFloat / rightFinalValue.valueFloat
			return token{tokenType: NUMBER, valueFloat: result, value: fmt.Sprintf("%g", result), line: op.line}, nil
		}
		if stringOperation {
			return token{}, fmt.Errorf("line %d: cannot divide strings", leftFinalValue.line)
		}
		return token{}, fmt.Errorf("line %d, cannot apply '/' to %s and %s, only between numbers", op.line, leftFinalValue.name, rightFinalValue.name)
	case EQUALS_EQUALS:
		var result bool
		if !sameType {
			result = false
		} else if numericOperation {
			result = leftFinalValue.valueFloat == rightFinalValue.valueFloat
		} else {
			// strings, booleans, nil: comparar por tokenType y value
			result = leftFinalValue.tokenType == rightFinalValue.tokenType && leftFinalValue.value == rightFinalValue.value
		}
		if result {
			return token{tokenType: TRUE, value: "true", line: op.line}, nil
		} else {
			return token{tokenType: FALSE, value: "false", line: op.line}, nil
		}
	case NOTEQUAL:
		var result bool
		if !sameType {
			result = true
		} else if numericOperation {
			result = leftFinalValue.valueFloat != rightFinalValue.valueFloat
		} else {
			// strings, booleans, nil: comparar por tokenType y value
			result = leftFinalValue.tokenType != rightFinalValue.tokenType || leftFinalValue.value != rightFinalValue.value
		}
		if result {
			return token{tokenType: TRUE, value: "true", line: op.line}, nil
		} else {
			return token{tokenType: FALSE, value: "false", line: op.line}, nil
		}
	case GREATER_THAN:
		if numericOperation {
			result := leftFinalValue.valueFloat > rightFinalValue.valueFloat
			if result {
				return token{tokenType: TRUE, value: "true", line: op.line}, nil
			} else {
				return token{tokenType: FALSE, value: "false", line: op.line}, nil
			}
		}
		if stringOperation {
			return token{}, fmt.Errorf("line %d: cannot compare strings with '>' operator", leftFinalValue.line)
		}
		return token{}, fmt.Errorf("line %d, cannot apply '>' to %s and %s, only between numbers", op.line, leftFinalValue.name, rightFinalValue.name)
	case GREATER_EQUALS:
		if numericOperation {
			result := leftFinalValue.valueFloat >= rightFinalValue.valueFloat
			if result {
				return token{tokenType: TRUE, value: "true", line: op.line}, nil
			} else {
				return token{tokenType: FALSE, value: "false", line: op.line}, nil
			}
		}
		if stringOperation {
			return token{}, fmt.Errorf("line %d: cannot compare strings with '>=' operator", leftFinalValue.line)
		}
		return token{}, fmt.Errorf("line %d, cannot apply '>=' to %s and %s, only between numbers", op.line, leftFinalValue.name, rightFinalValue.name)
	case LESS_THAN:
		if numericOperation {
			result := leftFinalValue.valueFloat < rightFinalValue.valueFloat
			if result {
				return token{tokenType: TRUE, value: "true", line: op.line}, nil
			} else {
				return token{tokenType: FALSE, value: "false", line: op.line}, nil
			}
		}
		if stringOperation {
			return token{}, fmt.Errorf("line %d: cannot compare strings with '<' operator", leftFinalValue.line)
		}
		return token{}, fmt.Errorf("line %d, cannot apply '<' to %s and %s, only between numbers", op.line, leftFinalValue.name, rightFinalValue.name)
	case LESS_EQUALS:
		if numericOperation {
			result := leftFinalValue.valueFloat <= rightFinalValue.valueFloat
			if result {
				return token{tokenType: TRUE, value: "true", line: op.line}, nil
			} else {
				return token{tokenType: FALSE, value: "false", line: op.line}, nil
			}
		}
		if stringOperation {
			return token{}, fmt.Errorf("line %d: cannot compare strings with '<=' operator", leftFinalValue.line)
		}
		return token{}, fmt.Errorf("line %d, cannot apply '<=' to %s and %s, only between numbers", op.line, leftFinalValue.name, rightFinalValue.name)
	}

	return token{}, fmt.Errorf("line %d: unknown operator '%s'", op.line, op.value)
}

func (i *interpreter) executePrintStmt(expression Expr) error {
	token, err := i.evaluate(expression)
	if err != nil {
		return err
	}
	println(token.value)
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
