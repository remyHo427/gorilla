package cpp

import (
	"bytes"
	"fmt"
)

type Parser struct {
	l      *Lexer
	curr   Token
	next   Token
	macros map[string]string
	err    []error
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}

	p.adv()
	p.adv()

	return p
}

func (p *Parser) Expand() (string, []error) {
	var out bytes.Buffer

	for !p.is(EOF) {
		if s, ok := p.parseGroup(); !ok {
			for !p.is(NEWLINE) && !p.is(EOF) {
				p.adv()
			}
			p.adv()
		} else {
			out.WriteString(s)
		}
	}

	if len(p.err) != 0 {
		return "", p.err
	} else {
		return out.String(), nil
	}
}

func (p *Parser) parseGroup() (string, bool) {
	if !p.is(HASH) {
		p.adv()
		return p.parseTextline()
	}

	return "", true
}

func (p *Parser) parseTextline() (string, bool) {
	var out bytes.Buffer

	for !p.is(NEWLINE) {
	}

	return out.String(), true
}

func (p *Parser) peek() uint {
	return p.curr.Type
}
func (p *Parser) is(ttype uint) bool {
	return p.curr.Type == ttype
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
	return Tmap[ttype]
}
func (p *Parser) error(format string, rest ...any) {
	p.err = append(p.err, fmt.Errorf(format, rest...))
}
