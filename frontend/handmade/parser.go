package handmade

import (
	. "Grobbit/common"
	"Grobbit/frontend/stub"
	"fmt"
	"slices"
)

type Parser struct {
	lexer LexerInterface
	token Token
}

///////////////////////////////////// Helper methods ///////////////////////////////////////

func (parser *Parser) oneOf(ttcheck TokenType, tts ...TokenType) bool {
	return slices.Contains(tts, ttcheck)
}

func (parser *Parser) matchError(msg string) {
	fmt.Println("Parsing error:", msg)
	panic("Parser.match: invalid token")
}

func (parser *Parser) match(tt TokenType) {
	if parser.token.Type != tt {
		msg := fmt.Sprintf("found '%s' but expected '%s' at position %s.", parser.token.Type, tt, parser.token.Pos)
		parser.matchError(msg)
	}
	parser.token = parser.lexer.NextToken()
}

func (parser *Parser) matchIf(tt TokenType) bool {
	if parser.token.Type == tt {
		parser.match(tt)
		return true
	}
	return false
}

///////////////////////////////////// Exported methods ///////////////////////////////////////

func (parser *Parser) ParseSrc(src []byte, handler ErrorHandler) string {
	var lex stub.Lexer // NOTE: You can change to your own lexer, if you like.
	lex.Init(src, handler)
	parser.Init(&lex)
	astree := parser.Parse()
	var pv PrintVisitor
	//ast.Print(nil, astree)  // For a more detailed printing.
	return pv.Root(astree)
}

func (parser *Parser) Init(lex LexerInterface) {
	parser.lexer = lex
	parser.token = parser.lexer.NextToken()
}

func (parser *Parser) Parse() Node {
	var (
		imports []*ImportSpecNode
		decls   []DeclNode
	)
	token := parser.token
	parser.match(TtKwPackage)
	ident := parser.Identifier()
	parser.match(TtSemicolon)
	for parser.token.Type == TtKwImport {
		specs := parser.ImportDecl()
		imports = append(imports, specs...)
		parser.match(TtSemicolon)
	}
	for parser.oneOf(parser.token.Type, TtKwConst, TtKwVar, TtKwFunc) {
		var decl DeclNode
		switch parser.token.Type {
		case TtKwConst:
			decl = parser.ConstDecl()
		case TtKwVar:
			decl = parser.VarDecl()
		case TtKwFunc:
			decl = parser.FuncDecl()
		default:
			// do nothing (should not happen)
		}
		decls = append(decls, decl)
		parser.match(TtSemicolon)
	}
	parser.match(TtEOI)
	return &FileNode{Tok: token, Name: ident, Imports: imports, Decls: decls}
}

func (parser *Parser) ImportDecl() []*ImportSpecNode {
	var (
		spec  *ImportSpecNode
		specs []*ImportSpecNode
	)
	parser.match(TtKwImport)
	if parser.matchIf(TtLParen) {
		for parser.oneOf(parser.token.Type, TtIdentifier, TtString) {
			spec = parser.ImportSpec()
			specs = append(specs, spec)
			parser.match(TtSemicolon)
		}
		parser.match(TtRParen)
	} else {
		spec = parser.ImportSpec()
		specs = append(specs, spec)
	}
	return specs
}

func (parser *Parser) ImportSpec() *ImportSpecNode {
	var identifier *IdentifierNode = nil
	token := parser.token
	if parser.token.Type == TtIdentifier {
		identifier = parser.Identifier()
	}
	path := parser.token.Lexeme
	parser.match(TtString)
	// Token is TtString, or TtIdentifier depending on whether path is prefixed or not.
	return &ImportSpecNode{Tok: token, Name: identifier, Path: path}
}

func (parser *Parser) VarDecl() DeclNode {
	var (
		spec  SpecNode
		specs []SpecNode
	)
	token := parser.token
	parser.match(TtKwVar)
	if parser.matchIf(TtLParen) {
		for parser.token.Type == TtIdentifier {
			spec = parser.VarSpec()
			specs = append(specs, spec)
			parser.match(TtSemicolon)
		}
		parser.match(TtRParen)
	} else {
		spec = parser.VarSpec()
		specs = append(specs, spec)
	}
	return &GenDeclNode{Tok: token, Specs: specs}
}

func (parser *Parser) VarSpec() *ValueSpecNode {
	var (
		exprs []ExprNode
	)
	token := parser.token
	idents := parser.Identifiers()
	typeExpr := parser.TypeName()
	if parser.matchIf(TtOpAssign) {
		expr := parser.BasicLiteral()
		exprs = append(exprs, expr)
		for parser.matchIf(TtComma) {
			expr = parser.BasicLiteral()
			exprs = append(exprs, expr)
		}
	}
	return &ValueSpecNode{Token: token, Names: idents, Type: typeExpr, Values: exprs}
}

func (parser *Parser) FuncDecl() DeclNode {
	var (
		params []*Field
		result ExprNode
	)
	token := parser.token
	parser.match(TtKwFunc)
	ident := parser.Identifier()
	parser.match(TtLParen)
	if parser.token.Type != TtRParen {
		params = parser.ParamDecl()
	}
	parser.match(TtRParen)
	if parser.token.Type != TtLBrace {
		result = parser.TypeName()
	}
	body := parser.BlockStatement()
	return &FuncDeclNode{Tok: token, Name: ident, Type: &FuncType{Params: params, Type: result}, Body: body}
}

func (parser *Parser) ParamDecl() []*Field {
	var fields []*Field
	field := parser.ParamSpec()
	fields = append(fields, field)
	for parser.token.Type == TtComma {
		field = parser.ParamSpec()
		fields = append(fields, field)
	}
	return fields
}

func (parser *Parser) ParamSpec() *Field {
	idents := parser.Identifiers()
	expr := parser.TypeName()
	return &Field{Names: idents, Type: expr}
}

func (parser *Parser) TypeName() ExprNode {
	return parser.Identifier()
}

func (parser *Parser) Identifiers() []*IdentifierNode {
	var idents []*IdentifierNode
	ident := parser.Identifier()
	idents = append(idents, ident)
	for parser.matchIf(TtComma) {
		ident := parser.Identifier()
		idents = append(idents, ident)
	}
	return idents
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func (parser *Parser) Statement() StmtNode {
	var stmt StmtNode = nil
	switch parser.token.Type {
	case TtKwConst, TtKwVar:
		stmt = parser.DeclStatement()
	case TtKwReturn:
		stmt = parser.ReturnStatement()
	case TtKwBreak:
		stmt = parser.BreakStatement()
	case TtKwContinue:
		stmt = parser.ContinueStatement()
	case TtKwIf:
		stmt = parser.IfStatement()
	case TtKwFor:
		stmt = parser.ForStatement()
	case TtLBrace:
		stmt = parser.BlockStatement()
	default:
		stmt = parser.SimpleStatement()
	}
	return stmt
}

func (parser *Parser) DeclStatement() *DeclStmtNode {
	var decl DeclNode
	token := parser.token
	if parser.matchIf(TtKwConst) {
		decl = parser.ConstDecl()
	} else {
		decl = parser.VarDecl()
	}
	return &DeclStmtNode{Tok: token, Decl: decl}
}

func (parser *Parser) BlockStatement() *BlockStmtNode {
	var stmts []StmtNode
	token := parser.token
	parser.match(TtLBrace)
	for parser.token.Type != TtRBrace {
		stmt := parser.Statement()
		stmts = append(stmts, stmt)
		parser.matchIf(TtSemicolon)
	}
	parser.match(TtRBrace)
	return &BlockStmtNode{Tok: token, List: stmts}
}

func (parser *Parser) ReturnStatement() *ReturnStmtNode {
	var expr ExprNode
	token := parser.token
	parser.match(TtKwReturn)
	if !parser.oneOf(parser.token.Type, TtSemicolon, TtRBrace) {
		expr = parser.Expression()
	}
	return &ReturnStmtNode{Token: token, Result: expr}
}

func (parser *Parser) ContinueStatement() StmtNode {
	token := parser.token
	parser.match(TtKwContinue)
	return &BranchStmtNode{Tok: token}
}

func (parser *Parser) ForStatement() *ForStmtNode {
	var init, cond, post StmtNode = nil, nil, nil
	token := parser.token
	parser.match(TtKwFor)
	if parser.token.Type != TtLBrace {
		init = parser.SimpleStatement()
		if parser.matchIf(TtSemicolon) { // ... for loop with a for clause
			if parser.token.Type != TtSemicolon {
				cond = parser.SimpleStatement()
			}
			parser.match(TtSemicolon)
			if parser.token.Type != TtLBrace {
				post = parser.SimpleStatement()
			}
		} else { // ... for loop with only a condition
			cond = init
			init = nil
		}
	}
	var condExpr ExprNode = nil
	if cond != nil {
		if tmp, isExpr := cond.(*ExprStmtNode); isExpr { // Ensure, condition is an ExpressionStmt
			condExpr = tmp.Expr
		} else {
			msg := fmt.Sprintf("Condition needs to be an expression at line %d.", parser.token.Pos.Line)
			parser.matchError(msg)
			// terminates
		}
	}
	block := parser.BlockStatement()
	return &ForStmtNode{Tok: token, Init: init, Cond: condExpr, Post: post, Body: block}
}

func (parser *Parser) SimpleStatement() StmtNode {
	var stmt StmtNode = nil
	var exprsLhs, exprsRhs []ExprNode
	exprsLhs = parser.Expressions()
	token := parser.token
	if parser.oneOf(parser.token.Type, TtOpDefine, TtOpAssign) {
		parser.match(parser.token.Type)
		exprsRhs = parser.Expressions()
		// We leave to it to semantic analysis to check that only identifier expressions on the left-hand-side
		stmt = &AssignStmtNode{Tok: token, Lhs: exprsLhs, Rhs: exprsRhs}
	} else {
		if len(exprsLhs) != 1 {
			msg := fmt.Sprintf("Expected only a single expression at position %s.", parser.token.Pos)
			parser.matchError(msg)
			// terminates program
		} else { // Expression statement
			// We leave to it to semantic analysis to check the correct expression type (call expression)
			stmt = &ExprStmtNode{Tok: token, Expr: exprsLhs[0]}
		}
	}
	return stmt
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func (parser *Parser) Expressions() []ExprNode {
	var exprs []ExprNode
	expr := parser.Expression()
	exprs = append(exprs, expr)
	for parser.matchIf(TtComma) {
		expr = parser.Expression()
		exprs = append(exprs, expr)
	}
	return exprs
}

func (parser *Parser) Expression() ExprNode {
	var expr ExprNode = nil
	expr = parser.ExprAnd()
	token := parser.token
	for parser.matchIf(TtOpOr) { // TtOpOr has the lowest precedence (see Go specs)
		rhs := parser.ExprAnd()
		expr = &BinaryExprNode{Tok: token, Lhs: expr, Rhs: rhs}
		token = parser.token
	}
	return expr
}

func (parser *Parser) Identifier() *IdentifierNode {
	node := &IdentifierNode{Tok: parser.token}
	parser.match(TtIdentifier)
	return node
}

func (parser *Parser) BasicLiteral() ExprNode {
	if !parser.oneOf(parser.token.Type, TtInt, TtFloat, TtString) {
		msg := fmt.Sprintf("Expected a literal but found '%s' at position %s.", parser.token.Type, parser.token.Pos)
		parser.matchError(msg)
	}
	token := parser.token
	parser.match(parser.token.Type)
	return &LiteralNode{Tok: token}
}

////////////////// You implement the methods below (add methods as needed) ////////////////////

func (parser *Parser) ConstDecl() DeclNode {
	// TO DO ...
	return nil
}

func (parser *Parser) ConstSpec() *ValueSpecNode {
	// TO DO ...
	return nil
}

func (parser *Parser) BreakStatement() StmtNode {
	// TO DO ...
	return nil
}

func (parser *Parser) IfStatement() StmtNode {
	// TO DO ...
	return nil
}

func (parser *Parser) ExprAnd() ExprNode {
	// TO DO ...

	// This code is a place-holder, replace with your implementation.
	var expr ExprNode = nil
	if parser.token.Type == TtIdentifier {
		expr = parser.Identifier()
	} else {
		expr = parser.BasicLiteral()
	}
	return expr
}

// Add functions as needed to parse expressions with the precedence (and associativity) of Grobbit operators correct
// (same as in Go). Note that you need to rewrite the grammar for reflecting the correct operator precedence.
