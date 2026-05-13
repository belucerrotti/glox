package main

import (
	"strings"
	"testing"
)

// helper: escanea y parsea un string, retorna los statements o el error
func parseSource(src string) ([]Stmt, error) {
	s := createScanner()
	tokens, err := s.scan([]byte(src))
	if err != nil {
		return nil, err
	}
	p := createParser(tokens)
	return p.parse()
}

// helper: verifica que el parseo sea exitoso y retorna el primer statement
func mustParse(t *testing.T, src string) []Stmt {
	t.Helper()
	stmts, err := parseSource(src)
	if err != nil {
		t.Fatalf("unexpected error parsing %q: %v", src, err)
	}
	return stmts
}

// helper: verifica que el parseo falle con un error que contenga substr
func mustFail(t *testing.T, src string, substr string) {
	t.Helper()
	_, err := parseSource(src)
	if err == nil {
		t.Fatalf("expected error parsing %q but got none", src)
	}
	if !strings.Contains(err.Error(), substr) {
		t.Fatalf("expected error containing %q, got: %v", substr, err)
	}
}

// ---------- Literales ----------

func TestLiteralNumber(t *testing.T) {
	stmts := mustParse(t, "42;")
	if len(stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(stmts))
	}
	expr, ok := stmts[0].(*ExpressionStmt)
	if !ok {
		t.Fatalf("expected ExpressionStmt")
	}
	lit, ok := expr.expr.(*LiteralExpr)
	if !ok {
		t.Fatalf("expected LiteralExpr")
	}
	if lit.value.value != "42" {
		t.Fatalf("expected value 42, got %s", lit.value.value)
	}
}

func TestLiteralString(t *testing.T) {
	stmts := mustParse(t, `"hola";`)
	expr := stmts[0].(*ExpressionStmt)
	lit := expr.expr.(*LiteralExpr)
	if lit.value.value != "hola" {
		t.Fatalf("expected 'hola', got %s", lit.value.value)
	}
}

func TestLiteralTrue(t *testing.T) {
	stmts := mustParse(t, "true;")
	expr := stmts[0].(*ExpressionStmt)
	lit := expr.expr.(*LiteralExpr)
	if lit.value.tokenType != TRUE {
		t.Fatalf("expected TRUE token")
	}
}

func TestLiteralFalse(t *testing.T) {
	stmts := mustParse(t, "false;")
	expr := stmts[0].(*ExpressionStmt)
	lit := expr.expr.(*LiteralExpr)
	if lit.value.tokenType != FALSE {
		t.Fatalf("expected FALSE token")
	}
}

func TestLiteralNil(t *testing.T) {
	stmts := mustParse(t, "nil;")
	expr := stmts[0].(*ExpressionStmt)
	lit := expr.expr.(*LiteralExpr)
	if lit.value.tokenType != NIL {
		t.Fatalf("expected NIL token")
	}
}

// ---------- Expresiones binarias ----------

func TestBinaryAdd(t *testing.T) {
	stmts := mustParse(t, "1 + 2;")
	expr := stmts[0].(*ExpressionStmt)
	bin := expr.expr.(*BinaryExpr)
	if bin.operator.tokenType != PLUS {
		t.Fatalf("expected PLUS operator")
	}
}

func TestBinaryChain(t *testing.T) {
	// 1 + 2 + 3 debe parsear como (1 + 2) + 3 (asociatividad izquierda)
	stmts := mustParse(t, "1 + 2 + 3;")
	expr := stmts[0].(*ExpressionStmt)
	outer := expr.expr.(*BinaryExpr)
	_, ok := outer.left.(*BinaryExpr)
	if !ok {
		t.Fatalf("expected left to be BinaryExpr (left associativity)")
	}
}

func TestBinaryPrecedence(t *testing.T) {
	// 1 + 2 * 3 debe parsear como 1 + (2 * 3)
	stmts := mustParse(t, "1 + 2 * 3;")
	expr := stmts[0].(*ExpressionStmt)
	outer := expr.expr.(*BinaryExpr)
	if outer.operator.tokenType != PLUS {
		t.Fatalf("expected outer operator to be PLUS")
	}
	_, ok := outer.right.(*BinaryExpr)
	if !ok {
		t.Fatalf("expected right to be BinaryExpr (higher precedence)")
	}
}

func TestGroupingOverridesPrecedence(t *testing.T) {
	// (1 + 2) * 3: el * debe ser el operador externo
	stmts := mustParse(t, "(1 + 2) * 3;")
	expr := stmts[0].(*ExpressionStmt)
	outer := expr.expr.(*BinaryExpr)
	if outer.operator.tokenType != STAR {
		t.Fatalf("expected outer operator to be STAR")
	}
	_, ok := outer.left.(*GroupingExpr)
	if !ok {
		t.Fatalf("expected left to be GroupingExpr")
	}
}

func TestComparison(t *testing.T) {
	mustParse(t, "a < b;")
	mustParse(t, "a <= b;")
	mustParse(t, "a > b;")
	mustParse(t, "a >= b;")
}

func TestEquality(t *testing.T) {
	mustParse(t, "a == b;")
	mustParse(t, "a != b;")
}

// ---------- Unario ----------

func TestUnaryMinus(t *testing.T) {
	stmts := mustParse(t, "-1;")
	expr := stmts[0].(*ExpressionStmt)
	u, ok := expr.expr.(*UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr")
	}
	if u.operator.tokenType != MINUS {
		t.Fatalf("expected MINUS operator")
	}
}

func TestUnaryBang(t *testing.T) {
	stmts := mustParse(t, "!true;")
	expr := stmts[0].(*ExpressionStmt)
	u, ok := expr.expr.(*UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr")
	}
	if u.operator.tokenType != BANG {
		t.Fatalf("expected BANG operator")
	}
}

func TestDoubleUnary(t *testing.T) {
	// !!true debe parsear como !(!true)
	stmts := mustParse(t, "!!true;")
	expr := stmts[0].(*ExpressionStmt)
	outer, ok := expr.expr.(*UnaryExpr)
	if !ok {
		t.Fatalf("expected outer UnaryExpr")
	}
	_, ok = outer.right.(*UnaryExpr)
	if !ok {
		t.Fatalf("expected inner UnaryExpr")
	}
}

// ---------- Variables ----------

func TestVariableExpr(t *testing.T) {
	stmts := mustParse(t, "x;")
	expr := stmts[0].(*ExpressionStmt)
	_, ok := expr.expr.(*VariableExpr)
	if !ok {
		t.Fatalf("expected VariableExpr")
	}
}

func TestAssignExpr(t *testing.T) {
	stmts := mustParse(t, "x = 5;")
	expr := stmts[0].(*ExpressionStmt)
	a, ok := expr.expr.(*AssignExpr)
	if !ok {
		t.Fatalf("expected AssignExpr")
	}
	if a.name.value != "x" {
		t.Fatalf("expected name 'x', got %s", a.name.value)
	}
}

func TestChainedAssign(t *testing.T) {
	// a = b = 5 debe parsear como a = (b = 5)
	stmts := mustParse(t, "a = b = 5;")
	expr := stmts[0].(*ExpressionStmt)
	outer := expr.expr.(*AssignExpr)
	_, ok := outer.value.(*AssignExpr)
	if !ok {
		t.Fatalf("expected right-associative chained assign")
	}
}

// ---------- Lógicos ----------

func TestLogicalAnd(t *testing.T) {
	stmts := mustParse(t, "a and b;")
	expr := stmts[0].(*ExpressionStmt)
	l, ok := expr.expr.(*LogicalExpr)
	if !ok {
		t.Fatalf("expected LogicalExpr")
	}
	if l.operator.tokenType != AND {
		t.Fatalf("expected AND operator")
	}
}

func TestLogicalOr(t *testing.T) {
	stmts := mustParse(t, "a or b;")
	expr := stmts[0].(*ExpressionStmt)
	l, ok := expr.expr.(*LogicalExpr)
	if !ok {
		t.Fatalf("expected LogicalExpr")
	}
	if l.operator.tokenType != OR {
		t.Fatalf("expected OR operator")
	}
}

// ---------- Llamadas ----------

func TestCallNoArgs(t *testing.T) {
	stmts := mustParse(t, "foo();")
	expr := stmts[0].(*ExpressionStmt)
	c, ok := expr.expr.(*CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr")
	}
	if len(c.arguments) != 0 {
		t.Fatalf("expected 0 arguments")
	}
}

func TestCallWithArgs(t *testing.T) {
	stmts := mustParse(t, "foo(1, 2, 3);")
	expr := stmts[0].(*ExpressionStmt)
	c := expr.expr.(*CallExpr)
	if len(c.arguments) != 3 {
		t.Fatalf("expected 3 arguments, got %d", len(c.arguments))
	}
}

func TestCallChained(t *testing.T) {
	// foo()() debe parsear como CallExpr{ callee: CallExpr{...} }
	stmts := mustParse(t, "foo()();")
	expr := stmts[0].(*ExpressionStmt)
	outer := expr.expr.(*CallExpr)
	_, ok := outer.callee.(*CallExpr)
	if !ok {
		t.Fatalf("expected chained CallExpr")
	}
}

// ---------- Statements ----------

func TestPrintStmt(t *testing.T) {
	stmts := mustParse(t, "print 42;")
	_, ok := stmts[0].(*PrintStmt)
	if !ok {
		t.Fatalf("expected PrintStmt")
	}
}

func TestVarDeclNoValue(t *testing.T) {
	stmts := mustParse(t, "var x;")
	v, ok := stmts[0].(*VarDecl)
	if !ok {
		t.Fatalf("expected VarDecl")
	}
	if v.value != nil {
		t.Fatalf("expected nil value")
	}
}

func TestVarDeclWithValue(t *testing.T) {
	stmts := mustParse(t, "var x = 5;")
	v, ok := stmts[0].(*VarDecl)
	if !ok {
		t.Fatalf("expected VarDecl")
	}
	if v.value == nil {
		t.Fatalf("expected non-nil value")
	}
}

func TestBlockStmt(t *testing.T) {
	stmts := mustParse(t, "{ var x = 1; print x; }")
	b, ok := stmts[0].(*BlockStmt)
	if !ok {
		t.Fatalf("expected BlockStmt")
	}
	if len(b.statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(b.statements))
	}
}

func TestIfStmt(t *testing.T) {
	stmts := mustParse(t, "if (x) print x;")
	_, ok := stmts[0].(*IfStmt)
	if !ok {
		t.Fatalf("expected IfStmt")
	}
}

func TestIfElseStmt(t *testing.T) {
	stmts := mustParse(t, "if (x) print x; else print 0;")
	i, ok := stmts[0].(*IfStmt)
	if !ok {
		t.Fatalf("expected IfStmt")
	}
	if i.elseBranch == nil {
		t.Fatalf("expected else branch")
	}
}

func TestWhileStmt(t *testing.T) {
	stmts := mustParse(t, "while (x > 0) x = x - 1;")
	_, ok := stmts[0].(*WhileStmt)
	if !ok {
		t.Fatalf("expected WhileStmt")
	}
}

func TestForStmt(t *testing.T) {
	stmts := mustParse(t, "for (var i = 0; i < 10; i = i + 1) print i;")
	_, ok := stmts[0].(*ForStmt)
	if !ok {
		t.Fatalf("expected ForStmt")
	}
}

func TestFunDecl(t *testing.T) {
	stmts := mustParse(t, "fun sumar(a, b) { return a + b; }")
	f, ok := stmts[0].(*FunDecl)
	if !ok {
		t.Fatalf("expected FunDecl")
	}
	if f.name.value != "sumar" {
		t.Fatalf("expected name 'sumar'")
	}
	if len(f.parameters) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(f.parameters))
	}
}

func TestReturnStmt(t *testing.T) {
	stmts := mustParse(t, "fun f() { return 1; }")
	f := stmts[0].(*FunDecl)
	_, ok := f.body[0].(*ReturnStmt)
	if !ok {
		t.Fatalf("expected ReturnStmt")
	}
}

func TestReturnEmpty(t *testing.T) {
	stmts := mustParse(t, "fun f() { return; }")
	f := stmts[0].(*FunDecl)
	r := f.body[0].(*ReturnStmt)
	if r.value != nil {
		t.Fatalf("expected nil return value")
	}
}

// ---------- Casos de error ----------

func TestErrorMissingSemicolon(t *testing.T) {
	mustFail(t, "1 + 2", ";")
}

func TestErrorMissingCloseParen(t *testing.T) {
	mustFail(t, "(1 + 2;", ")")
}

func TestErrorInvalidAssignTarget(t *testing.T) {
	mustFail(t, "1 + 2 = 5;", "invalid assignment target")
}

func TestErrorMissingVarName(t *testing.T) {
	mustFail(t, "var = 5;", "IDENTIFIER")
}

func TestErrorMissingIfParen(t *testing.T) {
	mustFail(t, "if x > 0 print x;", "(")
}

func TestErrorMissingWhileParen(t *testing.T) {
	mustFail(t, "while x > 0 print x;", "(")
}

func TestErrorUnclosedBlock(t *testing.T) {
	mustFail(t, "{ var x = 1;", "}")
}

func TestErrorMissingFunName(t *testing.T) {
	mustFail(t, "fun (a) { }", "IDENTIFIER")
}

func TestErrorMissingFunBrace(t *testing.T) {
	mustFail(t, "fun f(a) return a;", "{")
}

func TestErrorUnexpectedToken(t *testing.T) {
	mustFail(t, "var x = ;", "expression")
}
