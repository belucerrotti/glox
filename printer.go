// este archivo es simplemente para que se imprima más legible al usar la opción "--parser"

package main

import (
	"fmt"
	"strings"
)

func printExpr(expr Expr, indent int) string {
	pad := strings.Repeat("  ", indent)
	inner := pad + "  "
	switch e := expr.(type) {
	case *LiteralExpr:
		return fmt.Sprintf("LiteralExpr { value: %s }", e.value.value)
	case *VariableExpr:
		return fmt.Sprintf("VariableExpr { name: %s }", e.name.value)
	case *UnaryExpr:
		return fmt.Sprintf("UnaryExpr {\n%sop: %s\n%sright: %s\n%s}",
			inner, e.operator.value,
			inner, printExpr(e.right, indent+1),
			pad)
	case *BinaryExpr:
		return fmt.Sprintf("BinaryExpr {\n%sleft: %s\n%sop: %s\n%sright: %s\n%s}",
			inner, printExpr(e.left, indent+1),
			inner, e.operator.value,
			inner, printExpr(e.right, indent+1),
			pad)
	case *GroupingExpr:
		return fmt.Sprintf("GroupingExpr {\n%sexpr: %s\n%s}",
			inner, printExpr(e.expression, indent+1),
			pad)
	case *AssignExpr:
		return fmt.Sprintf("AssignExpr {\n%sname: %s\n%svalue: %s\n%s}",
			inner, e.name.value,
			inner, printExpr(e.value, indent+1),
			pad)
	default:
		return "UnknownExpr"
	}
}

func printStmt(stmt Stmt, indent int) string {
	pad := strings.Repeat("  ", indent)
	inner := pad + "  "
	switch s := stmt.(type) {
	case *PrintStmt:
		return fmt.Sprintf("PrintStmt {\n%sexpr: %s\n%s}", inner, printExpr(s.expr, indent+1), pad)
	case *ExpressionStmt:
		return fmt.Sprintf("ExpressionStmt {\n%sexpr: %s\n%s}", inner, printExpr(s.expr, indent+1), pad)
	case *ReturnStmt:
		if s.value == nil {
			return "ReturnStmt { nil }"
		}
		return fmt.Sprintf("ReturnStmt {\n%svalue: %s\n%s}", inner, printExpr(s.value, indent+1), pad)
	case *VarDecl:
		if s.value == nil {
			return fmt.Sprintf("VarDecl { name: %s, value: nil }", s.name.value)
		}
		return fmt.Sprintf("VarDecl {\n%sname: %s\n%svalue: %s\n%s}", inner, s.name.value, inner, printExpr(s.value, indent+1), pad)
	case *BlockStmt:
		lines := []string{"BlockStmt {"}
		for i, st := range s.statements {
			lines = append(lines, fmt.Sprintf("%s[%d]: %s", inner, i, printStmt(st, indent+1)))
		}
		lines = append(lines, pad+"}")
		return strings.Join(lines, "\n")
	case *IfStmt:
		result := fmt.Sprintf("IfStmt {\n%scondition: %s\n%sthen: %s",
			inner, printExpr(s.condition, indent+1),
			inner, printStmt(s.thenBranch, indent+1))
		if s.elseBranch != nil {
			result += fmt.Sprintf("\n%selse: %s", inner, printStmt(s.elseBranch, indent+1))
		}
		return result + "\n" + pad + "}"
	case *FunDecl:
		params := []string{}
		for _, p := range s.parameters {
			params = append(params, p.value)
		}
		lines := []string{"FunDecl {"}
		lines = append(lines, fmt.Sprintf("%sname: %s", inner, s.name.value))
		lines = append(lines, fmt.Sprintf("%sparams: [%s]", inner, strings.Join(params, ", ")))
		lines = append(lines, fmt.Sprintf("%sbody:", inner))
		for i, st := range s.body {
			lines = append(lines, fmt.Sprintf("%s  [%d]: %s", inner, i, printStmt(st, indent+2)))
		}
		lines = append(lines, pad+"}")
		return strings.Join(lines, "\n")
	case *ForStmt:
		init := "nil"
		if s.initializer != nil {
			init = printStmt(s.initializer, indent+1)
		}
		cond := "nil"
		if s.condition != nil {
			cond = printExpr(s.condition, indent+1)
		}
		inc := "nil"
		if s.increment != nil {
			inc = printExpr(s.increment, indent+1)
		}
		return fmt.Sprintf("ForStmt {\n%sinitializer: %s\n%scondition: %s\n%sincrement: %s\n%sbody: %s\n%s}",
			inner, init,
			inner, cond,
			inner, inc,
			inner, printStmt(s.body, indent+1),
			pad)
	case *WhileStmt:
		return fmt.Sprintf("WhileStmt {\n%scondition: %s\n%sbody: %s\n%s}",
			inner, printExpr(s.condition, indent+1),
			inner, printStmt(s.body, indent+1),
			pad)
	default:
		return "UnknownStmt"
	}
}
