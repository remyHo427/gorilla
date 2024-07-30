package parse

import (
	"fmt"
	"gorilla/lex"
	"strconv"
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
	ORDER   // < <= > =						(left-to-right)
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
	l    *lex.Lexer
	curr lex.Token
	next lex.Token
	err  []error
}

func New(l *lex.Lexer) *Parser {
	p := &Parser{l: l}
	p.adv()
	p.adv()

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

// stmt
func (p *Parser) parseStmt() Stmt {
	switch p.peek() {
	default:
		expr := p.parseExpr(LOWEST)
		if expr == nil {
			return nil
		}
		p.adv()

		if !p.expect(lex.SCOLON) {
			return nil
		}
		p.adv()

		return &ExprStmt{Expr: expr}
	}
}

// expr
func (p *Parser) parseExpr(currPrec uint) Expr {
	left := p.parsePrefix()

	if left == nil {
		return nil
	}

	for !p.is(lex.SCOLON) && currPrec < p.precn() {
		p.adv()

		switch p.peek() {
		case lex.ADD, lex.SUB, lex.MUL, lex.DIV, lex.MOD, lex.RSHIFT,
			lex.LSHIFT, lex.LT, lex.GT, lex.LEQ, lex.GEQ, lex.EQ,
			lex.NEQ, lex.BAND, lex.BXOR, lex.BOR, lex.AND, lex.OR,
			lex.DOT, lex.COMMA:
			left = p.parseInfixOperator(left)
		case lex.MOD_ASSIGN, lex.LS_ASSIGN, lex.RS_ASSIGN,
			lex.BO_ASSIGN, lex.BA_ASSIGN, lex.XO_ASSIGN,
			lex.DIV_ASSIGN, lex.ADD_ASSIGN, lex.SUB_ASSIGN,
			lex.MUL_ASSIGN, lex.ASSIGN:
			left = p.parseAssign(left)
		case lex.INC, lex.DEC:
			left = p.parsePostfixArithmetic(left)
		// case lex.LPAREN:
		// 	left = p.parseFunctionCall(left)
		// case lex.LBRACKET:
		// 	left = p.parseArrayIndexing(left)
		case lex.QMARK:
			left = p.parseTernaryOperator(left)
		}
	}

	return left
}
func (p *Parser) parsePrefix() Expr {
	switch p.peek() {
	case lex.IDENT:
		return &Ident{Name: p.curr.Literal}
	case lex.INT_CONST:
		n, _ := strconv.Atoi(p.curr.Literal)
		return &Int{Value: int64(n)}
	case lex.ADD, lex.SUB, lex.NOT, lex.INC, lex.DEC,
		lex.BAND, lex.BCOMP:
		return p.parsePrefixOperator()
	case lex.LPAREN:
		p.adv()
		if expr := p.parseExpr(LOWEST); expr == nil {
			return nil
		} else {
			p.adv()
			if !p.expect(lex.RPAREN) {
				return nil
			}
			return expr
		}
	default:
		return nil
	}
}

func (p *Parser) parseAssign(left Expr) Expr {
	expr := &AssignExpr{
		Type: p.peek(),
		Expr: left,
	}
	p.adv()

	if value := p.parseExpr(LOWEST); value == nil {
		return nil
	} else {
		expr.Value = value
	}

	return expr
}

func (p *Parser) parseInfixOperator(left Expr) Expr {
	expr := &InfixExpr{
		Type: p.peek(),
		Left: left,
	}

	currPrec := p.prec()
	p.adv()

	if right := p.parseExpr(currPrec); right == nil {
		return nil
	} else {
		expr.Right = right
	}

	return expr
}

func (p *Parser) parsePrefixOperator() Expr {
	expr := &PrefixExpr{
		Type: p.peek(),
	}
	p.adv()

	if right := p.parseExpr(PREFIX); right == nil {
		return nil
	} else {
		expr.Right = right
	}

	return expr
}

func (p *Parser) parseTernaryOperator(left Expr) Expr {
	expr := &TernaryExpr{
		Cond: left,
	}
	p.adv()

	if Then := p.parseExpr(COND); Then == nil {
		return nil
	} else {
		expr.Then = Then
	}
	p.adv()

	if !p.expect(lex.COLON) {
		return nil
	}
	p.adv()

	if Else := p.parseExpr(COND); Else == nil {
		return nil
	} else {
		expr.Else = Else
	}

	return expr
}

func (p *Parser) parsePostfixArithmetic(left Expr) Expr {
	expr := &PostfixArithmeticExpr{
		Type: p.peek(),
		Left: left,
	}

	return expr
}

func (p *Parser) peek() uint {
	return p.curr.Type
}
func (p *Parser) is(ttype uint) bool {
	return p.curr.Type == ttype
}
func (p *Parser) isn(ttype uint) bool {
	return p.next.Type == ttype
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
		p.error(fmt.Sprintf("expect %s got %s",
			tok_strmap[ttype], tok_strmap[tok.Type]))
		return false
	} else {
		return true
	}
}

func (p *Parser) error(msg string) {
	p.err = append(p.err, fmt.Errorf(msg))
}
