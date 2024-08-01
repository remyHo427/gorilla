package parse

import "gorilla/lex"

func (p *Parser) parseDeclStmt() Stmt {
	stmt := &DeclStmt{}

	decls := []Decl{}
	for !p.is(lex.SCOLON) && !p.is(lex.EOF) {
		var decl Decl

		switch ttype := p.peek(); ttype {
		case lex.TYPEDEF, lex.EXTERN, lex.STATIC, lex.AUTO, lex.REGISTER:
			decl = p.parseStorageClass(ttype)
		case lex.CONST, lex.VOLATILE:
			decl = p.parseTypeQualifier(ttype)
		case lex.VOID, lex.CHAR, lex.SHORT, lex.INT, lex.LONG,
			lex.FLOAT, lex.DOUBLE, lex.SIGNED, lex.UNSIGNED:
			decl = p.parseDefaultTypeSpecifier(ttype)
		case lex.ENUM:
			decl = p.parseEnum()
		case lex.IDENT:
			decl = p.parseTypeSpecifier()
		default:
			decl = nil
		}

		if decl == nil {
			return nil
		} else {
			decls = append(decls, decl)
		}
	}

	if !p.expect(lex.SCOLON) {
		return nil
	}
	p.adv()

	stmt.Decls = decls
	return stmt
}

func (p *Parser) parseStorageClass(ttype uint) Decl {
	p.adv()
	return &StorageClass{
		Type: ttype,
	}
}

func (p *Parser) parseTypeQualifier(ttype uint) Decl {
	p.adv()
	return &TypeQualifer{
		Type: ttype,
	}
}

func (p *Parser) parseDefaultTypeSpecifier(ttype uint) Decl {
	p.adv()
	return &DefaultTypeSpecifier{
		Type: ttype,
	}
}

func (p *Parser) parseTypeSpecifier() Decl {
	if id := p.curr.Literal; p.types[id] {
		p.adv()
		return &TypeSpecifier{
			Literal: id,
		}
	} else {
		return nil
	}
}

func (p *Parser) parseEnum() Decl {
	enum := &Enum{}
	p.adv()

	if !p.expectid() {
		return nil
	} else {
		enum.name = p.curr.Literal
		p.adv()
	}

	if !p.expect(lex.LBRACE) {
		return nil
	}
	p.adv()

	for {
		if p.expectid() {
			enum.enums = append(enum.enums, p.curr.Literal)
			p.adv()
		} else {
			return nil
		}

		if p.is(lex.COMMA) {
			p.adv()
		} else {
			break
		}
	}

	if !p.expect(lex.RBRACE) {
		return nil
	}
	p.adv()

	return enum
}
