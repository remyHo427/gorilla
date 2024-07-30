package lex

import "testing"

var tmap = map[uint]string{
	EOF:         "EOF",
	ERR:         "ERR",
	INC:         "INC",
	DEC:         "DEC",
	LSHIFT:      "LSHIFT",
	RSHIFT:      "RSHIFT",
	LT:          "LT",
	GT:          "GT",
	LEQ:         "LEQ",
	GEQ:         "GEQ",
	EQ:          "EQ",
	NEQ:         "NEQ",
	BAND:        "BAND",
	BOR:         "BOR",
	MUL:         "MUL",
	DIV:         "DIV",
	MOD:         "MOD",
	ADD:         "ADD",
	SUB:         "SUB",
	AND:         "AND",
	OR:          "OR",
	BXOR:        "XOR",
	NOT:         "NOT",
	ASSIGN:      "ASSIGN",
	MUL_ASSIGN:  "MUL_ASSIGN",
	DIV_ASSIGN:  "DIV_ASSIGN",
	MOD_ASSIGN:  "MOD_ASSIGN",
	ADD_ASSIGN:  "ADD_ASSIGN",
	SUB_ASSIGN:  "SUB_ASSIGN",
	RS_ASSIGN:   "RS_ASSIGN",
	LS_ASSIGN:   "LS_ASSIGN",
	BA_ASSIGN:   "BA_ASSIGN",
	XO_ASSIGN:   "XO_ASSIGN",
	BO_ASSIGN:   "BO_ASSIGN",
	QMARK:       "QMARK",
	SCOLON:      "SCOLON",
	DOT:         "DOT",
	LBRACE:      "LBRACE",
	RBRACE:      "RBRACE",
	LPAREN:      "LPAREN",
	RPAREN:      "RPAREN",
	LBRACKET:    "LBRACKET",
	RBRACKET:    "RBRACCKET",
	COMMA:       "COMMA",
	IDENT:       "IDENT",
	CHAR:        "CHAR",
	STRING:      "STRING",
	INT_CONST:   "INT_CONST",
	FLOAT_CONST: "FLOAT_CONST",
}

func TestEmptyString(t *testing.T) {
	l := New("")
	seq := []uint{
		EOF,
	}
	tokseq(*l, seq, t)
}
func TestSpaceOnly(t *testing.T) {
	l := New("\f\n\r\t\v ")
	seq := []uint{
		EOF,
	}
	tokseq(*l, seq, t)
}
func TestKeywords(t *testing.T) {
	l := New(`auto break case char const continue default do
	double else enum extern float for goto if int long
	register return short signed sizeof static struct
	switch typedef union unsigned void volatile while`)
	seq := []uint{
		AUTO, BREAK, CASE, CHAR, CONST, CONTINUE, DEFAULT, DO,
		DOUBLE, ELSE, ENUM, EXTERN, FLOAT, FOR, GOTO, IF, INT, LONG,
		REGISTER, RETURN, SHORT, SIGNED, SIZEOF, STATIC, STRUCT,
		SWITCH, TYPEDEF, UNION, UNSIGNED, VOID, VOLATILE, WHILE,
	}
	tokseq(*l, seq, t)
}
func TestIdentifiers(t *testing.T) {
	l := New(`a Uint a0 a00 __test__`)
	seq := []uint{
		IDENT, IDENT, IDENT, IDENT,
	}
	tokseq(*l, seq, t)
}
func TestIntegers(t *testing.T) {
	l := New(`0 1 5 100 0100`)
	seq := []uint{
		INT_CONST, INT_CONST, INT_CONST,
		INT_CONST, INT_CONST, INT_CONST, EOF,
	}
	tokseq(*l, seq, t)
}
func TestPunctuators(t *testing.T) {
	l := New(`? ; : . { } ( ) [ ] ,`)
	seq := []uint{
		QMARK, SCOLON, COLON, DOT, LBRACE, RBRACE,
		LPAREN, RPAREN, LBRACKET, RBRACKET, COMMA, EOF,
	}
	tokseq(*l, seq, t)
}
func TestOperators(t *testing.T) {
	l := New(`
		++ += + -- -= - *= * %= % /= / <<= << <= <
		>>= >> >= > && &= & || |= | ^= ^ != ! == =
	`)
	seq := []uint{
		INC, ADD_ASSIGN, ADD, DEC, SUB_ASSIGN, SUB,
		MUL_ASSIGN, MUL, MOD_ASSIGN, MOD, DIV_ASSIGN,
		DIV, LS_ASSIGN, LSHIFT, LEQ, LT, RS_ASSIGN,
		RSHIFT, GEQ, GT, AND, BA_ASSIGN, BAND,
		OR, BO_ASSIGN, BOR, XO_ASSIGN, BXOR, NEQ,
		NOT, EQ, ASSIGN, EOF,
	}
	tokseq(*l, seq, t)
}
func TestComment(t *testing.T) {
	l := New(`
		// a
		// // a
		/*
		b
		*/
		c
	`)
	seq := []uint{
		IDENT, EOF,
	}
	tokseq(*l, seq, t)
}

func tokseq(l Lexer, seq []uint, t *testing.T) {
	for i, ttype := range seq {
		tok := l.Lex()
		if tok.Type != ttype {
			t.Errorf("expected tok type %s, got %s at seq[%d]",
				tmap[ttype], tmap[tok.Type], i)
		}
	}
}
