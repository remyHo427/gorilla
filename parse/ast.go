package parse

import (
	"bytes"
	"gorilla/lex"
	"strconv"
	"strings"
)

type Node interface {
	String() string
}

type Expr interface {
	exprNode()
	Node
}
type Stmt interface {
	stmtNode()
	Node
}
type Decl interface {
	declNode()
	Node
}

// stmt
type ExprStmt struct {
	Expr Expr
}

func (s *ExprStmt) stmtNode() {}
func (s *ExprStmt) String() string {
	return s.Expr.String()
}

type IfStmt struct {
	If   Expr
	Then Stmt
	Else Stmt
}

func (s *IfStmt) stmtNode() {}
func (s *IfStmt) String() string {
	return join("if", s.If, s.Then, s.Else)
}

type BlockStmt struct {
	Stmts []Stmt
}

func (s *BlockStmt) stmtNode() {}
func (s *BlockStmt) String() string {
	return join("block", s.Stmts)
}

type WhileStmt struct {
	Cond Expr
	Loop Stmt
}

func (s *WhileStmt) stmtNode() {}
func (s *WhileStmt) String() string {
	return join("while", s.Cond, s.Loop)
}

type ReturnStmt struct {
	Return Expr
}

func (s *ReturnStmt) stmtNode() {}
func (s *ReturnStmt) String() string {
	return join("return", s.Return)
}

type BreakStmt struct{}

func (s *BreakStmt) stmtNode() {}
func (s *BreakStmt) String() string {
	return join("break")
}

type ContinueStmt struct{}

func (s *ContinueStmt) stmtNode() {}
func (s *ContinueStmt) String() string {
	return join("continue")
}

type NullStmt struct{}

func (s *NullStmt) stmtNode() {}
func (s *NullStmt) String() string {
	return join("null")
}

type DoStmt struct {
	Cond Expr
	Loop Stmt
}

func (s *DoStmt) stmtNode() {}
func (s *DoStmt) String() string {
	return join("do", s.Cond, s.Loop)
}

type ForStmt struct {
	Init Stmt
	Cond Stmt
	Post Expr
	Loop Stmt
}

func (s *ForStmt) stmtNode() {}
func (s *ForStmt) String() string {
	return join("for", s.Init, s.Cond, s.Post, s.Loop)
}

type SwitchStmt struct {
	Cond Expr
	Stmt Stmt
}

func (s *SwitchStmt) stmtNode() {}
func (s *SwitchStmt) String() string {
	return join("switch", s.Cond, s.Stmt)
}

type CaseStmt struct {
	Cond Expr
	Stmt Stmt
}

func (s *CaseStmt) stmtNode() {}
func (s *CaseStmt) String() string {
	return join("case", s.Cond, s.Stmt)
}

type DefaultStmt struct {
	Stmt Stmt
}

func (s *DefaultStmt) stmtNode() {}
func (s *DefaultStmt) String() string {
	return join("default", s.Stmt)
}

type DeclStmt struct {
	Decls []Decl
}

func (s *DeclStmt) stmtNode() {}
func (s *DeclStmt) String() string {
	return join("decl", s.Decls)
}

// decl
type StorageClass struct {
	Type uint
}

func (d *StorageClass) declNode() {}
func (d *StorageClass) String() string {
	return join("storage_class", d.Type)
}

type TypeQualifer struct {
	Type uint
}

func (d *TypeQualifer) declNode() {}
func (d *TypeQualifer) String() string {
	return join("type_qualifier", d.Type)
}

type DefaultTypeSpecifier struct {
	Type uint
}

func (d *DefaultTypeSpecifier) declNode() {}
func (d *DefaultTypeSpecifier) String() string {
	return join("default_type_specifier", d.Type)
}

type TypeSpecifier struct {
	Literal string
}

func (d *TypeSpecifier) declNode() {}
func (d *TypeSpecifier) String() string {
	return join("type_specifier", d.Literal)
}

type Enum struct {
	name  string
	enums []string
}

func (d *Enum) declNode() {}
func (d *Enum) String() string {
	return join("enum", d.name, d.enums)
}

// expr
type InfixExpr struct {
	Type  uint
	Left  Expr
	Right Expr
}

func (e *InfixExpr) exprNode() {}
func (e *InfixExpr) String() string {
	return join(e.Left, e.Type, e.Right)
}

type AssignExpr struct {
	Type  uint
	Expr  Expr
	Value Expr
}

func (e *AssignExpr) exprNode() {}
func (e *AssignExpr) String() string {
	return join(e.Expr, e.Type, e.Value)
}

type PrefixExpr struct {
	Type  uint
	Right Expr
}

func (e *PrefixExpr) exprNode() {}
func (e *PrefixExpr) String() string {
	return join(e.Type, e.Right)
}

type TernaryExpr struct {
	Cond Expr
	Then Expr
	Else Expr
}

func (e *TernaryExpr) exprNode() {}
func (e *TernaryExpr) String() string {
	return join(e.Cond, e.Then, e.Else)
}

type PostfixArithmeticExpr struct {
	Type uint
	Left Expr
}

func (e *PostfixArithmeticExpr) exprNode() {}
func (e *PostfixArithmeticExpr) String() string {
	return join(e.Left, e.Type)
}

type CallExpr struct {
	Callee Expr
	Args   []Expr
}

func (e *CallExpr) exprNode() {}
func (e *CallExpr) String() string {
	return join(e.Callee, e.Args)
}

type IndexExpr struct {
	Arr   Expr
	Index Expr
}

func (e *IndexExpr) exprNode() {}
func (e *IndexExpr) String() string {
	return join(e.Arr, e.Index)
}

type Int struct {
	Value int64
}

func (e *Int) exprNode() {}
func (e *Int) String() string {
	return strconv.Itoa(int(e.Value))
}

type Ident struct {
	Name string
}

func (e *Ident) exprNode() {}
func (e *Ident) String() string {
	return e.Name
}

// helper
func join(args ...any) string {
	var out bytes.Buffer

	out.WriteString("(")
	for i, a := range args {
		switch t := a.(type) {
		case Expr:
			out.WriteString(t.String())
		case Stmt:
			out.WriteString(t.String())
		case bytes.Buffer:
			out.WriteString(t.String())
		case string:
			out.WriteString(t)
		case uint:
			out.WriteString(lex.Tmap[t])
		case Decl:
			out.WriteString(t.String())
		case []Expr:
			for i, e := range t {
				out.WriteString(e.String())
				if i < len(t)-1 {
					out.WriteString(" ")
				}
			}
		case []Stmt:
			for i, s := range t {
				out.WriteString(s.String())
				if i < len(t)-1 {
					out.WriteString(" ")
				}
			}
		case []Decl:
			for i, d := range t {
				out.WriteString(d.String())
				if i < len(t)-1 {
					out.WriteString(" ")
				}
			}
		case []string:
			out.WriteString(strings.Join(t, " "))
		}

		if i < len(args)-1 {
			out.WriteString(" ")
		}
	}
	out.WriteString(")")

	return out.String()
}
