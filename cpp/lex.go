package cpp

import (
	"fmt"
	"unicode"
)

var kw_map = map[string]uint{
	"define":  DEFINE,
	"elif":    ELIF,
	"else":    ELSE,
	"error":   ERROR,
	"if":      IF,
	"ifdef":   IFDEF,
	"ifndef":  IFNDEF,
	"include": INCLUDE,
	"line":    LINE,
	"pragma":  PRAGMA,
	"undef":   UNDEF,
}

var Tmap = map[uint]string{
	EOF:        "EOF",
	DEFINE:     "define",
	ELIF:       "elif",
	ELSE:       "else",
	ERROR:      "error",
	IF:         "if",
	IFDEF:      "ifdef",
	IFNDEF:     "ifndef",
	INCLUDE:    "include",
	LINE:       "line",
	PRAGMA:     "pragma",
	UNDEF:      "undef",
	HEADER:     "HEADER",
	PPNUM:      "PPNUM",
	HASH:       "#",
	NEWLINE:    `\n`,
	WS:         "WS",
	COMMA:      ",",
	PUNCT:      "PUNCT",
	ELLIP:      "...",
	IDENT:      "ident",
	CHAR_CONST: "char_const",
	STRING:     "string",
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

		if unicode.IsLetter(c) {
			return l.word()
		}

		var ttype uint
		switch c {
		// phase 3.3: newlines are kept
		case '\n':
			l.adv()
			return tok(NEWLINE)
		case '#':
			l.adv()
			return tok(HASH)
		case ',':
			l.adv()
			return tok(COMMA)
		case ')':
			l.adv()
			return tok(RPAREN)
		case '.':
			l.adv()
			ttype = l.match("..", ELLIP, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '?', ';', ':', '{', '}', '(', '[', ']', '~':
			l.adv()
			return tok(PUNCT)
		case '+':
			l.adv()
			ttype = l.match("+", PUNCT, ttype)
			ttype = l.match("=", PUNCT, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '-':
			l.adv()
			ttype = l.match("-", PUNCT, ttype)
			ttype = l.match("=", PUNCT, ttype)
			ttype = l.match(">", PUNCT, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '*':
			l.adv()
			ttype = l.match("=", PUNCT, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '%':
			l.adv()
			ttype = l.match("=", PUNCT, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '/':
			l.adv()
			ttype = l.match("=", PUNCT, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '<':
			l.adv()
			ttype = l.match("<=", PUNCT, ttype)
			ttype = l.match("<", PUNCT, ttype)
			ttype = l.match("=", PUNCT, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '>':
			l.adv()
			ttype = l.match(">=", PUNCT, ttype)
			ttype = l.match(">", PUNCT, ttype)
			ttype = l.match("=", PUNCT, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '&':
			l.adv()
			ttype = l.match("&", PUNCT, ttype)
			ttype = l.match("=", PUNCT, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '|':
			l.adv()
			ttype = l.match("|", PUNCT, ttype)
			ttype = l.match("=", PUNCT, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '^':
			l.adv()
			ttype = l.match("=", PUNCT, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '!':
			l.adv()
			ttype = l.match("=", PUNCT, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '=':
			l.adv()
			ttype = l.match("=", PUNCT, ttype)
			ttype = l.match("", PUNCT, ttype)
			return tok(ttype)
		case '_':
			return l.word()
		case '"':
			return l.string()

		default:
			// phase 3.2: non newline ws are collapsed into one space character
			if unicode.IsSpace(c) {
				l.group_ws()
				return tok(WS)
			} else {
				l.adv()
				return Token{
					Type:    ERR,
					Literal: fmt.Sprintf("unknown character %q", c),
				}
			}
		}
	}

	return tok(EOF)
}

func (l *Lexer) group_ws() {
	for {
		l.adv()
		if c := l.peek(); !unicode.IsSpace(c) || c == '\n' {
			break
		}
	}
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

	if kword, ok := l.keyword[s]; ok {
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
func (l *Lexer) adv() {
	l.sp++
}
func (l *Lexer) isend() bool {
	return l.sp >= len(l.src)
}

func tok(ttype uint) Token {
	return Token{Type: ttype, Literal: ""}
}
