package handmade

import (
	. "Grobbit/common"
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
	r rune
	t TokenType
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
					return
				}
			}
		}
		lexer.token = Token{Type: Tt, Lexeme: string(startCh), Pos: startPos, PosEnd: lexer.pos}
	}
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

	// TODO: Modify the code in here.

	for lexer.processWhiteSpaces() {
		// Nothing ...
	}

	if lexer.eoi {
		lexer.setToken(TtEOI)
		return lexer.token
	}
	switch lexer.ch {
	case '(':
		lexer.setToken(TtLParen)
	case '*':
		lexer.setToken(TtOpMul, pair{'=', TtOpMulAssign})
	default:
		lexer.setToken(TtUnknown)
	}
	return lexer.token
}
