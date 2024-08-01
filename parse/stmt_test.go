package parse

import "testing"

func TestDefaultStmt(t *testing.T) {
	tt := []Pair{
		{"default: a;", "(default a)"},
	}

	check(t, tt)
}
func TestCaseStmt(t *testing.T) {
	tt := []Pair{
		{"case 1: a;", "(case 1 a)"},
	}

	check(t, tt)
}
func TestSwitchStmt(t *testing.T) {
	tt := []Pair{
		{"switch (c) { }", "(switch c (block ))"},
		{"switch (c) { case 1: a; case 2: b; default: c; }", "(switch c (block (case 1 a) (case 2 b) (default c)))"},
	}

	check(t, tt)
}

func TestForStmt(t *testing.T) {
	tt := []Pair{
		{"for (;;) {}", "(for (null) (null)  (block ))"},
		{"for (i = 0; i < 10; i++) a;", "(for (i = 0) (i < 10) (i ++) a)"},
		{"for (i = 0;; i++) a;", "(for (i = 0) (null) (i ++) a)"},
		{"for (;i < 10; i++) a;", "(for (null) (i < 10) (i ++) a)"},
		{"for (;;i++) a;", "(for (null) (null) (i ++) a)"},
	}

	check(t, tt)
}
func TestDoStmt(t *testing.T) {
	tt := []Pair{
		{"do { a; } while (b);", "(do b (block a))"},
	}

	check(t, tt)
}

func TestNullStmt(t *testing.T) {
	tt := []Pair{
		{";", "(null)"},
		{"while (true) ;", "(while true (null))"},
		{"{;;}", "(block (null) (null))"},
	}

	check(t, tt)
}

func TestContinueStmt(t *testing.T) {
	tt := []Pair{
		{"continue;", "(continue)"},
		{"while (true) continue;", "(while true (continue))"},
	}

	check(t, tt)
}

func TestBreakStmt(t *testing.T) {
	tt := []Pair{
		{"break;", "(break)"},
		{"while (false) break;", "(while false (break))"},
	}

	check(t, tt)
}

func TestReturnStmt(t *testing.T) {
	tt := []Pair{
		{"return 1;", "(return 1)"},
		{"return a + (b * 2);", "(return (a + (b * 2)))"},
		{"return;", "(return )"},
	}

	check(t, tt)
}

func TestWhileStmt(t *testing.T) {
	tt := []Pair{
		{"while (true) { a; }", "(while true (block a))"},
		{"while (true) { if (a) { b; } }", "(while true (block (if a (block b) )))"},
		{"while (true) a; b;", "(while true a)"},
	}

	check(t, tt)
}

func TestBlockStmt(t *testing.T) {
	tt := []Pair{
		{"{ a; b; }", "(block a b)"},
		{"{ if (a) { b; } }", "(block (if a (block b) ))"},
	}

	check(t, tt)
}

func TestIfStmt(t *testing.T) {
	tt := []Pair{
		{"if (true) { true; }", "(if true (block true) )"},
		{"if (true) { true; } else { false; }", "(if true (block true) (block false))"},
		{"if (a) { if (b) { c; } else { d; }}", "(if a (block (if b (block c) (block d))) )"},
		{"if (a) { b; } else if (c) { d; } else { e; }", "(if a (block b) (if c (block d) (block e)))"},
		{"if (a) {} else if (b) {} else if (c) {} else {}", "(if a (block ) (if b (block ) (if c (block ) (block ))))"},
		{"if (a) b; else c;", "(if a b c)"},
		{"if (a) if (b) c; else d;", "(if a (if b c d) )"},
		{"if (a) b; c;", "(if a b )"},
	}

	check(t, tt)
}
