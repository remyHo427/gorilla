package parse

import (
	"gorilla/lex"
	"strconv"
)

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
			lex.DOT:
			left = p.parseInfixOperator(left)
		case lex.MOD_ASSIGN, lex.LS_ASSIGN, lex.RS_ASSIGN,
			lex.BO_ASSIGN, lex.BA_ASSIGN, lex.XO_ASSIGN,
			lex.DIV_ASSIGN, lex.ADD_ASSIGN, lex.SUB_ASSIGN,
			lex.MUL_ASSIGN, lex.ASSIGN:
			left = p.parseAssign(left)
		case lex.INC, lex.DEC:
			left = p.parsePostfixArithmetic(left)
		case lex.LPAREN:
			left = p.parseFunctionCall(left)
		case lex.LBRACKET:
			left = p.parseArrayIndexing(left)
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

	if value := p.parseExpr(ASSIGN); value == nil {
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

func (p *Parser) parseFunctionCall(left Expr) Expr {
	expr := &CallExpr{
		Callee: left,
	}
	p.adv()

	args := []Expr{}
	for !p.is(lex.RPAREN) && !p.is(lex.EOF) {
		arg := p.parseExpr(LOWEST)

		if arg == nil {
			return nil
		} else {
			args = append(args, arg)
		}

		if !p.is(lex.COMMA) {
			p.adv()
			break
		} else {
			p.adv()
		}
	}
	expr.Args = args

	if !p.expect(lex.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseArrayIndexing(left Expr) Expr {
	expr := &IndexExpr{
		Arr: left,
	}
	p.adv()

	if idx := p.parseExpr(LOWEST); idx == nil {
		p.error("expression expected")
		return nil
	} else {
		expr.Index = idx
		p.adv()
	}

	if !p.expect(lex.RBRACKET) {
		return nil
	}

	return expr
}
