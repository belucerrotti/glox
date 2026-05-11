package main

type Expr interface {
	exprNode()
}

// expression operator expression
type BinaryExpr struct {
	left     Expr
	operator token
	right    Expr
}

// operator expression
type UnaryExpr struct {
	operator token
	right    Expr
}

// ( expression )
type GroupingExpr struct {
	expression Expr
}

// number | string | true | false | nil
type LiteralExpr struct {
	value token
}

type VariableExpr struct {
	name token
}

type AssignExpr struct {
	name  token
	value Expr
}

func (b *BinaryExpr) exprNode()   {}
func (u *UnaryExpr) exprNode()    {}
func (g *GroupingExpr) exprNode() {}
func (l *LiteralExpr) exprNode()  {}
func (v *VariableExpr) exprNode() {}
func (a *AssignExpr) exprNode()   {}
