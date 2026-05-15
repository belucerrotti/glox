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
	case *ExpressionStmt:
		_, err = i.evaluate(s.expr)
	case *VarDecl:
		err = i.executeVarDecl(s)
	case *BlockStmt:
		err = i.executeBlockStmt(s.statements)
	case *IfStmt:
		err = i.executeIfStmt(s)
	case *WhileStmt:
		err = i.executeWhileStmt(s)
		// case *ForStmt:
		// 	err = i.executeForStmt(s)
		// case *FunDecl:
		// 	err = i.executeFunDecl(s)
		// case *ReturnStmt:
		// 	err = i.executeReturnStmt(s)
	}

	if err != nil {
		return err
	}

	return nil
}

func (i *interpreter) evaluateCondition(e Expr) (bool, error) {
	condition, err := i.evaluate(e)
	if err != nil {
		return false, err
	}

	return condition.tokenType == TRUE, nil
}

func (i *interpreter) executeWhileStmt(s *WhileStmt) error {
	for true {
		cond, err := i.evaluateCondition(s.condition)
		if err != nil {
			return err
		}
		if cond {
			err := i.execute(s.body)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
	}
	return nil
}

func (i *interpreter) executeIfStmt(s *IfStmt) error {
	condition, err := i.evaluateCondition(s.condition)
	if err != nil {
		return err
	}

	if condition {
		err := i.execute(s.thenBranch)
		if err != nil {
			return err
		}
	} else if s.thenBranch != nil {
		err := i.execute(s.elseBranch)
		if err != nil {
			return err
		}
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
	case *VariableExpr:
		token, err = i.evaluateVariableExpression(e.name)
	case *AssignExpr:
		token, err = i.evaluateAssignExpression(e.name, e.value)
	case *LogicalExpr:
		token, err = i.evaluateLogicalExpression(e.left, e.operator, e.right)
		// case *CallExpr:
		// 	token, err = i.evaluateCallExpression(e.callee, e.paren, e.arguments)
	}

	if err != nil {
		return token, err
	}
	return token, nil
}

func (i *interpreter) executeVarDecl(s *VarDecl) error {
	value := token{tokenType: NIL, value: "nil"}
	if s.value != nil {
		v, err := i.evaluate(s.value)
		if err != nil {
			return err
		}
		value = v
	}
	i.environment.variables[s.name.value] = value
	return nil
}

func (i *interpreter) executeBlockStmt(statements []Stmt) error {
	previous := i.environment
	i.environment = createEnvironment(previous)
	defer func() { i.environment = previous }()

	for _, stmt := range statements {
		if err := i.execute(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (i *interpreter) evaluateCallExpression(callee Expr, paren token, arguments []Expr) (token, error) {
	// todo
	return token{}, nil
}

func (i *interpreter) evaluateVariableExpression(name token) (token, error) {
	env := i.environment
	for env != nil {
		if value, ok := env.variables[name.value]; ok {
			switch v := value.(type) {
			case token:
				return v, nil
			default:
				return token{}, fmt.Errorf("line %d: variable '%s' has an invalid value", name.line, name.value)
			}
		}
		env = env.father
	}
	return token{}, fmt.Errorf("line %d: variable '%s' is not defined", name.line, name.value)
}

func (i *interpreter) evaluateAssignExpression(name token, value Expr) (token, error) {
	valueToken, err := i.evaluate(value)
	if err != nil {
		return token{}, err
	}

	env := i.environment
	for env != nil {
		if _, ok := env.variables[name.value]; ok {
			env.variables[name.value] = valueToken
			return valueToken, nil
		}
		env = env.father
	}

	return token{}, fmt.Errorf("line %d: variable '%s' is not defined", name.line, name.value)
}

func (i *interpreter) evaluateLogicalExpression(left Expr, op token, right Expr) (token, error) {
	leftValue, err := i.evaluate(left)
	if err != nil {
		return token{}, err
	}

	isTruthy := leftValue.tokenType != FALSE && leftValue.tokenType != NIL

	switch op.tokenType {
	case AND:
		if !isTruthy {
			return leftValue, nil
		}
	case OR:
		if isTruthy {
			return leftValue, nil
		}
	default:
		return token{}, fmt.Errorf("line %d: unknown logical operator '%s'", op.line, op.value)
	}

	return i.evaluate(right)
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
