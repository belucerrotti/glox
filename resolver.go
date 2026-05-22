package main

import "fmt"

type resolver struct {
	scopes     []map[string]bool
	distances  map[Expr]int
	inFunction bool
}

func createResolver() *resolver {
	return &resolver{
		scopes:     []map[string]bool{},
		distances:  map[Expr]int{},
		inFunction: false,
	}
}

func (r *resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *resolver) declare(name string) {
	if len(r.scopes) == 0 {
		return // global scope, no se trackea
	}
	r.scopes[len(r.scopes)-1][name] = false
}

func (r *resolver) define(name string) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name] = true
}

func (r *resolver) resolveDistance(expr Expr, name string) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name]; ok {
			r.distances[expr] = len(r.scopes) - 1 - i
			return
		}
	}
}

func (r *resolver) resolve(statements []Stmt) error {
	for _, stmt := range statements {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *resolver) resolveStmt(stmt Stmt) error {
	switch s := stmt.(type) {
	case *VarDecl:
		return r.resolveVarDecl(s)
	case *FunDecl:
		return r.resolveFunDecl(s)
	case *ClassDecl:
		return r.resolveClassDecl(s)
	case *BlockStmt:
		return r.resolveBlockStmt(s)
	case *ExpressionStmt:
		return r.resolveExpr(s.expr)
	case *PrintStmt:
		return r.resolveExpr(s.expr)
	case *ReturnStmt:
		return r.resolveReturnStmt(s)
	case *IfStmt:
		return r.resolveIfStmt(s)
	case *WhileStmt:
		return r.resolveWhileStmt(s)
	case *ForStmt:
		return r.resolveForStmt(s)
	}
	return nil
}

func (r *resolver) resolveClassDecl(s *ClassDecl) error {
	r.declare(s.name.value)
	r.define(s.name.value)
	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true
	for _, method := range s.methods {
		if err := r.resolveFunction(method); err != nil {
			return err
		}
	}
	r.endScope()
	return nil
}

func (r *resolver) resolveVarDecl(s *VarDecl) error {
	r.declare(s.name.value)
	if s.value != nil {
		if err := r.resolveExpr(s.value); err != nil {
			return err
		}
	}
	r.define(s.name.value)
	return nil
}

func (r *resolver) resolveFunDecl(s *FunDecl) error {
	r.declare(s.name.value)
	r.define(s.name.value)
	return r.resolveFunction(s)
}

func (r *resolver) resolveFunction(s *FunDecl) error {
	wasInFunction := r.inFunction
	r.inFunction = true
	defer func() { r.inFunction = wasInFunction }()

	r.beginScope()
	for _, param := range s.parameters {
		r.declare(param.value)
		r.define(param.value)
	}
	for _, stmt := range s.body {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	r.endScope()
	return nil
}

func (r *resolver) resolveBlockStmt(s *BlockStmt) error {
	r.beginScope()
	for _, stmt := range s.statements {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	r.endScope()
	return nil
}

func (r *resolver) resolveReturnStmt(s *ReturnStmt) error {
	if !r.inFunction {
		return fmt.Errorf("cannot use 'return' outside a function")
	}
	if s.value != nil {
		return r.resolveExpr(s.value)
	}
	return nil
}

func (r *resolver) resolveIfStmt(s *IfStmt) error {
	if err := r.resolveExpr(s.condition); err != nil {
		return err
	}
	if err := r.resolveStmt(s.thenBranch); err != nil {
		return err
	}
	if s.elseBranch != nil {
		return r.resolveStmt(s.elseBranch)
	}
	return nil
}

func (r *resolver) resolveWhileStmt(s *WhileStmt) error {
	if err := r.resolveExpr(s.condition); err != nil {
		return err
	}
	return r.resolveStmt(s.body)
}

func (r *resolver) resolveForStmt(s *ForStmt) error {
	r.beginScope()
	defer r.endScope()

	if s.initializer != nil {
		if err := r.resolveStmt(s.initializer); err != nil {
			return err
		}
	}
	if s.condition != nil {
		if err := r.resolveExpr(s.condition); err != nil {
			return err
		}
	}
	if s.increment != nil {
		if err := r.resolveExpr(s.increment); err != nil {
			return err
		}
	}
	return r.resolveStmt(s.body)
}

func (r *resolver) resolveExpr(expr Expr) error {
	switch e := expr.(type) {
	case *VariableExpr:
		if len(r.scopes) > 0 {
			if initialized, exists := r.scopes[len(r.scopes)-1][e.name.value]; exists && !initialized {
				return fmt.Errorf("line %d: cannot read local variable '%s' before its initialization", e.name.line, e.name.value)
			}
		}
		r.resolveDistance(e, e.name.value)
	case *AssignExpr:
		if err := r.resolveExpr(e.value); err != nil {
			return err
		}
		r.resolveDistance(e, e.name.value)
	case *BinaryExpr:
		if err := r.resolveExpr(e.left); err != nil {
			return err
		}
		return r.resolveExpr(e.right)
	case *UnaryExpr:
		return r.resolveExpr(e.right)
	case *GroupingExpr:
		return r.resolveExpr(e.expression)
	case *LogicalExpr:
		if err := r.resolveExpr(e.left); err != nil {
			return err
		}
		return r.resolveExpr(e.right)
	case *CallExpr:
		if err := r.resolveExpr(e.callee); err != nil {
			return err
		}
		for _, arg := range e.arguments {
			if err := r.resolveExpr(arg); err != nil {
				return err
			}
		}
	case *GetExpr:
		return r.resolveExpr(e.object)
	case *SetExpr:
		if err := r.resolveExpr(e.object); err != nil {
			return err
		}
		return r.resolveExpr(e.value)
	case *LiteralExpr:
	}
	return nil
}
