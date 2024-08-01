package lex

import (
	"fmt"
	"unicode"
)

var Tmap = map[uint]string{
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
	BOR:         "|",
	AND:         "&&",
	OR:          "||",
	QMARK:       "?",
	COLON:       ":",
	ASSIGN:      "=",
	MUL_ASSIGN:  "*=",
	DIV_ASSIGN:  "/=",
	MOD_ASSIGN:  "%=",
	ADD_ASSIGN:  "+=",
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

var kw_map = map[string]uint{
	"auto":     AUTO,
	"break":    BREAK,
	"case":     CASE,
	"char":     CHAR,
	"const":    CONST,
	"continue": CONTINUE,
	"default":  DEFAULT,
	"do":       DO,
	"double":   DOUBLE,
	"else":     ELSE,
	"enum":     ENUM,
	"extern":   EXTERN,
	"float":    FLOAT,
	"for":      FOR,
	"goto":     GOTO,
	"if":       IF,
	"int":      INT,
	"long":     LONG,
	"register": REGISTER,
	"return":   RETURN,
	"short":    SHORT,
	"signed":   SIGNED,
	"sizeof":   SIZEOF,
	"static":   STATIC,
	"struct":   STRUCT,
	"switch":   SWITCH,
	"typedef":  TYPEDEF,
	"union":    UNION,
	"unsigned": UNSIGNED,
	"void":     VOID,
	"volatile": VOLATILE,
	"while":    WHILE,
}

type Lexer struct {
	src   []rune
	sp    int
	kword map[string]uint
}

func New(src string) *Lexer {
	l := &Lexer{src: []rune(src), kword: kw_map}
	return l
}
func (l *Lexer) Lex() Token {
	for !l.isend() {
		c := l.peek()

		if unicode.IsSpace(c) {
			l.adv()
			continue
		}
		if unicode.IsLetter(c) {
			return l.word()
		}
		if unicode.IsDigit(c) {
			return l.number()
		}

		var ttype uint = EOF
		switch c {
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
			if l.isend() {
				return tok(DIV)
			} else if c := l.peek(); c == '/' || c == '*' {
				l.skip_comment()
				continue
			} else {
				ttype = l.match("=", DIV_ASSIGN, ttype)
				ttype = l.match("", DIV, ttype)
				return tok(ttype)
			}
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
func (l *Lexer) skip_comment() {
	if c := l.peek(); c == '/' {
		for {
			l.adv()
			if l.isend() || l.peek() == '\n' {
				break
			}
		}
	} else if c == '*' {
		for {
			l.adv()
			if l.isend() {
				break
			} else if l.peek() == '*' && l.peekn() == '/' {
				l.adv()
				l.adv()
				break
			}
		}
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
	kword := l.kword[s]

	if kword != 0 {
		return tok(kword)
	} else {
		return Token{
			Type:    IDENT,
			Literal: s,
		}
	}
}
func (l *Lexer) number() Token {
	if l.peek() == '0' {
		l.adv()
		return Token{
			Type:    INT_CONST,
			Literal: "0",
		}
	}

	start := l.sp
	for {
		l.adv()
		if l.isend() {
			break
		} else if c := l.peek(); !unicode.IsDigit(c) {
			break
		}
	}
	end := l.sp
	s := string(l.src[start:end])

	return Token{
		Type:    INT_CONST,
		Literal: s,
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
