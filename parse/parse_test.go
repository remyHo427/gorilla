package parse

import (
	"gorilla/lex"
	"testing"
)

type Pair struct {
	input  string
	output string
}

func TestInfix(t *testing.T) {
	tt := []Pair{
		{"1 + 1;", "(1 + 1)"},
		// todo! expand tests to all infix operators
	}
	check(t, tt)
}

func TestAssign(t *testing.T) {
	tt := []Pair{
		{"a = 1;", "(a = 1)"},
		// todo! expand tests to all assign operators
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
		{"a = b, c;", "((a = b) , c)"},
	}
	check(t, tt)
}

func TestAssociativity(t *testing.T) {
	tt := []Pair{
		{"a++++;", "((a ++) ++)"},
		{"a----;", "((a --) --)"},
		// {"a()()", "((a ()) ())"},
		// {"a[][]", "((a []) [])"},
		{"a.b.c;", "((a . b) . c)"},
		{"++++a;", "(++ (++ a))"},
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
		}

		if test.output != tree[0].String() {
			t.Errorf("expected \"%s\", got \"%s\" at tt[%d]",
				test.output, tree[0].String(), i)
		}
	}
}
