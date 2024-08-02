package cpp

const (
	EOF = iota
	ERR
	// keywords
	DEFINE
	ELIF
	ELSE
	ERROR
	IF
	IFDEF
	IFNDEF
	INCLUDE
	LINE
	PRAGMA
	UNDEF
	// literal
	HEADER
	IDENT
	PPNUM
	CHAR_CONST
	STRING
	// pp misc
	HASH
	NEWLINE
	WS
	// significant punctuators
	RPAREN
	ELLIP
	COMMA
	MACRO_BEGIN
	// rest of C punctuators
	PUNCT
)

type Token struct {
	Type    uint
	Literal string
}
