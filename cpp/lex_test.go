package cpp

import "testing"

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
	l := New(`. ? ; : { } ( [ ] ~ ++ += + -- -= -> - *=
* %= % /= / <<= << <= < >>= >> >= > && &= & || |= | ^= ^ != ! =`)
	seq := []uint{
		PUNCT, WS, PUNCT, WS,
		PUNCT, WS, PUNCT, WS, PUNCT, WS,
		PUNCT, WS, PUNCT, WS, PUNCT, WS,
		PUNCT, WS, PUNCT, WS, PUNCT, WS,
		PUNCT, WS, PUNCT, WS, PUNCT, WS,
		PUNCT, WS, PUNCT, WS, PUNCT, WS,
		PUNCT, NEWLINE, PUNCT, WS,
		PUNCT, WS, PUNCT, WS, PUNCT, WS,
		PUNCT, WS, PUNCT, WS, PUNCT, WS,
		PUNCT, WS, PUNCT, WS, PUNCT, WS,
		PUNCT, WS, PUNCT, WS, PUNCT, WS,
		PUNCT, WS, PUNCT, WS, PUNCT, WS,
		PUNCT, WS, PUNCT, WS, PUNCT, WS,
		PUNCT, WS, PUNCT, WS, PUNCT, WS,
		PUNCT, WS, PUNCT,
	}

	tokseq(*l, seq, t)
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
			t.Errorf("expected tok type %s, got %s at seq[%d]",
				Tmap[ttype], Tmap[tok.Type], i)
		}
	}
}
