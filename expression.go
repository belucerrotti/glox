package main

type Expr interface {
	exprNode()
}

// expression operator expression
// ej: 1 + 2, "hola" == "chau", 3 * (4 - 1)
type BinaryExpr struct {
	left     Expr
	operator token
	right    Expr
}

// operator expression
// ej: -5, !true
type UnaryExpr struct {
	operator token
	right    Expr
}

// ( expression )
// ej: (1 + 2)
type GroupingExpr struct {
	expression Expr
}

// number | string | true | false | nil
// ej: 42, "hola", true
type LiteralExpr struct {
	value token
}

// and/or (si el left ya define, no se lee el right)
// ej: (3 > 5) and (4 < 2)
type LogicalExpr struct {
	left     Expr
	operator token
	right    Expr
}

// llamada a funcion
// ej: suma(1, 2)
type CallExpr struct {
	callee    Expr
	paren     token
	arguments []Expr
}

// lectura de variable
// ej: x, miVariable
type VariableExpr struct {
	name token
}

// asignacion de variable
// ej: x = 5, nombre = "lox"
type AssignExpr struct {
	name  token
	value Expr
}

func (b *BinaryExpr) exprNode()   {}
func (u *UnaryExpr) exprNode()    {}
func (g *GroupingExpr) exprNode() {}
func (l *LiteralExpr) exprNode()  {}
func (l *LogicalExpr) exprNode()  {}
func (c *CallExpr) exprNode()     {}
func (v *VariableExpr) exprNode() {}
func (a *AssignExpr) exprNode()   {}
