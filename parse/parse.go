package parse

import (
	"fmt"
	"gorilla/lex"
)

const (
	LOWEST  = iota
	COMMA   // , 							(left-to-right)
	ASSIGN  // = += -= etc.					(right-to-left)
	COND    // a ? 1 : 0					(left-to-right)
	OR      // ||							(left-to-right)
	AND     // &&							(left-to-right)
	BOR     // |							(left-to-right)
	BXOR    // ^							(left-to-right)
	BAND    // &							(left-to-right)
	EQ      // == !=						(left-to-right)
	ORDER   // < <= > >=					(left-to-right)
	SHIFT   // << >>						(left-to-right)
	SUM     // + -							(left-to-right)
	PRODUCT // * / %						(left-to-right)
	PREFIX  // ++i --i +i -i !a ~a 			(right-to-left)
	POSTFIX // i++ i-- printf() a[] . -> 	(left-to-right)
)

var prec = map[uint]uint{
	lex.INC:        POSTFIX,
	lex.DEC:        POSTFIX,
	lex.LPAREN:     POSTFIX,
	lex.LBRACKET:   POSTFIX,
	lex.DOT:        POSTFIX,
	lex.ARROW:      POSTFIX,
	lex.MUL:        PRODUCT,
	lex.DIV:        PRODUCT,
	lex.MOD:        PRODUCT,
	lex.ADD:        SUM,
	lex.SUB:        SUM,
	lex.LSHIFT:     SHIFT,
	lex.RSHIFT:     SHIFT,
	lex.GT:         ORDER,
	lex.LT:         ORDER,
	lex.GEQ:        ORDER,
	lex.LEQ:        ORDER,
	lex.EQ:         EQ,
	lex.NEQ:        EQ,
	lex.BAND:       BAND,
	lex.BXOR:       BXOR,
	lex.BOR:        BOR,
	lex.AND:        AND,
	lex.OR:         OR,
	lex.QMARK:      COND,
	lex.ASSIGN:     ASSIGN,
	lex.ADD_ASSIGN: ASSIGN,
	lex.SUB_ASSIGN: ASSIGN,
	lex.MUL_ASSIGN: ASSIGN,
	lex.DIV_ASSIGN: ASSIGN,
	lex.MOD_ASSIGN: ASSIGN,
	lex.RS_ASSIGN:  ASSIGN,
	lex.LS_ASSIGN:  ASSIGN,
	lex.BA_ASSIGN:  ASSIGN,
	lex.XO_ASSIGN:  ASSIGN,
	lex.BO_ASSIGN:  ASSIGN,
	lex.COMMA:      COMMA,
}

type Parser struct {
	l     *lex.Lexer
	curr  lex.Token
	next  lex.Token
	types map[string]bool
	err   []error
}

func New(l *lex.Lexer) *Parser {
	p := &Parser{l: l}

	p.adv()
	p.adv()

	p.types = map[string]bool{
		"bool": true,
	}

	return p
}

func (p *Parser) Parse() ([]Stmt, []error) {
	stmts := []Stmt{}

	for !p.is(lex.EOF) {
		if s := p.parseStmt(); s == nil {
			for !p.is(lex.SCOLON) && !p.is(lex.EOF) {
				p.adv()
			}
			p.adv()
		} else {
			stmts = append(stmts, s)
		}
	}

	if len(p.err) == 0 {
		return stmts, nil
	} else {
		return stmts, p.err
	}
}

func (p *Parser) expectid() bool {
	if p.curr.Type != lex.IDENT {
		p.error("expected IDENT, got %s", toks(p.curr.Type))
		return false
	} else if p.types[p.curr.Literal] {
		p.error("syntax error: using type name %s in place of identifier",
			p.curr.Literal)
		return false
	} else {
		return true
	}
}
func (p *Parser) peek() uint {
	return p.curr.Type
}
func (p *Parser) is(ttype uint) bool {
	return p.curr.Type == ttype
}
func (p *Parser) precn() uint {
	return prec[p.next.Type]
}
func (p *Parser) prec() uint {
	return prec[p.curr.Type]
}
func (p *Parser) adv() {
	p.curr = p.next
	p.next = p.l.Lex()
}
func (p *Parser) expect(ttype uint) bool {
	tok := p.curr
	if tok.Type != ttype {
		p.error("expected %s, got %s", toks(ttype), toks(tok.Type))
		return false
	} else {
		return true
	}
}
func toks(ttype uint) string {
	return lex.Tmap[ttype]
}
func (p *Parser) error(format string, rest ...any) {
	p.err = append(p.err, fmt.Errorf(format, rest...))
}
