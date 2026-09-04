package common

import (
	"fmt"
	"strconv"
)

/////////////////////////////////////////////////////////////////////////////////////////////

// TokenType tells the kind of the identified token.
type TokenType int

type Position struct {
	Line int
	Col  int
}

func (p Position) String() string {
	return fmt.Sprintf("(%d %d)", p.Line, p.Col)
}

type Token struct {
	Type   TokenType
	Lexeme string
	Pos    Position // location of start of token
	PosEnd Position // location immediately after end of token
}

type Error struct {
	Pos Position
	Msg string
}

func (e Error) String() string {
	return fmt.Sprintf("Error: '%s' '%s'.", e.Pos, e.Msg)
}

type ErrorHandler func(pos Position, msg string)

/////////////////////////////////////////////////////////////////////////////////////////////

const (
	TtEOI        TokenType = iota // End of input
	TtUnknown                     // Unknown token, returned if none other token matched
	TtIdentifier                  // Identifier

	beginLiteral // begin_ end_ constants do not serve as actual token types, used as markers (not exported)
	TtInt        // Integer literal
	TtFloat      // Float literal
	TtImag       // Imaginary literal
	TtChar       // Char/Rune literal
	TtString     // String literal
	endLiteral

	beginOperator
	TtOpAdd             // +
	TtOpSub             // -
	TtOpMul             // *
	TtOpDiv             // /
	TtOpMod             // %
	TtOpBitAnd          // &
	TtOpBitOr           //  |
	TtOpBitXor          // ^
	TtOpBitShl          // <<
	TtOpBitShr          // >>
	TtOpBitAndNot       // &^
	TtOpAssign          // =
	TtOpAddAssign       // +=
	TtOpSubAssign       // -=
	TtOpMulAssign       // *=
	TtOpDivAssign       // /=
	TtOpModAssign       // %=
	TtOpBitAndAssign    // &=
	TtOpBitOrAssign     //  |=
	TtOpBitXorAssign    // ^=
	TtOpBitShlAssign    // <<=
	TtOpBitShrAssign    // >>=
	TtOpBitAndNotAssign // &^=
	TtOpAnd             // &&
	TtOpOr              // ||
	TtOpNot             // !
	TtOpArrow           // <-
	TtOpInc             // ++
	TtOpDec             // --
	TtOpEq              // ==
	TtOpLt              // <
	TtOpGt              // >
	TtOpNe              // !=
	TtOpLe              // <=
	TtOpGe              // >=
	TtOpDefine          // :=
	TtOpEllipsis        // ...
	endOperator

	beginPunctuation
	TtLParen    // (
	TtLBracket  // [
	TtLBrace    // {
	TtRParen    // )
	TtRBracket  // ]
	TtRBrace    // }
	TtComma     // ,
	TtPeriod    // .
	TtSemicolon // ;
	TtColon     // :
	endPunctuation

	beginKeyword
	TtKwBreak
	TtKwCase
	TtKwChan
	TtKwConst
	TtKwContinue
	TtKwDefault
	TtKwDefer
	TtKwElse
	TtKwFallthrough
	TtKwFor
	TtKwFunc
	TtKwGo
	TtKwGoto
	TtKwIf
	TtKwImport
	TtKwInterface
	TtKwMap
	TtKwPackage
	TtKwRange
	TtKwReturn
	TtKwSelect
	TtKwStruct
	TtKwSwitch
	TtKwType
	TtKwVar
	endKeyword

	beginOther
	TtTilde // ~
	endOther
)

var tokenTypes = [...]string{
	TtEOI:               "EOI",
	TtUnknown:           "UNKNOWN",
	TtIdentifier:        "IDENT",
	TtInt:               "INT",
	TtFloat:             "FLOAT",
	TtChar:              "CHAR",
	TtString:            "STRING",
	TtOpAdd:             "+",
	TtOpSub:             "-",
	TtOpMul:             "*",
	TtOpDiv:             "/",
	TtOpMod:             "%",
	TtOpBitAnd:          "&",
	TtOpBitOr:           "|",
	TtOpBitXor:          "^",
	TtOpBitShl:          "<<",
	TtOpBitShr:          ">>",
	TtOpBitAndNot:       "&^",
	TtOpAddAssign:       "+=",
	TtOpSubAssign:       "-=",
	TtOpMulAssign:       "*=",
	TtOpDivAssign:       "/=",
	TtOpModAssign:       "%=",
	TtOpBitAndAssign:    "&=",
	TtOpBitOrAssign:     "|=",
	TtOpBitXorAssign:    "^=",
	TtOpBitShlAssign:    "<<=",
	TtOpBitShrAssign:    ">>=",
	TtOpBitAndNotAssign: "&^=",
	TtOpAnd:             "&&",
	TtOpOr:              "||",
	TtOpArrow:           "<-",
	TtOpInc:             "++",
	TtOpDec:             "--",
	TtOpEq:              "==",
	TtOpLt:              "<",
	TtOpGt:              ">",
	TtOpAssign:          "=",
	TtOpNot:             "!",
	TtOpNe:              "!=",
	TtOpLe:              "<=",
	TtOpGe:              ">=",
	TtOpDefine:          ":=",
	TtOpEllipsis:        "...",
	TtLParen:            "(",
	TtLBracket:          "[",
	TtLBrace:            "{",
	TtComma:             ",",
	TtPeriod:            ".",
	TtRParen:            ")",
	TtRBracket:          "]",
	TtRBrace:            "}",
	TtSemicolon:         ";",
	TtColon:             ":",
	TtKwBreak:           "break",
	TtKwCase:            "case",
	TtKwChan:            "chan",
	TtKwConst:           "const",
	TtKwContinue:        "continue",
	TtKwDefault:         "default",
	TtKwDefer:           "defer",
	TtKwElse:            "else",
	TtKwFallthrough:     "fallthrough",
	TtKwFor:             "for",
	TtKwFunc:            "func",
	TtKwGo:              "go",
	TtKwGoto:            "goto",
	TtKwIf:              "if",
	TtKwImport:          "import",
	TtKwInterface:       "interface",
	TtKwMap:             "map",
	TtKwPackage:         "package",
	TtKwRange:           "range",
	TtKwReturn:          "return",
	TtKwSelect:          "select",
	TtKwStruct:          "struct",
	TtKwSwitch:          "switch",
	TtKwType:            "type",
	TtKwVar:             "var",
	TtTilde:             "~",
}

var keywords map[string]TokenType

// String gives a string representation of the token type.
func (Tt TokenType) String() string {
	s := ""
	if 0 <= Tt && Tt < TokenType(len(tokenTypes)) {
		s = tokenTypes[Tt]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(Tt)) + ")"
	}
	return s
}

// Lookup tells whether a string is a keyword or not.
// If the string represents a keyword then its corresponding token type is returned, otherwise identifier token type.
func Lookup(ident string) TokenType {
	if tok, isKeyword := keywords[ident]; isKeyword {
		return tok
	}
	return TtIdentifier
}

/////////////////////////////////////////////////////////////////////////////////////////////

func init() {
	// Initializes (called automatically on startup)
	keywords = make(map[string]TokenType, endKeyword-(beginKeyword+1))
	for i := beginKeyword + 1; i < endKeyword; i++ {
		keywords[tokenTypes[i]] = i
	}
}
