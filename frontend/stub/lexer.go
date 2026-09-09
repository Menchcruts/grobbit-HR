package stub

import (
	. "Grobbit/common"
	"go/scanner"
	"go/token"
)

var tokenTypes = [...]TokenType{
	token.EOF:            TtEOI,
	token.ILLEGAL:        TtUnknown,
	token.IDENT:          TtIdentifier,
	token.INT:            TtInt,
	token.FLOAT:          TtFloat,
	token.CHAR:           TtChar,
	token.STRING:         TtString,
	token.ADD:            TtOpAdd,
	token.SUB:            TtOpSub,
	token.MUL:            TtOpMul,
	token.QUO:            TtOpDiv,
	token.REM:            TtOpMod,
	token.AND:            TtOpBitAnd,
	token.OR:             TtOpBitOr,
	token.XOR:            TtOpBitXor,
	token.SHL:            TtOpBitShl,
	token.SHR:            TtOpBitShr,
	token.AND_NOT:        TtOpBitAndNot,
	token.ADD_ASSIGN:     TtOpAddAssign,
	token.SUB_ASSIGN:     TtOpSubAssign,
	token.MUL_ASSIGN:     TtOpMulAssign,
	token.QUO_ASSIGN:     TtOpDivAssign,
	token.REM_ASSIGN:     TtOpModAssign,
	token.AND_ASSIGN:     TtOpBitAndAssign,
	token.OR_ASSIGN:      TtOpBitOrAssign,
	token.XOR_ASSIGN:     TtOpBitXorAssign,
	token.SHL_ASSIGN:     TtOpBitShlAssign,
	token.SHR_ASSIGN:     TtOpBitShrAssign,
	token.AND_NOT_ASSIGN: TtOpBitAndNotAssign,
	token.LAND:           TtOpAnd,
	token.LOR:            TtOpOr,
	token.ARROW:          TtOpArrow,
	token.INC:            TtOpInc,
	token.DEC:            TtOpDec,
	token.EQL:            TtOpEq,
	token.LSS:            TtOpLt,
	token.GTR:            TtOpGt,
	token.ASSIGN:         TtOpAssign,
	token.NOT:            TtOpNot,
	token.NEQ:            TtOpNe,
	token.LEQ:            TtOpLe,
	token.GEQ:            TtOpGe,
	token.DEFINE:         TtOpDefine,
	token.ELLIPSIS:       TtOpEllipsis,
	token.LPAREN:         TtLParen,
	token.LBRACK:         TtLBracket,
	token.LBRACE:         TtLBrace,
	token.COMMA:          TtComma,
	token.PERIOD:         TtPeriod,
	token.RPAREN:         TtRParen,
	token.RBRACK:         TtRBracket,
	token.RBRACE:         TtRBrace,
	token.SEMICOLON:      TtSemicolon,
	token.COLON:          TtColon,
	token.BREAK:          TtKwBreak,
	token.CASE:           TtKwCase,
	token.CHAN:           TtKwChan,
	token.CONST:          TtKwConst,
	token.CONTINUE:       TtKwContinue,
	token.DEFAULT:        TtKwDefault,
	token.DEFER:          TtKwDefer,
	token.ELSE:           TtKwElse,
	token.FALLTHROUGH:    TtKwFallthrough,
	token.FOR:            TtKwFor,
	token.FUNC:           TtKwFunc,
	token.GO:             TtKwGo,
	token.GOTO:           TtKwGo,
	token.IF:             TtKwIf,
	token.IMPORT:         TtKwImport,
	token.INTERFACE:      TtKwInterface,
	token.MAP:            TtKwMap,
	token.PACKAGE:        TtKwPackage,
	token.RANGE:          TtKwRange,
	token.RETURN:         TtKwReturn,
	token.SELECT:         TtKwSelect,
	token.STRUCT:         TtKwStruct,
	token.SWITCH:         TtKwSwitch,
	token.TYPE:           TtKwType,
	token.VAR:            TtKwVar,
	token.TILDE:          TtTilde,
}

type Lexer struct {
	fset    *token.FileSet
	scanner scanner.Scanner
}

// Init function initialises the lexical analysis.
func (lexer *Lexer) Init(src []byte, handler ErrorHandler) {
	lexer.fset = token.NewFileSet()
	file := lexer.fset.AddFile("", lexer.fset.Base(), len(src))
	errorHandler := func(pos token.Position, msg string) {
		cpos := Position{Line: pos.Line, Col: pos.Column}
		handler(cpos, msg)
	}
	lexer.scanner.Init(file, src, errorHandler, 2)
}

func (lexer *Lexer) NextToken() Token {
	pos, tok, lit := lexer.scanner.Scan()
	position := lexer.fset.Position(pos)
	endPosition := lexer.fset.Position(lexer.scanner.End())
	cpos := Position{Line: position.Line, Col: position.Column}
	clit := lit
	if len(lit) == 0 {
		clit = tok.String()
	}
	cendPos := Position{Line: endPosition.Line, Col: endPosition.Column}
	ctype := tokenTypes[tok]
	return Token{Type: ctype, Lexeme: clit, Pos: cpos, PosEnd: cendPos}
}
