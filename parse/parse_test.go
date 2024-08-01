package parse

import (
	"gorilla/lex"
	"testing"
)

type Pair struct {
	input  string
	output string
}

func TestEnum(t *testing.T) {
	tt := []Pair{
		{"enum Op { EOF, JMP, POP };", "(decl (enum Op EOF JMP POP))"},
	}

	check(t, tt)
}
func TestTypeSpecifier(t *testing.T) {
	tt := []Pair{
		// will fail if parser isn't initialized with "bool" in p.types
		{"bool;", "(decl (type_specifier bool))"},
	}

	check(t, tt)
}
func TestDefaultTypeSpecifier(t *testing.T) {
	tt := []Pair{
		{"void;", "(decl (default_type_specifier void))"},
		{"char;", "(decl (default_type_specifier char))"},
		{"short;", "(decl (default_type_specifier short))"},
		{"int;", "(decl (default_type_specifier int))"},
		{"long;", "(decl (default_type_specifier long))"},
		{"float;", "(decl (default_type_specifier float))"},
		{"double;", "(decl (default_type_specifier double))"},
		{"signed;", "(decl (default_type_specifier signed))"},
		{"unsigned;", "(decl (default_type_specifier unsigned))"},
		{"unsigned int;", "(decl (default_type_specifier unsigned) (default_type_specifier int))"},
	}

	check(t, tt)
}
func TestTypeQualifer(t *testing.T) {
	tt := []Pair{
		{"const;", "(decl (type_qualifier const))"},
		{"volatile;", "(decl (type_qualifier volatile))"},
		{"const volatile;", "(decl (type_qualifier const) (type_qualifier volatile))"},
	}

	check(t, tt)
}

func TestStorageClass(t *testing.T) {
	tt := []Pair{
		{"typedef;", "(decl (storage_class typedef))"},
		{"extern;", "(decl (storage_class extern))"},
		{"static;", "(decl (storage_class static))"},
		{"auto;", "(decl (storage_class auto))"},
		{"register;", "(decl (storage_class register))"},
		{"extern register;", "(decl (storage_class extern) (storage_class register))"},
	}

	check(t, tt)
}

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
func TestInfix(t *testing.T) {
	tt := []Pair{
		{"1 . 1;", "(1 . 1)"},
		{"1 * 1;", "(1 * 1)"},
		{"1 / 1;", "(1 / 1)"},
		{"1 % 1;", "(1 % 1)"},
		{"1 + 1;", "(1 + 1)"},
		{"1 - 1;", "(1 - 1)"},
		{"1 << 1;", "(1 << 1)"},
		{"1 >> 1;", "(1 >> 1)"},
		{"1 < 1;", "(1 < 1)"},
		{"1 > 1;", "(1 > 1)"},
		{"1 <= 1;", "(1 <= 1)"},
		{"1 >= 1;", "(1 >= 1)"},
		{"1 == 1;", "(1 == 1)"},
		{"1 != 1;", "(1 != 1)"},
		{"1 & 1;", "(1 & 1)"},
		{"1 ^ 1;", "(1 ^ 1)"},
		{"1 | 1;", "(1 | 1)"},
		{"1 && 1;", "(1 && 1)"},
		{"1 || 1;", "(1 || 1)"},
	}
	check(t, tt)
}

func TestAssign(t *testing.T) {
	tt := []Pair{
		{"1 = 2;", "(1 = 2)"},
		{"1 += 2;", "(1 += 2)"},
		{"1 -= 2;", "(1 -= 2)"},
		{"1 /= 2;", "(1 /= 2)"},
		{"1 *= 2;", "(1 *= 2)"},
		{"1 %= 2;", "(1 %= 2)"},
		{"1 <<= 2;", "(1 <<= 2)"},
		{"1 >>= 2;", "(1 >>= 2)"},
		{"1 &= 2;", "(1 &= 2)"},
		{"1 ^= 2;", "(1 ^= 2)"},
		{"1 |= 2;", "(1 |= 2)"},
	}
	check(t, tt)
}

func TestPrefixOpt(t *testing.T) {
	tt := []Pair{
		{"!true;", "(! true)"},
		{"++a;", "(++ a)"},
		{"--a;", "(-- a)"},
		{"+a;", "(+ a)"},
		{"-a;", "(- a)"},
		{"&a;", "(& a)"},
		{"~a;", "(~ a)"},
	}
	check(t, tt)
}

func TestPostfixOpt(t *testing.T) {
	tt := []Pair{
		{"a++;", "(a ++)"},
	}
	check(t, tt)
}

func TestTernary(t *testing.T) {
	tt := []Pair{
		{"1 ? 2 : 3;", "(1 2 3)"},
	}
	check(t, tt)
}

func TestCall(t *testing.T) {
	tt := []Pair{
		{"a();", "(a )"},
		{"a(1);", "(a 1)"},
		{"a(1, 2);", "(a 1 2)"},
		{"a(1, b(2 + 3));", "(a 1 (b (2 + 3)))"},
	}
	check(t, tt)
}

func TestIndexing(t *testing.T) {
	tt := []Pair{
		{"a[i];", "(a i)"},
		{"a[i+1];", "(a (i + 1))"},
	}
	check(t, tt)
}
func TestPrecedence(t *testing.T) {
	tt := []Pair{
		{"++a++;", "(++ (a ++))"},
		{"++a--;", "(++ (a --))"},
		{"++a.b;", "(++ (a . b))"},
		{"++a * b;", "((++ a) * b)"},
		{"--a * b;", "((-- a) * b)"},
		{"!a * b;", "((! a) * b)"},
		{"~a * b;", "((~ a) * b)"},
		{"1 * 2 + 3;", "((1 * 2) + 3)"},
		{"1 / 2 + 3;", "((1 / 2) + 3)"},
		{"1 % 2 + 3;", "((1 % 2) + 3)"},
		{"1 + 2 << 3;", "((1 + 2) << 3)"},
		{"1 - 2 << 3;", "((1 - 2) << 3)"},
		{"1 << 2 < 3;", "((1 << 2) < 3)"},
		{"1 >> 2 < 3;", "((1 >> 2) < 3)"},
		{"1 <= 2 == 3;", "((1 <= 2) == 3)"},
		{"1 >= 2 == 3;", "((1 >= 2) == 3)"},
		{"1 < 2 == 3;", "((1 < 2) == 3)"},
		{"1 > 2 == 3;", "((1 > 2) == 3)"},
		{"1 == 2 & 3;", "((1 == 2) & 3)"},
		{"1 != 2 & 3;", "((1 != 2) & 3)"},
		{"1 & 2 ^ 3;", "((1 & 2) ^ 3)"},
		{"1 ^ 2 | 3;", "((1 ^ 2) | 3)"},
		{"1 && 2 || 3;", "((1 && 2) || 3)"},
		{"1 || 2 ? 3 : 4;", "((1 || 2) 3 4)"},
		{"a = 1 ? 3 : 4;", "(a = (1 3 4))"},
		{"a += 1 ? 3 : 4;", "(a += (1 3 4))"},
		{"a -= 1 ? 3 : 4;", "(a -= (1 3 4))"},
		{"a *= 1 ? 3 : 4;", "(a *= (1 3 4))"},
		{"a /= 1 ? 3 : 4;", "(a /= (1 3 4))"},
		{"a &= 1 ? 3 : 4;", "(a &= (1 3 4))"},
		{"a <<= 1 ? 3 : 4;", "(a <<= (1 3 4))"},
		{"a >>= 1 ? 3 : 4;", "(a >>= (1 3 4))"},
		{"a &= 1 ? 3 : 4;", "(a &= (1 3 4))"},
		{"a ^= 1 ? 3 : 4;", "(a ^= (1 3 4))"},
		{"a |= 1 ? 3 : 4;", "(a |= (1 3 4))"},
	}
	check(t, tt)
}

func TestAssociativity(t *testing.T) {
	tt := []Pair{
		{"a++++;", "((a ++) ++)"},
		{"a----;", "((a --) --)"},
		{"a()();", "((a ) )"},
		{"a[1][1];", "((a 1) 1)"},
		{"a.b.c;", "((a . b) . c)"},
		{"++++a;", "(++ (++ a))"},
		{"----a;", "(-- (-- a))"},
		{"- -a;", "(- (- a))"},
		{"+ +a;", "(+ (+ a))"},
		{"!!a;", "(! (! a))"},
		{"~~a;", "(~ (~ a))"},
		{"1 * 2 * 3;", "((1 * 2) * 3)"},
		{"1 / 2 / 3;", "((1 / 2) / 3)"},
		{"1 % 2 % 3;", "((1 % 2) % 3)"},
		{"1 + 2 + 3;", "((1 + 2) + 3)"},
		{"1 << 2 << 3;", "((1 << 2) << 3)"},
		{"1 >> 2 >> 3;", "((1 >> 2) >> 3)"},
		{"1 < 2 < 3;", "((1 < 2) < 3)"},
		{"1 > 2 > 3;", "((1 > 2) > 3)"},
		{"1 >= 2 >= 3;", "((1 >= 2) >= 3)"},
		{"1 <= 2 <= 3;", "((1 <= 2) <= 3)"},
		{"1 == 2 == 3;", "((1 == 2) == 3)"},
		{"1 != 2 != 3;", "((1 != 2) != 3)"},
		{"1 & 2 & 3;", "((1 & 2) & 3)"},
		{"1 ^ 2 ^ 3;", "((1 ^ 2) ^ 3)"},
		{"1 | 2 | 3;", "((1 | 2) | 3)"},
		{"1 && 2 && 3;", "((1 && 2) && 3)"},
		{"1 || 2 || 3;", "((1 || 2) || 3)"},
		{"1 ? 2 : 3 ? 2 : 3;", "((1 2 3) 2 3)"},
		{"1 = 2 = 3;", "((1 = 2) = 3)"},
		{"1 += 2 += 3;", "((1 += 2) += 3)"},
		{"1 -= 2 -= 3;", "((1 -= 2) -= 3)"},
		{"1 *= 2 *= 3;", "((1 *= 2) *= 3)"},
		{"1 /= 2 /= 3;", "((1 /= 2) /= 3)"},
		{"1 %= 2 %= 3;", "((1 %= 2) %= 3)"},
		{"1 <<= 2 <<= 3;", "((1 <<= 2) <<= 3)"},
		{"1 >>= 2 >>= 3;", "((1 >>= 2) >>= 3)"},
		{"1 &= 2 &= 3;", "((1 &= 2) &= 3)"},
		{"1 ^= 2 ^= 3;", "((1 ^= 2) ^= 3)"},
		{"1 |= 2 |= 3;", "((1 |= 2) |= 3)"},
	}
	check(t, tt)
}

func check(t *testing.T, tt []Pair) {
	for i, test := range tt {
		l := lex.New(test.input)
		p := New(l)
		tree, err := p.Parse()

		for _, e := range err {
			t.Errorf(e.Error())
		}

		if len(p.err) > 0 {
			continue
		} else if len(tree) == 0 {
			t.Errorf("no ast produced")
			continue
		}

		if test.output != tree[0].String() {
			t.Errorf("expected \"%s\", got \"%s\" at tt[%d]",
				test.output, tree[0].String(), i)
		}
	}
}
