package lex

import "testing"

var tmap = map[uint]string{
	EOF:         "EOF",
	ERR:         "ERR",
	AUTO:        "auto",
	BREAK:       "break",
	CASE:        "case",
	CHAR:        "char",
	CONST:       "const",
	CONTINUE:    "continue",
	DEFAULT:     "default",
	DO:          "do",
	DOUBLE:      "double",
	ELSE:        "else",
	ENUM:        "enum",
	EXTERN:      "extern",
	FLOAT:       "float",
	FOR:         "for",
	GOTO:        "goto",
	IF:          "if",
	INT:         "int",
	LONG:        "long",
	REGISTER:    "register",
	RETURN:      "return",
	SHORT:       "short",
	SIGNED:      "signed",
	SIZEOF:      "sizeof",
	STATIC:      "static",
	STRUCT:      "struct",
	SWITCH:      "switch",
	TYPEDEF:     "typedef",
	UNION:       "union",
	UNSIGNED:    "unsigned",
	VOID:        "void",
	VOLATILE:    "volatile",
	WHILE:       "while",
	LBRACKET:    "[",
	RBRACKET:    "]",
	LPAREN:      "(",
	RPAREN:      ")",
	DOT:         ".",
	ARROW:       "->",
	INC:         "++",
	DEC:         "--",
	BAND:        "&",
	MUL:         "*",
	ADD:         "+",
	SUB:         "-",
	BCOMP:       "~",
	NOT:         "!",
	DIV:         "/",
	MOD:         "%",
	LSHIFT:      "<<",
	RSHIFT:      ">>",
	LT:          "<",
	GT:          ">",
	LEQ:         "<=",
	GEQ:         ">=",
	EQ:          "==",
	NEQ:         "!=",
	BXOR:        "^",
	AND:         "&&",
	OR:          "||",
	QMARK:       "?",
	COLON:       ":",
	ASSIGN:      "=",
	MUL_ASSIGN:  "*=",
	DIV_ASSIGN:  "/=",
	MOD_ASSIGN:  "%=",
	SUB_ASSIGN:  "-=",
	LS_ASSIGN:   "<<=",
	RS_ASSIGN:   ">>=",
	BA_ASSIGN:   "&=",
	XO_ASSIGN:   "^=",
	BO_ASSIGN:   "|=",
	COMMA:       ",",
	LBRACE:      "{",
	RBRACE:      "}",
	SCOLON:      ";",
	ELLIP:       "...",
	IDENT:       "ident",
	CHAR_CONST:  "char_const",
	STRING:      "string",
	INT_CONST:   "int_const",
	FLOAT_CONST: "float_const",
}

func TestString(t *testing.T) {
	l := New(`"test"`)
	if l.Lex().Literal != "test" {
		t.Errorf("lex string failed\n")
	}
	l = New(`""`)
	if l.Lex().Literal != "" {
		t.Errorf("lex string failed\n")
	}
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
func TestOperators(t *testing.T) {
	l := New(`
		[ ] ( ) . -> ++ -- & * + - ~ ! / % << >> < > <=
		>= == != ^ && || ? : = *= /= %= -= <<= >>= &= ^=
		|= , { } ; ...
	`)
	seq := []uint{
		LBRACKET, RBRACKET, LPAREN, RPAREN, DOT, ARROW, INC, DEC,
		BAND, MUL, ADD, SUB, BCOMP, NOT, DIV, MOD, LSHIFT, RSHIFT,
		LT, GT, LEQ, GEQ, EQ, NEQ, BXOR, AND, OR, QMARK, COLON,
		ASSIGN, MUL_ASSIGN, DIV_ASSIGN, MOD_ASSIGN, SUB_ASSIGN,
		LS_ASSIGN, RS_ASSIGN, BA_ASSIGN, XO_ASSIGN, BO_ASSIGN,
		COMMA, LBRACE, RBRACE, SCOLON, ELLIP,
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
