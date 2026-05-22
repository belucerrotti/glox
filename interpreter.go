package main

import (
	"fmt"
)

type returnValue struct {
	value token
}

func (r returnValue) Error() string { return "" }

type interpreter struct {
	statements  []Stmt
	current     int
	environment *environment
	distances   map[Expr]int
}

type environment struct {
	variables map[string]interface{}
	father    *environment
}

type loxFunction struct {
	parameters []token
	body       []Stmt
	closure    *environment
}

type loxClass struct {
	name    string
	methods map[string]loxFunction
}

type loxInstance struct {
	class  *loxClass
	fields map[string]interface{}
}

func createEnvironment(father *environment) *environment {
	return &environment{
		variables: map[string]interface{}{},
		father:    father,
	}
}

func createInterpreter(statements []Stmt, distances map[Expr]int) *interpreter {
	return &interpreter{
		statements:  statements,
		current:     0,
		environment: createEnvironment(nil),
		distances:   distances,
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
	case *ForStmt:
		err = i.executeForStmt(s)
	case *FunDecl:
		err = i.executeFunDecl(s)
	case *ReturnStmt:
		err = i.executeReturnStmt(s)
	case *ClassDecl:
		err = i.executeClassDecl(s)
	}

	if err != nil {
		return err
	}

	return nil
}

func (i *interpreter) executeClassDecl(c *ClassDecl) error {
	methods := make(map[string]loxFunction)
	for _, method := range c.methods {
		methods[method.name.value] = loxFunction{
			parameters: method.parameters,
			body:       method.body,
			closure:    i.environment,
		}
	}
	i.environment.variables[c.name.value] = loxClass{name: c.name.value, methods: methods}
	return nil
}

func (i *interpreter) executeReturnStmt(s *ReturnStmt) error {
	if s.value == nil {
		return returnValue{value: token{tokenType: NIL, value: "nil"}}
	}

	value, err := i.evaluate(s.value)
	if err != nil {
		return err
	}
	return returnValue{value: value}
}

func (i *interpreter) executeFunDecl(s *FunDecl) error {
	i.environment.variables[s.name.value] = loxFunction{parameters: s.parameters, body: s.body, closure: i.environment}
	return nil
}

func (i *interpreter) executeForStmt(s *ForStmt) error {
	previous := i.environment
	i.environment = createEnvironment(previous)
	defer func() { i.environment = previous }()

	if s.initializer != nil {
		err := i.execute(s.initializer)
		if err != nil {
			return err
		}
	}

	for {
		var cond bool
		var err error
		if s.condition == nil {
			cond = true
		} else {
			cond, err = i.evaluateCondition(s.condition)
			if err != nil {
				return err
			}
		}

		if cond {
			err := i.execute(s.body)
			if err != nil {
				return err
			}

			if s.increment != nil {
				_, e := i.evaluate(s.increment)
				if e != nil {
					return e
				}
			}
		} else {
			return nil
		}
	}
}

func isTruthy(t token) bool {
	return t.tokenType != FALSE && t.tokenType != NIL
}

func (i *interpreter) evaluateCondition(e Expr) (bool, error) {
	condition, err := i.evaluate(e)
	if err != nil {
		return false, err
	}

	return isTruthy(condition), nil
}

func (i *interpreter) executeWhileStmt(s *WhileStmt) error {
	for {
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
	} else if s.elseBranch != nil {
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
		token, err = i.evaluateVariableExpression(e.name, e)
	case *AssignExpr:
		token, err = i.evaluateAssignExpression(e.name, e.value, e)
	case *LogicalExpr:
		token, err = i.evaluateLogicalExpression(e.left, e.operator, e.right)
	case *CallExpr:
		token, err = i.evaluateCallExpression(e.callee, e.paren, e.arguments)
	case *GetExpr:
		token, err = i.evaluateGetExpr(e.object, e.name)
	case *SetExpr:
		token, err = i.evaluateSetExpr(e.object, e.name, e.value)
	}

	if err != nil {
		return token, err
	}
	return token, nil
}

func (i *interpreter) evaluateGetExpr(object Expr, name token) (token, error) {
	ob, err := i.evaluate(object)
	if err != nil {
		return token{}, err
	}
	if ob.instance == nil {
		return token{}, fmt.Errorf("line %d: only instances have properties", name.line)
	}
	instance := ob.instance
	if value, ok := instance.fields[name.value]; ok {
		switch v := value.(type) {
		case token:
			return v, nil
		case loxFunction:
			return token{value: name.value, line: name.line, callable: &v}, nil
		default:
			return token{}, fmt.Errorf("line %d: property '%s' has an invalid value", name.line, name.value)
		}
	}

	if method, ok := instance.class.methods[name.value]; ok {
		thisEnv := createEnvironment(method.closure)
		thisEnv.variables["this"] = ob
		m := loxFunction{parameters: method.parameters, body: method.body, closure: thisEnv}
		return token{value: name.value, line: name.line, callable: &m}, nil
	}
	return token{}, fmt.Errorf("line %d: undefined property '%s'", name.line, name.value)
}

func (i *interpreter) evaluateSetExpr(object Expr, name token, value Expr) (token, error) {
	ob, err := i.evaluate(object)
	if err != nil {
		return token{}, err
	}
	if ob.instance == nil {
		return token{}, fmt.Errorf("line %d: only instances have properties", name.line)
	}
	instance := ob.instance

	val, err := i.evaluate(value)
	if err != nil {
		return token{}, err
	}

	// crea el campo si no existe, o lo sobreescribe si ya existe
	instance.fields[name.value] = val
	return val, nil
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
	name, err := i.evaluate(callee)
	if err != nil {
		return token{}, err
	}

	if name.callable != nil {
		return i.callFunction(name.callable, name.value, paren, arguments)
	}

	if name.class != nil {
		return i.instanceClass(name.class, name.value, paren, arguments)
	}

	env := i.environment
	for env != nil {
		if value, ok := env.variables[name.value]; ok {
			switch v := value.(type) {
			case loxFunction:
				return i.callFunction(&v, name.value, paren, arguments)
			case loxClass:
				return i.instanceClass(&v, name.value, paren, arguments)
			default:
				return token{}, fmt.Errorf("line %d: '%s' is not a function", paren.line, name.value)
			}
		}
		env = env.father
	}
	return token{}, fmt.Errorf("line %d: function '%s' is not defined", paren.line, name.value)
}

func (i *interpreter) instanceClass(v *loxClass, name string, paren token, arguments []Expr) (token, error) {
	instance := &loxInstance{class: v, fields: map[string]interface{}{}}
	instanceToken := token{tokenType: IDENTIFIER, value: name, line: paren.line, instance: instance}

	if init, ok := v.methods["init"]; ok {
		_, err := i.callMethod(&init, instanceToken, paren, arguments)
		if err != nil {
			return token{}, err
		}
	} else if len(arguments) > 0 {
		return token{}, fmt.Errorf("line %d: class '%s' expects 0 argument(s) but got %d", paren.line, name, len(arguments))
	}

	return instanceToken, nil
}

func (i *interpreter) callMethod(v *loxFunction, receiver token, paren token, arguments []Expr) (token, error) {
	if len(arguments) != len(v.parameters) {
		return token{}, fmt.Errorf("line %d: method expects %d argument(s) but got %d", paren.line, len(v.parameters), len(arguments))
	}

	thisEnv := createEnvironment(v.closure)
	thisEnv.variables["this"] = receiver
	e := createEnvironment(thisEnv)
	for idx, param := range v.parameters {
		argVal, err := i.evaluate(arguments[idx])
		if err != nil {
			return token{}, err
		}
		e.variables[param.value] = argVal
	}
	prev := i.environment
	i.environment = e
	var execErr error
	for _, stmt := range v.body {
		if err := i.execute(stmt); err != nil {
			execErr = err
			break
		}
	}
	i.environment = prev
	if rv, ok := execErr.(returnValue); ok {
		return rv.value, nil
	}
	if execErr != nil {
		return token{}, execErr
	}
	return receiver, nil
}

func (i *interpreter) callFunction(v *loxFunction, name string, paren token, arguments []Expr) (token, error) {
	if len(arguments) != len(v.parameters) {
		return token{}, fmt.Errorf("line %d: function '%s' expects %d argument(s) but got %d", paren.line, name, len(v.parameters), len(arguments))
	}
	e := createEnvironment(v.closure)
	for idx := 0; idx < len(arguments); idx++ {
		argVal, err := i.evaluate(arguments[idx])
		if err != nil {
			return token{}, err
		}
		e.variables[v.parameters[idx].value] = argVal
	}
	prev := i.environment
	i.environment = e
	var execErr error
	for _, stmt := range v.body {
		if err := i.execute(stmt); err != nil {
			execErr = err
			break
		}
	}
	i.environment = prev
	if rv, ok := execErr.(returnValue); ok {
		return rv.value, nil
	}
	if execErr != nil {
		return token{}, execErr
	}
	return token{tokenType: NIL, value: "nil"}, nil
}

func (i *interpreter) environmentAt(distance int) *environment {
	env := i.environment
	for d := 0; d < distance; d++ {
		env = env.father
	}
	return env
}

func (i *interpreter) lookupVariable(name token, expr Expr) (token, error) {
	var env *environment
	if distance, ok := i.distances[expr]; ok {
		env = i.environmentAt(distance)
	} else {
		env = i.environment
		for env.father != nil {
			env = env.father
		}
	}

	if value, ok := env.variables[name.value]; ok {
		switch stored := value.(type) {
		case token:
			return stored, nil
		case loxFunction:
			return token{value: name.value, line: name.line, callable: &stored}, nil
		case loxClass:
			return token{value: name.value, line: name.line, class: &stored}, nil
		default:
			return token{}, fmt.Errorf("line %d: variable '%s' has an invalid value", name.line, name.value)
		}
	}
	return token{}, fmt.Errorf("line %d: variable '%s' is not defined", name.line, name.value)
}

func (i *interpreter) evaluateVariableExpression(name token, expr Expr) (token, error) {
	return i.lookupVariable(name, expr)
}

func (i *interpreter) evaluateAssignExpression(name token, value Expr, expr Expr) (token, error) {
	valueToken, err := i.evaluate(value)
	if err != nil {
		return token{}, err
	}

	var env *environment
	if distance, ok := i.distances[expr]; ok {
		env = i.environmentAt(distance)
	} else {
		env = i.environment
		for env.father != nil {
			env = env.father
		}
	}

	if _, ok := env.variables[name.value]; ok {
		env.variables[name.value] = valueToken
		return valueToken, nil
	}

	return token{}, fmt.Errorf("line %d: variable '%s' is not defined", name.line, name.value)
}

func (i *interpreter) evaluateLogicalExpression(left Expr, op token, right Expr) (token, error) {
	leftValue, err := i.evaluate(left)
	if err != nil {
		return token{}, err
	}

	truthy := isTruthy(leftValue)

	switch op.tokenType {
	case AND:
		if !truthy {
			return leftValue, nil
		}
	case OR:
		if truthy {
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
		if !isTruthy(value) {
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
	case MOD:
		if numericOperation {
			if rightFinalValue.valueFloat == 0 {
				return token{}, fmt.Errorf("line %d: modulo by zero", op.line)
			}
			result := float64(int(leftFinalValue.valueFloat) % int(rightFinalValue.valueFloat))
			return token{tokenType: NUMBER, valueFloat: result, value: fmt.Sprintf("%g", result), line: op.line}, nil
		}
		return token{}, fmt.Errorf("line %d: '%%' only applies to numbers", op.line)
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
		if _, ok := err.(returnValue); ok {
			return fmt.Errorf("cannot use 'return' outside a function")
		}
		if err != nil {
			return err
		}

		i.current++
	}

	return nil
}
