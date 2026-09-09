package common

import (
	"errors"
	"fmt"
	"reflect"
)

type Visitor interface {
	VisitIdentifierNode(*IdentifierNode)
	VisitLiteralNode(*LiteralNode)
	VisitCallExprNode(*CallExprNode)
	VisitUnaryExprNode(*UnaryExprNode)
	VisitBinaryExprNode(*BinaryExprNode)
	VisitDeclStmtNode(*DeclStmtNode)
	VisitExprStmtNode(*ExprStmtNode)
	VisitAssignStmtNode(*AssignStmtNode)
	VisitReturnStmtNode(*ReturnStmtNode)
	VisitBranchStmtNode(*BranchStmtNode)
	VisitBlockStmtNode(*BlockStmtNode)
	VisitIfStmtNode(*IfStmtNode)
	VisitForStmtNode(*ForStmtNode)
	VisitImportSpecNode(*ImportSpecNode)
	VisitValueSpecNode(*ValueSpecNode)
	VisitGenDeclNode(*GenDeclNode)
	VisitFuncDeclNode(*FuncDeclNode)
	VisitFileNode(*FileNode)
}

// IsNil checks if an interface is nil
func IsNil(i any) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	return v.Kind() == reflect.Pointer && v.IsNil()
}

type VisitorWalker struct {
	Level int
}

func (w *VisitorWalker) VisitDispatcher(visitor Visitor, node Node) error {
	if IsNil(node) {
		return errors.New("dispatcher called with <nil> node, should not happen")
	}
	w.Level++
	switch n := node.(type) {
	case *IdentifierNode:
		visitor.VisitIdentifierNode(n)
	case *LiteralNode:
		visitor.VisitLiteralNode(n)
	case *CallExprNode:
		visitor.VisitCallExprNode(n)
	case *UnaryExprNode:
		visitor.VisitUnaryExprNode(n)
	case *BinaryExprNode:
		visitor.VisitBinaryExprNode(n)
	case *DeclStmtNode:
		visitor.VisitDeclStmtNode(n)
	case *ExprStmtNode:
		visitor.VisitExprStmtNode(n)
	case *AssignStmtNode:
		visitor.VisitAssignStmtNode(n)
	case *ReturnStmtNode:
		visitor.VisitReturnStmtNode(n)
	case *BranchStmtNode:
		visitor.VisitBranchStmtNode(n)
	case *BlockStmtNode:
		visitor.VisitBlockStmtNode(n)
	case *IfStmtNode:
		visitor.VisitIfStmtNode(n)
	case *ForStmtNode:
		visitor.VisitForStmtNode(n)
	case *ImportSpecNode:
		visitor.VisitImportSpecNode(n)
	case *ValueSpecNode:
		visitor.VisitValueSpecNode(n)
	case *GenDeclNode:
		visitor.VisitGenDeclNode(n)
	case *FuncDeclNode:
		visitor.VisitFuncDeclNode(n)
	case *FileNode:
		visitor.VisitFileNode(n)
	default:
		txt := fmt.Sprintf("unknown node type in dispatcher: %T, should not happen", node)
		return errors.New(txt)
	}
	w.Level--
	return nil
}
