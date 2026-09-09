package common

//////////////////////////////////////////////////////////////////////////////////////

type Node any

// ExprNode implement the Expr interface.
type ExprNode interface {
	Node
	exprNode()
}

// StmtNode implement the Stmt interface.
type StmtNode interface {
	Node
	stmtNode()
}

// SpecNode implement the Spec interface.
type SpecNode interface {
	Node
	specNode()
}

// DeclNode implement the Decl interface.
type DeclNode interface {
	Node
	declNode()
}

//////////////////////////////////////////////////////////////////////////////////////

type (
	IdentifierNode struct {
		Tok Token
	}

	LiteralNode struct {
		Tok Token
	}

	CallExprNode struct {
		Tok  Token
		Fun  ExprNode   // function expression
		Args []ExprNode // function arguments; or nil
	}

	UnaryExprNode struct {
		Tok  Token
		Expr ExprNode // operand
	}

	BinaryExprNode struct {
		Tok Token
		Lhs ExprNode // left operand
		Rhs ExprNode // right operand
	}
)

// func (*BadExpr) exprNode()        {}
func (*IdentifierNode) exprNode() {}
func (*LiteralNode) exprNode()    {}
func (*CallExprNode) exprNode()   {}
func (*UnaryExprNode) exprNode()  {}
func (*BinaryExprNode) exprNode() {}

//////////////////////////////////////////////////////////////////////////////////////

type (

	// DeclStmtNode represents a declaration in a statement list.
	DeclStmtNode struct {
		Tok  Token
		Decl DeclNode // *GenDecl for 'const' and 'var'
	}

	ExprStmtNode struct {
		Tok  Token
		Expr ExprNode
	}

	AssignStmtNode struct {
		Tok Token
		Lhs []ExprNode
		Rhs []ExprNode
	}

	ReturnStmtNode struct {
		Token  Token
		Result ExprNode
	}

	BranchStmtNode struct {
		Tok Token // either 'break' or 'continue' statement
	}

	BlockStmtNode struct {
		Tok  Token
		List []StmtNode
	}

	IfStmtNode struct {
		Tok  Token
		Init StmtNode // initialization statement or nil
		Cond ExprNode // condition
		Body *BlockStmtNode
		Else StmtNode // else block, if-statement or nil
	}

	ForStmtNode struct {
		Tok  Token
		Init StmtNode // initialization statement; or nil
		Cond ExprNode // condition; or nil
		Post StmtNode // post iteration statement; or nil
		Body *BlockStmtNode
	}
)

// func (*BadStmt) stmtNode()        {}
func (*DeclStmtNode) stmtNode()   {}
func (*ExprStmtNode) stmtNode()   {}
func (*AssignStmtNode) stmtNode() {}
func (*ReturnStmtNode) stmtNode() {}
func (*BranchStmtNode) stmtNode() {}
func (*BlockStmtNode) stmtNode()  {}
func (*IfStmtNode) stmtNode()     {}
func (*ForStmtNode) stmtNode()    {}

//////////////////////////////////////////////////////////////////////////////////////

type (
	Field struct {
		Names []*IdentifierNode
		Type  ExprNode
	}

	FuncType struct {
		Params []*Field // (incoming) parameters; non-nil
		Type   ExprNode // (outgoing) result; or nil
	}

	ImportSpecNode struct { // Used for ImportSpec
		Tok  Token           // first token in spec: TtString or TtIdentifier
		Name *IdentifierNode // identifier or nil
		Path string          // import path
	}

	ValueSpecNode struct { // Used for ConstSpec and VarSpec
		Token  Token
		Names  []*IdentifierNode // value names (len(Names) > 0)
		Type   ExprNode          // value type; or nil
		Values []ExprNode        // initial values; or nil
	}
)

func (*ImportSpecNode) specNode() {}
func (*ValueSpecNode) specNode()  {}

type (
	GenDeclNode struct {
		Tok   Token // TtKwImport, TtKwConst, TtKwVar
		Specs []SpecNode
	}

	FuncDeclNode struct {
		Tok  Token           // TtKwFunc
		Name *IdentifierNode // function/method name
		Type *FuncType       // function signature: type and value parameters, result, and token of "func" keyword
		Body *BlockStmtNode  // function body; or nil for external (non-Go) function
	}
)

func (*GenDeclNode) declNode()  {}
func (*FuncDeclNode) declNode() {}

type FileNode struct {
	Tok     Token
	Name    *IdentifierNode   // package name
	Imports []*ImportSpecNode // imports in this file
	Decls   []DeclNode        // top-level declarations; or nil
}
