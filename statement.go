package main

type Stmt interface {
	stmtNode()
}

type ExpressionStmt struct {
	expr Expr
}

type PrintStmt struct {
	expr Expr
}

type VarDecl struct {
	name  token
	value Expr // puede ser nil
}

type BlockStmt struct {
	statements []Stmt
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt // puede ser nil
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

type ForStmt struct {
	initializer Stmt // puede ser nil
	condition   Expr // puede ser nil
	increment   Expr // puede ser nil
	body        Stmt
}

type FunDecl struct {
	name       token
	parameters []token
	body       []Stmt
}

type ReturnStmt struct {
	value Expr
} // value puede ser nil

func (e *ExpressionStmt) stmtNode() {}
func (p *PrintStmt) stmtNode()      {}
func (v *VarDecl) stmtNode()        {}
func (b *BlockStmt) stmtNode()      {}
func (i *IfStmt) stmtNode()         {}
func (w *WhileStmt) stmtNode()      {}
func (fr *ForStmt) stmtNode()       {}
func (f *FunDecl) stmtNode()        {}
func (r *ReturnStmt) stmtNode()     {}
