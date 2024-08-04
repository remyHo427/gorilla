package cpp

import "testing"

func TestPpnum(t *testing.T) {
	l := New("0 1 .1 0.1 0.1f 1 1 100ul .1ab.+-..-+")
	seq := []uint{PPNUM, WS, PPNUM, WS, PPNUM, WS, PPNUM,
		WS, PPNUM, WS, PPNUM, WS, PPNUM, WS, PPNUM, WS, PPNUM, EOF}

	tokseq(*l, seq, t)
}

func TestNewline(t *testing.T) {
	l := New("\n\n\n")
	seq := []uint{NEWLINE, NEWLINE, NEWLINE, EOF}

	tokseq(*l, seq, t)
}

func TestNonNewlineWS(t *testing.T) {
	l := New("\t\r\v \n \t\t\n")
	seq := []uint{WS, NEWLINE, WS, NEWLINE, EOF}

	tokseq(*l, seq, t)
}

func TestKeywords(t *testing.T) {
	l := New(`define elif else error if ifdef ifndef
		include line pragma undef`)
	seq := []uint{DEFINE, WS, ELIF, WS, ELSE, WS, ERROR, WS, IF, WS,
		IFDEF, WS, IFNDEF, NEWLINE, WS, INCLUDE, WS, LINE, WS,
		PRAGMA, WS, UNDEF}

	tokseq(*l, seq, t)
}

func TestInsignificantPunctuators(t *testing.T) {
	l := New(`. ? ; : { } ( [ ] ~ ++ += + -- -= -> - *= * %= % /= / <<= << <= < >>= >> >= > && &= & || |= | ^= ^ != ! = `)
	seq := []string{
		".", "?", ";", ":", "{", "}", "(", "[", "]", "~", "++",
		"+=", "+", "--", "-=", "->", "-", "*=", "*", "%=", "%",
		"/=", "/", "<<=", "<<", "<=", "<", ">>=", ">>", ">=", ">",
		"&&", "&=", "&", "||", "|=", "|", "^=", "^", "!=", "!", "=",
	}

	for _, s := range seq {
		if tok := l.Lex(); tok.Type != PUNCT {
			t.Errorf("expects a punctuator")
		} else if tok.Literal != s {
			t.Errorf(`want="%s", got="%s"`, s, tok.Literal)
		} else if l.Lex().Type != WS {
			t.Errorf(`following WS not parsed`)
		}
	}

	l = New(`=`)
	if tok := l.Lex(); tok.Type != PUNCT {
		t.Errorf("did not match a single punctuator, got: %s", Tmap[tok.Type])
	} else if tok.Literal != "=" {
		t.Errorf(`expect "=", got="%s"`, tok.Literal)
	}
}

func TestSignificantPunctuators(t *testing.T) {
	l := New(`, ... # )`)
	seq := []uint{
		COMMA, WS, ELLIP, WS, HASH, WS,
		RPAREN,
	}

	tokseq(*l, seq, t)
}

func tokseq(l Lexer, seq []uint, t *testing.T) {
	for i, ttype := range seq {
		tok := l.Lex()
		if tok.Type != ttype {
			t.Errorf("expected tok type %s, got %s at seq[%d] (%s)",
				Tmap[ttype], Tmap[tok.Type], i, tok.Literal)
		}
	}
}
