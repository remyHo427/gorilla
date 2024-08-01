package parse

import "gorilla/lex"

// stmt
func (p *Parser) parseStmt() Stmt {
	switch p.peek() {
	case lex.IF:
		return p.parseIfStmt()
	case lex.WHILE:
		return p.parseWhileStmt()
	case lex.RETURN:
		return p.parseReturnStmt()
	case lex.BREAK:
		return p.parseBreakStmt()
	case lex.CONTINUE:
		return p.parseContinueStmt()
	case lex.LBRACE:
		return p.parseBlockStmt()
	case lex.DO:
		return p.parseDoStmt()
	case lex.FOR:
		return p.parseForStmt()
	case lex.SWITCH:
		return p.parseSwitchStmt()
	case lex.CASE:
		return p.parseCaseStmt()
	case lex.DEFAULT:
		return p.parseDefaultStmt()
	case lex.SCOLON:
		p.adv()
		return &NullStmt{}
	default:
		if decl := p.parseDeclStmt(); decl != nil {
			return decl
		} else if expr := p.parseExprStmt(); expr != nil {
			return expr
		} else {
			return nil
		}
	}
}

func (p *Parser) parseIfStmt() Stmt {
	stmt := &IfStmt{}
	p.adv()

	if !p.expect(lex.LPAREN) {
		return nil
	}
	p.adv()

	if If := p.parseExpr(LOWEST); If == nil {
		return nil
	} else {
		stmt.If = If
	}
	p.adv()

	if !p.expect(lex.RPAREN) {
		return nil
	}
	p.adv()

	if Then := p.parseStmt(); Then == nil {
		return nil
	} else {
		stmt.Then = Then
	}

	if !p.is(lex.ELSE) {
		return stmt
	}
	p.adv()

	if Else := p.parseStmt(); Else == nil {
		return nil
	} else {
		stmt.Else = Else
	}

	return stmt
}

func (p *Parser) parseWhileStmt() Stmt {
	stmt := &WhileStmt{}
	p.adv()

	if !p.expect(lex.LPAREN) {
		return nil
	}
	p.adv()

	if cond := p.parseExpr(LOWEST); cond == nil {
		return nil
	} else {
		stmt.Cond = cond
	}
	p.adv()

	if !p.expect(lex.RPAREN) {
		return nil
	}
	p.adv()

	if loop := p.parseStmt(); loop == nil {
		return nil
	} else {
		stmt.Loop = loop
	}

	return stmt
}

func (p *Parser) parseReturnStmt() Stmt {
	stmt := &ReturnStmt{}
	p.adv()

	if p.is(lex.SCOLON) {
		p.adv()
		return stmt
	}

	if rtrn := p.parseExpr(LOWEST); rtrn == nil {
		return nil
	} else {
		stmt.Return = rtrn
	}
	p.adv()

	if !p.expect(lex.SCOLON) {
		return nil
	}
	p.adv()

	return stmt
}

func (p *Parser) parseBreakStmt() Stmt {
	stmt := &BreakStmt{}
	p.adv()

	if !p.expect(lex.SCOLON) {
		return nil
	}
	p.adv()

	return stmt
}

func (p *Parser) parseContinueStmt() Stmt {
	stmt := &ContinueStmt{}
	p.adv()

	if !p.expect(lex.SCOLON) {
		return nil
	}
	p.adv()

	return stmt
}

func (p *Parser) parseDoStmt() Stmt {
	stmt := &DoStmt{}
	p.adv()

	if loop := p.parseStmt(); loop == nil {
		return nil
	} else {
		stmt.Loop = loop
	}

	if !p.expect(lex.WHILE) {
		return nil
	}
	p.adv()

	if !p.expect(lex.LPAREN) {
		return nil
	}
	p.adv()

	if cond := p.parseExpr(LOWEST); cond == nil {
		return nil
	} else {
		stmt.Cond = cond
	}
	p.adv()

	if !p.expect(lex.RPAREN) {
		return nil
	}
	p.adv()

	if !p.expect(lex.SCOLON) {
		return nil
	}
	p.adv()

	return stmt
}

func (p *Parser) parseForStmt() Stmt {
	stmt := &ForStmt{}
	p.adv()

	if !p.expect(lex.LPAREN) {
		return nil
	}
	p.adv()

	if init := p.parseExprStmt(); init == nil {
		return nil
	} else {
		stmt.Init = init
	}

	if cond := p.parseExprStmt(); cond == nil {
		return nil
	} else {
		stmt.Cond = cond
	}

	if post := p.parseExpr(LOWEST); post != nil {
		p.adv()
		stmt.Post = post
	}

	if !p.expect(lex.RPAREN) {
		return nil
	}
	p.adv()

	if loop := p.parseStmt(); loop == nil {
		return nil
	} else {
		stmt.Loop = loop
	}

	return stmt
}

func (p *Parser) parseBlockStmt() Stmt {
	block := &BlockStmt{}
	p.adv()

	stmts := []Stmt{}
	for !p.is(lex.RBRACE) && !p.is(lex.EOF) {
		if s := p.parseStmt(); s == nil {
			return nil
		} else {
			stmts = append(stmts, s)
		}
	}
	block.Stmts = stmts

	if !p.expect(lex.RBRACE) {
		return nil
	}
	p.adv()

	return block
}

func (p *Parser) parseExprStmt() Stmt {
	// required because when parsing inside for (...) there
	// might be null statements
	if p.is(lex.SCOLON) {
		p.adv()
		return &NullStmt{}
	}

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

func (p *Parser) parseSwitchStmt() Stmt {
	stmt := &SwitchStmt{}
	p.adv()

	if !p.expect(lex.LPAREN) {
		return nil
	}
	p.adv()

	if expr := p.parseExpr(LOWEST); expr == nil {
		return nil
	} else {
		stmt.Cond = expr
	}
	p.adv()

	if !p.expect(lex.RPAREN) {
		return nil
	}
	p.adv()

	if body := p.parseStmt(); body == nil {
		return nil
	} else {
		stmt.Stmt = body
	}

	return stmt
}

func (p *Parser) parseCaseStmt() Stmt {
	stmt := &CaseStmt{}
	p.adv()

	if expr := p.parseExpr(LOWEST); expr == nil {
		return nil
	} else {
		stmt.Cond = expr
	}
	p.adv()

	if !p.expect(lex.COLON) {
		return nil
	}
	p.adv()

	if s := p.parseStmt(); s == nil {
		return nil
	} else {
		stmt.Stmt = s
	}

	return stmt
}

func (p *Parser) parseDefaultStmt() Stmt {
	stmt := &DefaultStmt{}
	p.adv()

	if !p.expect(lex.COLON) {
		return nil
	}
	p.adv()

	if s := p.parseStmt(); s == nil {
		return nil
	} else {
		stmt.Stmt = s
	}

	return stmt
}
