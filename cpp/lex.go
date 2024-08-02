package cpp

import (
	"fmt"
	"unicode"
)

var kw_map = map[string]uint{}
var Tmap = map[uint]string{
	EOF:            "EOF",
	DEFINE:         "define",
	ELIF:           "elif",
	ELSE:           "else",
	ERROR:          "error",
	IF:             "if",
	IFDEF:          "ifdef",
	IFNDEF:         "ifndef",
	INCLUDE:        "include",
	LINE:           "line",
	PRAGMA:         "pragma",
	UNDEF:          "undef",
	HEADER:         "HEADER",
	PPNUM:          "PPNUM",
	HASH:           "#",
	NEWLINE:        `\n`,
	NON_NEWLINE_WS: "NON_NEWLINE_WS",
	LBRACKET:       "[",
	RBRACKET:       "]",
	LPAREN:         "(",
	RPAREN:         ")",
	DOT:            ".",
	ARROW:          "->",
	INC:            "++",
	DEC:            "--",
	BAND:           "&",
	MUL:            "*",
	ADD:            "+",
	SUB:            "-",
	BCOMP:          "~",
	NOT:            "!",
	DIV:            "/",
	MOD:            "%",
	LSHIFT:         "<<",
	RSHIFT:         ">>",
	LT:             "<",
	GT:             ">",
	LEQ:            "<=",
	GEQ:            ">=",
	EQ:             "==",
	NEQ:            "!=",
	BXOR:           "^",
	BOR:            "|",
	AND:            "&&",
	OR:             "||",
	QMARK:          "?",
	COLON:          ":",
	ASSIGN:         "=",
	MUL_ASSIGN:     "*=",
	DIV_ASSIGN:     "/=",
	MOD_ASSIGN:     "%=",
	ADD_ASSIGN:     "+=",
	SUB_ASSIGN:     "-=",
	LS_ASSIGN:      "<<=",
	RS_ASSIGN:      ">>=",
	BA_ASSIGN:      "&=",
	XO_ASSIGN:      "^=",
	BO_ASSIGN:      "|=",
	COMMA:          ",",
	LBRACE:         "{",
	RBRACE:         "}",
	SCOLON:         ";",
	ELLIP:          "...",
	IDENT:          "ident",
	CHAR_CONST:     "char_const",
	STRING:         "string",
}

type Lexer struct {
	src     []rune
	sp      int
	keyword map[string]uint
}

func New(src string) *Lexer {
	l := &Lexer{src: []rune(pre(src)), keyword: kw_map}
	return l
}

func (l *Lexer) Lex() Token {
	for !l.isend() {
		c := l.peek()

		var ttype uint
		switch c {
		case '\t', '\r', '\v', '\f', ' ':
			l.adv()
			return tok(NON_NEWLINE_WS)
		case '\n':
			l.adv()
			return tok(NEWLINE)
		case '?':
			l.adv()
			return tok(QMARK)
		case ';':
			l.adv()
			return tok(SCOLON)
		case ':':
			l.adv()
			return tok(COLON)
		case '{':
			l.adv()
			return tok(LBRACE)
		case '}':
			l.adv()
			return tok(RBRACE)
		case '(':
			l.adv()
			return tok(LPAREN)
		case ')':
			l.adv()
			return tok(RPAREN)
		case '[':
			l.adv()
			return tok(LBRACKET)
		case ']':
			l.adv()
			return tok(RBRACKET)
		case ',':
			l.adv()
			return tok(COMMA)
		case '~':
			l.adv()
			return tok(BCOMP)
		case '.':
			l.adv()
			ttype = l.match("..", ELLIP, ttype)
			ttype = l.match("", DOT, ttype)
			return tok(ttype)
		case '+':
			l.adv()
			ttype = l.match("+", INC, ttype)
			ttype = l.match("=", ADD_ASSIGN, ttype)
			ttype = l.match("", ADD, ttype)
			return tok(ttype)
		case '-':
			l.adv()
			ttype = l.match("-", DEC, ttype)
			ttype = l.match("=", SUB_ASSIGN, ttype)
			ttype = l.match(">", ARROW, ttype)
			ttype = l.match("", SUB, ttype)
			return tok(ttype)
		case '*':
			l.adv()
			ttype = l.match("=", MUL_ASSIGN, ttype)
			ttype = l.match("", MUL, ttype)
			return tok(ttype)
		case '%':
			l.adv()
			ttype = l.match("=", MOD_ASSIGN, ttype)
			ttype = l.match("", MOD, ttype)
			return tok(ttype)
		case '/':
			l.adv()
			ttype = l.match("=", DIV_ASSIGN, ttype)
			ttype = l.match("", DIV, ttype)
			return tok(ttype)
		case '<':
			l.adv()
			ttype = l.match("<=", LS_ASSIGN, ttype)
			ttype = l.match("<", LSHIFT, ttype)
			ttype = l.match("=", LEQ, ttype)
			ttype = l.match("", LT, ttype)
			return tok(ttype)
		case '>':
			l.adv()
			ttype = l.match(">=", RS_ASSIGN, ttype)
			ttype = l.match(">", RSHIFT, ttype)
			ttype = l.match("=", GEQ, ttype)
			ttype = l.match("", GT, ttype)
			return tok(ttype)
		case '&':
			l.adv()
			ttype = l.match("&", AND, ttype)
			ttype = l.match("=", BA_ASSIGN, ttype)
			ttype = l.match("", BAND, ttype)
			return tok(ttype)
		case '|':
			l.adv()
			ttype = l.match("|", OR, ttype)
			ttype = l.match("=", BO_ASSIGN, ttype)
			ttype = l.match("", BOR, ttype)
			return tok(ttype)
		case '^':
			l.adv()
			ttype = l.match("=", XO_ASSIGN, ttype)
			ttype = l.match("", BXOR, ttype)
			return tok(ttype)
		case '!':
			l.adv()
			ttype = l.match("=", NEQ, ttype)
			ttype = l.match("", NOT, ttype)
			return tok(ttype)
		case '=':
			l.adv()
			ttype = l.match("=", EQ, ttype)
			ttype = l.match("", ASSIGN, ttype)
			return tok(ttype)
		case '_':
			return l.word()
		case '"':
			return l.string()

		default:
			l.adv()
			return Token{
				Type:    ERR,
				Literal: fmt.Sprintf("unknown character %q", c),
			}
		}
	}

	return tok(EOF)
}

func (l *Lexer) string() Token {
	l.adv()
	start := l.sp

	for {
		if l.isend() {
			break
		} else if l.peek() == '"' {
			break
		} else {
			l.adv()
		}
	}

	if l.isend() {
		return Token{
			Type:    ERR,
			Literal: "unterminated string",
		}
	}

	end := l.sp
	l.adv()

	return Token{
		Type:    STRING,
		Literal: string(l.src[start:end]),
	}
}
func (l *Lexer) word() Token {
	start := l.sp
	for {
		l.adv()
		if l.isend() {
			break
		} else if c := l.peek(); !unicode.IsLetter(c) &&
			!unicode.IsDigit(c) &&
			c != '_' {
			break
		}
	}
	end := l.sp

	s := string(l.src[start:end])
	kword := l.keyword[s]

	if kword != 0 {
		return tok(kword)
	} else {
		return Token{
			Type:    IDENT,
			Literal: s,
		}
	}
}
func (l *Lexer) match(s string, ttype uint, curr uint) uint {
	if curr != EOF {
		return curr
	}

	start := l.sp
	for _, c := range s {
		if l.isend() || l.peek() != rune(c) {
			l.sp = start
			return EOF
		} else {
			l.adv()
		}
	}

	return ttype
}
func (l *Lexer) peek() rune {
	return l.src[l.sp]
}
func (l *Lexer) peekn() rune {
	return l.src[l.sp+1]
}
func (l *Lexer) adv() {
	l.sp++
}
func (l *Lexer) isend() bool {
	return l.sp >= len(l.src)
}

func tok(ttype uint) Token {
	return Token{Type: ttype, Literal: ""}
}
