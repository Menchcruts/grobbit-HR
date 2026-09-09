package common

import (
	"fmt"
	"strings"
)

type PrintVisitor struct {
	VisitorWalker
	output strings.Builder
}

func (v *PrintVisitor) Root(node Node) string {
	v.output.Reset()
	v.visit(node)
	return v.output.String()
}

func (v *PrintVisitor) printf(format string, a ...any) {
	builder := strings.Builder{}
	for i := 0; i < v.Level; i++ {
		builder.WriteString(". ")
	}
	builder.WriteString(fmt.Sprintf(format, a...))
	v.output.WriteString(builder.String())
}

func (v *PrintVisitor) in(node Node) {
	v.printf("%T {\n", node)
}

func (v *PrintVisitor) out(node Node) {
	v.printf("}\n")
}

func (v *PrintVisitor) visit(node Node) {
	err := v.VisitDispatcher(v, node)
	if err != nil {
		v.printf("%s\n", err.Error())
	}
}

func (v *PrintVisitor) VisitIdentifierNode(node *IdentifierNode) {
	v.in(node)
	v.printf("  %s\n", node.Tok.Lexeme)
	v.out(node)
}

func (v *PrintVisitor) VisitLiteralNode(node *LiteralNode) {
	v.in(node)
	v.printf("  %s\n", node.Tok.Lexeme)
	v.out(node)
}

func (v *PrintVisitor) VisitCallExprNode(node *CallExprNode) {
	v.in(node)
	v.visit(node.Fun)
	for _, arg := range node.Args {
		v.visit(arg)
	}
	v.out(node)
}

func (v *PrintVisitor) VisitUnaryExprNode(node *UnaryExprNode) {
	v.in(node)
	v.printf("  %s\n", node.Tok.Lexeme)
	v.visit(node.Expr)
	v.out(node)
}

func (v *PrintVisitor) VisitBinaryExprNode(node *BinaryExprNode) {
	v.in(node)
	v.visit(node.Lhs)
	v.printf("  %s\n", node.Tok.Lexeme)
	v.visit(node.Rhs)
	v.out(node)
}

func (v *PrintVisitor) VisitDeclStmtNode(node *DeclStmtNode) {
	v.in(node)
	v.visit(node.Decl)
	v.out(node)
}

func (v *PrintVisitor) VisitExprStmtNode(node *ExprStmtNode) {
	v.in(node)
	v.visit(node.Expr)
	v.out(node)
}

func (v *PrintVisitor) VisitAssignStmtNode(node *AssignStmtNode) {
	v.in(node)
	for _, expr := range node.Lhs {
		v.visit(expr)
	}
	v.printf("  =\n")
	for _, expr := range node.Rhs {
		v.visit(expr)
	}
	v.out(node)
}

func (v *PrintVisitor) VisitReturnStmtNode(node *ReturnStmtNode) {
	v.in(node)
	v.printf("  %s\n", node.Token.Lexeme)
	if !IsNil(node.Result) {
		v.visit(node.Result)
	}
	v.out(node)
}

func (v *PrintVisitor) VisitBranchStmtNode(node *BranchStmtNode) {
	v.in(node)
	v.printf("  %s\n", node.Tok.Lexeme)
	v.out(node)
}

func (v *PrintVisitor) VisitBlockStmtNode(node *BlockStmtNode) {
	v.in(node)
	for _, stmt := range node.List {
		v.visit(stmt)
	}
	v.out(node)
}

func (v *PrintVisitor) VisitIfStmtNode(node *IfStmtNode) {
	v.in(node)
	if !IsNil(node.Init) {
		v.visit(node.Init)
	}
	v.visit(node.Cond)
	v.visit(node.Body)
	if !IsNil(node.Else) {
		v.visit(node.Else)
	}
	v.out(node)
}

func (v *PrintVisitor) VisitForStmtNode(node *ForStmtNode) {
	v.in(node)
	if !IsNil(node.Init) {
		v.visit(node.Init)
	}
	if !IsNil(node.Cond) {
		v.visit(node.Cond)
	}
	if !IsNil(node.Post) {
		v.visit(node.Post)
	}
	v.visit(node.Body)
	v.out(node)
}

func (v *PrintVisitor) VisitImportSpecNode(node *ImportSpecNode) {
	v.in(node)
	// v.visit(v, node.Name)  Always nil, Tok lexeme used instead to identify prefix.
	var prefix = ""
	if node.Tok.Type != TtString {
		prefix = node.Tok.Lexeme
	}
	v.printf("  %s %s\n", prefix, node.Path)
	v.out(node)
}

func (v *PrintVisitor) VisitValueSpecNode(node *ValueSpecNode) {
	v.in(node)
	for _, id := range node.Names {
		v.visit(id)
	}
	if !IsNil(node.Type) {
		v.visit(node.Type)
	}
	for _, expr := range node.Values {
		v.visit(expr)
	}
	v.out(node)
}

func (v *PrintVisitor) VisitGenDeclNode(node *GenDeclNode) {
	v.in(node)
	v.printf("  %s\n", node.Tok.Lexeme)
	for _, spec := range node.Specs {
		v.visit(spec)
	}
	v.out(node)
}

func (v *PrintVisitor) VisitFuncDeclNode(node *FuncDeclNode) {
	v.in(node)
	v.visit(node.Name)
	v.printf("  (\n")
	for _, field := range node.Type.Params {
		for _, id := range field.Names {
			v.visit(id)
		}
		v.visit(field.Type)
	}
	v.printf("  )\n")
	v.printf("  (\n")
	if !IsNil(node.Type.Type) {
		v.visit(node.Type.Type)
	}
	v.visit(node.Body)
	v.out(node)
}

func (v *PrintVisitor) VisitFileNode(node *FileNode) {
	v.in(node)
	v.visit(node.Name)
	for _, spec := range node.Imports {
		v.visit(spec)
	}
	for _, decl := range node.Decls {
		v.visit(decl)
	}
	v.out(node)
}
