package handmade

import (
	. "Grobbit/common"
	"strings"
	"unicode"
)

type Lexer struct {
	src          []rune
	n            int
	i            int
	ch           rune
	eoi          bool
	token        Token
	pos          Position
	errorHandler ErrorHandler
}

var charEscapeSet = map[rune]bool{
	't':  true,
	'n':  true,
	'\\': true,
	'\'': true,
}

var stringEscapeSet = map[rune]bool{
	't':  true,
	'n':  true,
	'\\': true,
	'"':  true,
}

///////////////////// Helper functions ////////////////////////////

func (lexer *Lexer) peekRune() (rune, bool) {
	if lexer.eoi || lexer.i+1 >= lexer.n {
		return rune(0), true
	}
	return lexer.src[lexer.i+1], false
}

func (lexer *Lexer) peekRuneIs(check rune) bool {
	r, err := lexer.peekRune()
	return !err && r == check
}

func (lexer *Lexer) nextRune() {
	if lexer.eoi {
		return
	}
	lexer.i++
	if lexer.i >= lexer.n {
		lexer.ch, lexer.eoi = rune(0), true
		return
	}
	ch := lexer.src[lexer.i]
	if ch == '\n' { // Note: '\r'
		lexer.pos.Line += 1
		lexer.pos.Col = -1
	}
	lexer.pos.Col += 1
	lexer.ch, lexer.eoi = ch, false
}

type pair struct {
	r  rune
	t  TokenType
	sp []pair // Subpairs to check for
}

func (lexer *Lexer) setToken(Tt TokenType, pairs ...pair) {
	if len(pairs) == 0 {
		lexer.token = Token{Type: Tt, Lexeme: string(lexer.ch), Pos: lexer.pos}
		lexer.nextRune()
		lexer.token.PosEnd = lexer.pos
	} else {
		startCh, startPos := lexer.ch, lexer.pos
		lexer.nextRune()
		if !lexer.eoi {
			for _, p := range pairs {
				if lexer.ch == p.r {
					lexer.token = Token{Type: p.t, Lexeme: string(startCh) + string(lexer.ch), Pos: startPos}
					lexer.nextRune()
					lexer.token.PosEnd = lexer.pos
					lexer.token = lexer.mapToken(lexer.token, p.sp...)
					return
				}
			}
		}
		lexer.token = Token{Type: Tt, Lexeme: string(startCh), Pos: startPos, PosEnd: lexer.pos}
	}
}

func (lexer *Lexer) mapToken(token Token, subpairs ...pair) Token {
	if len(subpairs) != 0 {
		for _, p := range subpairs {
			if lexer.ch == p.r {
				token.Lexeme += string(lexer.ch)
				token.Type = p.t
				lexer.nextRune()
				token = lexer.mapToken(token, p.sp...)
			}
		}
	}

	return token
}

func (lexer *Lexer) resetToken(Tt TokenType) {
	lexer.token = Token{Type: Tt, Lexeme: lexer.token.Lexeme + string(lexer.ch), Pos: lexer.token.Pos}
	lexer.nextRune()
	lexer.token.PosEnd = lexer.pos
}

func (lexer *Lexer) processWhiteSpaces() bool {
	n := 0
	for !lexer.eoi && unicode.IsSpace(lexer.ch) {
		n += 1
		lexer.nextRune()
	}
	return n > 0
}

func (lexer *Lexer) processComments() bool {
	return lexer.processLineComment() || lexer.processBlockComment()
}

func (lexer *Lexer) processLineComment() bool {
	n := 0
	for !lexer.eoi && lexer.ch == '/' && lexer.peekRuneIs('/') {
		n += 1
		for !lexer.eoi {
			ch := lexer.ch
			lexer.nextRune()
			if ch == '\n' {
				break
			}
		}
	}
	return n > 0
}

func (lexer *Lexer) processBlockComment() bool {
	n := 0
	for !lexer.eoi && lexer.ch == '/' && lexer.peekRuneIs('*') {
		n += 1
		start := lexer.pos
		closed := false
		for !lexer.eoi {
			if lexer.ch == '*' && lexer.peekRuneIs('/') {
				lexer.nextRune() // *
				lexer.nextRune() // /
				closed = true
				break
			}
			lexer.nextRune()
		}
		if !closed {
			lexer.errorHandler(start, "comment not terminated")
		}
	}
	return n > 0
}

func (lexer *Lexer) matchOperator() bool {
	switch lexer.ch {
	case '+':
		lexer.setToken(
			TtOpAdd,
			pair{'=', TtOpAddAssign, nil},
			pair{'+', TtOpInc, nil},
		)
		return true
	case '-':
		lexer.setToken(
			TtOpSub,
			pair{'=', TtOpSubAssign, nil},
			pair{'-', TtOpDec, nil},
		)
		return true
	case '*':
		lexer.setToken(TtOpMul, pair{'=', TtOpMulAssign, nil})
		return true
	case '/':
		lexer.setToken(TtOpDiv, pair{'=', TtOpDivAssign, nil})
		return true
	case '%':
		lexer.setToken(TtOpMod, pair{'=', TtOpModAssign, nil})
		return true
	case '&':
		lexer.setToken(
			TtOpBitAnd,
			pair{'=', TtOpBitAndAssign, nil},
			pair{'&', TtOpAnd, nil},
			pair{'^', TtOpBitAndNot, []pair{
				{'=', TtOpBitAndNotAssign, nil},
			},
			})
		return true
	case '|':
		lexer.setToken(
			TtOpBitOr,
			pair{'|', TtOpOr, nil},
			pair{'=', TtOpBitOrAssign, nil},
		)
		return true
	case '<':
		lexer.setToken(
			TtOpLt,
			pair{'<', TtOpBitShl, []pair{
				{'=', TtOpBitShlAssign, nil},
			}},
			pair{'=', TtOpLe, nil},
			pair{'-', TtOpArrow, nil},
		)
		return true
	case '>':
		lexer.setToken(
			TtOpGt,
			pair{'>', TtOpBitShr, []pair{
				{'=', TtOpBitShrAssign, nil},
			}},
			pair{'=', TtOpGe, nil},
		)
		return true
	case '=':
		lexer.setToken(
			TtOpAssign,
			pair{'=', TtOpEq, nil},
		)
		return true
	case '!':
		lexer.setToken(
			TtOpNot,
			pair{'=', TtOpNe, nil},
		)
		return true
	case '^':
		lexer.setToken(
			TtOpBitXor,
			pair{'=', TtOpBitXorAssign, nil},
		)
		return true
	case '.':
		lexer.setToken(
			TtPeriod,
			pair{'.', TtUnknown, []pair{
				{'.', TtOpEllipsis, nil},
			}},
		)
		return true
	case '(':
		lexer.setToken(TtLParen)
		return true
	case ')':
		lexer.setToken(TtRParen)
		return true
	case '{':
		lexer.setToken(TtLBrace)
		return true
	case '}':
		lexer.setToken(TtRBrace)
		return true
	case '[':
		lexer.setToken(TtLBracket)
		return true
	case ']':
		lexer.setToken(TtRBracket)
		return true
	case ':':
		lexer.setToken(
			TtColon,
			pair{'=', TtOpDefine, nil},
		)
		return true
	case ';':
		lexer.setToken(TtSemicolon)
		return true
	case ',':
		lexer.setToken(TtComma)
		return true
	}
	return false
}

func (lexer *Lexer) buildLiteral() Token {
	if unicode.IsDigit(lexer.ch) {
		return lexer.buildNumericLiteral()
	} else if lexer.ch == '"' {
		return lexer.buildStringLiteral()
	} else {
		return lexer.buildWord()
	}
}

func (lexer *Lexer) buildWord() Token {
	start := lexer.pos
	sb := strings.Builder{}

	if unicode.IsLetter(lexer.ch) || lexer.ch == '_' {
		sb.WriteRune(lexer.ch)
		lexer.nextRune()
	}

	for unicode.IsLetter(lexer.ch) || lexer.ch == '_' || unicode.IsDigit(lexer.ch) {
		sb.WriteRune(lexer.ch)
		lexer.nextRune()
	}
	lexeme := sb.String()
	return Token{Type: TtIdentifier, Lexeme: lexeme, Pos: start, PosEnd: lexer.pos}
}

func (lexer *Lexer) buildStringLiteral() Token {
	start := lexer.pos
	sb := strings.Builder{}
	sb.WriteRune(lexer.ch)

	lexer.nextRune()
	for !lexer.eoi && lexer.ch != '"' && lexer.ch != '\n' {
		if lexer.ch == '\\' {
			sb.WriteRune(lexer.ch)
			lexer.nextRune()
			if lexer.eoi || lexer.ch == '\n' {
				break
			}
		}
		sb.WriteRune(lexer.ch)
		lexer.nextRune()
	}
	if lexer.ch == '"' {
		sb.WriteRune(lexer.ch)
		lexer.nextRune()
	} else {
		lexer.errorHandler(start, "string literal not terminated")
	}

	lexeme := sb.String()
	return Token{Type: TtString, Lexeme: lexeme, Pos: start, PosEnd: lexer.pos}
}

func (lexer *Lexer) buildNumericLiteral() Token {
	if lexer.ch != '0' {
		return lexer.buildDecimalLiteral()
	}

	// 0 or 0x/0b/0o prefix
	next, invalid := lexer.peekRune()
	if invalid {
		return lexer.buildDecimalLiteral()
	}

	switch next {
	case 'x', 'X':
		return lexer.buildHexLiteral()
	case 'b', 'B':
		return lexer.buildBinaryLiteral()
	case 'o', 'O':
		return lexer.buildOctalLiteral()
	default:
		return lexer.buildDecimalLiteral()
	}
}

func (lexer *Lexer) buildDecimalLiteral() Token {
	start := lexer.pos
	tt := TtInt

	sb := strings.Builder{}
	if lexer.ch == '0' && !lexer.peekRuneIs('.') {
		lexer.nextRune()
		return Token{Type: tt, Lexeme: "0", Pos: start, PosEnd: start}
	}

	sb.WriteRune(lexer.ch)
	lexer.nextRune()

	for unicode.IsDigit(lexer.ch) {
		sb.WriteRune(lexer.ch)
		lexer.nextRune()
	}

	if lexer.ch == '.' {
		sb.WriteRune(lexer.ch)
		lexer.nextRune()
		for unicode.IsDigit(lexer.ch) {
			sb.WriteRune(lexer.ch)
			lexer.nextRune()
		}
		tt = TtFloat
	}

	if lexer.ch == 'e' || lexer.ch == 'E' {
		sb.WriteRune(lexer.ch)
		lexer.nextRune()
		if lexer.ch == '+' || lexer.ch == '-' {
			sb.WriteRune(lexer.ch)
			lexer.nextRune()
		}
		for unicode.IsDigit(lexer.ch) {
			sb.WriteRune(lexer.ch)
			lexer.nextRune()
		}
		tt = TtFloat
	}

	lexeme := sb.String()
	return Token{Type: tt, Lexeme: lexeme, Pos: start, PosEnd: lexer.pos}
}

func (lexer *Lexer) buildHexLiteral() Token {
	start := lexer.pos
	tt := TtInt

	sb := strings.Builder{}
	sb.WriteRune(lexer.ch) // 0
	lexer.nextRune()

	sb.WriteRune(lexer.ch) // Should be 'x' or 'X'
	lexer.nextRune()

	for unicode.IsDigit(lexer.ch) || strings.ContainsRune("abcdefABCDEF", lexer.ch) || lexer.ch == '_' {
		sb.WriteRune(lexer.ch)
		lexer.nextRune()
	}

	if lexer.ch == '.' {
		sb.WriteRune(lexer.ch)
		lexer.nextRune()
		for unicode.IsDigit(lexer.ch) {
			sb.WriteRune(lexer.ch)
			lexer.nextRune()
		}
		tt = TtFloat
	}

	lexeme := sb.String()
	return Token{Type: tt, Lexeme: lexeme, Pos: start, PosEnd: lexer.pos}
}

func (lexer *Lexer) buildOctalLiteral() Token {
	start := lexer.pos

	sb := strings.Builder{}
	sb.WriteRune(lexer.ch) // 0
	lexer.nextRune()

	sb.WriteRune(lexer.ch) // Should be 'o' or 'O'
	lexer.nextRune()

	for unicode.IsDigit(lexer.ch) {
		if lexer.ch == '8' || lexer.ch == '9' {
			lexer.errorHandler(start, "invalid digit '"+string(lexer.ch)+"' in octal literal")
		}

		sb.WriteRune(lexer.ch)
		lexer.nextRune()
	}

	lexeme := sb.String()
	return Token{Type: TtInt, Lexeme: lexeme, Pos: start, PosEnd: lexer.pos}
}

func (lexer *Lexer) buildBinaryLiteral() Token {
	start := lexer.pos

	sb := strings.Builder{}
	sb.WriteRune(lexer.ch) // 0
	lexer.nextRune()

	sb.WriteRune(lexer.ch) // Should be 'b' or 'B'
	lexer.nextRune()

	tt := TtInt
	for unicode.IsDigit(lexer.ch) {
		if lexer.ch != '0' && lexer.ch != '1' {
			lexer.errorHandler(start, "invalid digit '"+string(lexer.ch)+"' in binary literal")
		}
		sb.WriteRune(lexer.ch)
		lexer.nextRune()
	}

	lexeme := sb.String()
	return Token{Type: tt, Lexeme: lexeme, Pos: start, PosEnd: lexer.pos}
}

///////////////////// Public functions ////////////////////////////

// Init function initialises the lexical analysis.
func (lexer *Lexer) Init(src []byte, handler ErrorHandler) {
	lexer.src, lexer.n, lexer.i = []rune(string(src)), len(src), -1
	lexer.eoi = false
	lexer.pos = Position{Line: 1, Col: 0}
	lexer.errorHandler = handler
	lexer.nextRune()
}

// NextToken reads and returns the next token.
func (lexer *Lexer) NextToken() Token {
	for lexer.processWhiteSpaces() || lexer.processComments() {
		// Nothing ...
	}

	if lexer.eoi {
		lexer.setToken(TtEOI)
		return lexer.token
	}

	if lexer.matchOperator() {
		return lexer.token
	}

	poss_literal := lexer.buildLiteral()

	keyword_tt := Lookup(poss_literal.Lexeme)
	if keyword_tt != TtIdentifier {
		poss_literal.Type = keyword_tt
	}
	lexer.token = poss_literal

	return lexer.token
}
