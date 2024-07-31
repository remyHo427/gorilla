package parse

import (
	"bytes"
	"strconv"
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
			out.WriteString(tok_strmap[t])
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
		}

		if i < len(args)-1 {
			out.WriteString(" ")
		}
	}
	out.WriteString(")")

	return out.String()
}
