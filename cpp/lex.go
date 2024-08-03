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

		switch c {
		// phase 3.3: newlines are kept
		case '\n':
			l.adv()
			return tok(NEWLINE, "\n")
		case '#':
			l.adv()
			return tok(HASH, "#")
		case ',':
			l.adv()
			return tok(COMMA, ",")
		case ')':
			l.adv()
			return tok(RPAREN, "(")
		case '.':
			if p, matched := l.oneOf("..."); matched {
				return tok(ELLIP, p)
			} else {
				l.adv()
				return tok(PUNCT, ".")
			}
		case '?', ';', ':', '{', '}', '(', '[', ']', '~':
			l.adv()
			return tok(PUNCT, string(c))
		case '+':
			if p, matched := l.oneOf("++", "+="); matched {
				return tok(PUNCT, p)
			} else {
				l.adv()
				return tok(PUNCT, "+")
			}
		case '-':
			if p, matched := l.oneOf("--", "-=", "->"); matched {
				return tok(PUNCT, p)
			} else {
				l.adv()
				return tok(PUNCT, "-")
			}
		case '*':
			if p, matched := l.oneOf("*="); matched {
				return tok(PUNCT, p)
			} else {
				l.adv()
				return tok(PUNCT, "*")
			}
		case '%':
			if p, matched := l.oneOf("%="); matched {
				return tok(PUNCT, p)
			} else {
				l.adv()
				return tok(PUNCT, "%")
			}
		case '/':
			if p, matched := l.oneOf("/="); matched {
				return tok(PUNCT, p)
			} else {
				l.adv()
				return tok(PUNCT, "/")
			}
		case '<':
			if p, matched := l.oneOf("<<=", "<<", "<="); matched {
				return tok(PUNCT, p)
			} else {
				l.adv()
				return tok(PUNCT, "<")
			}
		case '>':
			if p, matched := l.oneOf(">>=", ">>", ">="); matched {
				return tok(PUNCT, p)
			} else {
				l.adv()
				return tok(PUNCT, ">")
			}
		case '&':
			if p, matched := l.oneOf("&=", "&&"); matched {
				return tok(PUNCT, p)
			} else {
				l.adv()
				return tok(PUNCT, "&")
			}
		case '|':
			if p, matched := l.oneOf("|=", "||"); matched {
				return tok(PUNCT, p)
			} else {
				l.adv()
				return tok(PUNCT, "|")
			}
		case '^':
			if p, matched := l.oneOf("^="); matched {
				return tok(PUNCT, p)
			} else {
				l.adv()
				return tok(PUNCT, "^")
			}
		case '!':
			if p, matched := l.oneOf("!="); matched {
				return tok(PUNCT, p)
			} else {
				l.adv()
				return tok(PUNCT, "!")
			}
		case '=':
			if p, matched := l.oneOf("=="); matched {
				return tok(PUNCT, p)
			} else {
				l.adv()
				return tok(PUNCT, "=")
			}
		case '_':
			return l.word()
		case '"':
			return l.string()

		default:
			// phase 3.2: non newline ws are collapsed into one space character
			if unicode.IsSpace(c) {
				return tok(WS, l.group_ws())
			} else {
				l.adv()
				return Token{
					Type:    ERR,
					Literal: fmt.Sprintf("unknown character %q", c),
				}
			}
		}
	}

	return tok(EOF, "")
}

func (l *Lexer) group_ws() string {
	start := l.sp

	for {
		l.adv()
		if l.isend() {
			break
		} else if c := l.peek(); !unicode.IsSpace(c) || c == '\n' {
			break
		}
	}

	end := l.sp
	return string(l.src[start:end])
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
		return tok(kword, s)
	} else {
		return Token{
			Type:    IDENT,
			Literal: s,
		}
	}
}
func (l *Lexer) oneOf(patterns ...string) (string, bool) {
	start := l.sp

	for _, p := range patterns {
		l.sp = start
		i := 0

		for i < len(p) {
			if l.isend() {
				return "", false
			} else if rune(p[i]) != l.peek() {
				break
			} else {
				l.adv()
				i++
			}
		}

		if i >= len(p) {
			return p, true
		}
	}

	l.sp = start
	return "", false
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

func tok(ttype uint, s string) Token {
	return Token{Type: ttype, Literal: s}
}
