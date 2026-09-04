package handmade

import (
	"Grobbit/common"
	"fmt"
	"go/scanner"
	"go/token"
	"os"
)

func outputString(file *os.File, line string) {
	fmt.Print(line)
	_, err := file.WriteString(line)
	if err != nil {
		panic(err)
	}
}

func runBuiltInLexer(filename, outFilename string) {

	fmt.Printf("\nLEXER BUILTIN ==>\n")

	src, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	ofile, err := os.Create(outFilename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer ofile.Close()

	fset := token.NewFileSet()
	var lex scanner.Scanner
	file := fset.AddFile(filename, fset.Base(), len(src))

	errorHandler := func(pos token.Position, msg string) {
		line := fmt.Sprintf("Error at (%d %d): %s\n", pos.Line, pos.Column, msg)
		outputString(ofile, line)
	}

	lex.Init(file, src, errorHandler, 2) // Mode 2: do not keep comments nor emit semicolons.

	for {
		pos, tok, lit := lex.Scan()
		if tok == token.EOF {
			break
		}
		position := fset.Position(pos)
		posstr := fmt.Sprintf("(%d %d)", position.Line, position.Column)
		if len(lit) == 0 {
			lit = tok.String()
		}
		line := fmt.Sprintf("%-12s %-15s %q\n", posstr, tok, lit)
		outputString(ofile, line)
	}
}

func runOwnLexer(filename, outFilename string) {

	fmt.Printf("\nLEXER OWN ==>\n")

	src, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	ofile, err := os.Create(outFilename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer ofile.Close()

	errorHandler := func(pos common.Position, msg string) {
		line := fmt.Sprintf("Error at (%d %d): %s\n", pos.Line, pos.Col, msg)
		outputString(ofile, line)
	}

	var lex Lexer
	lex.Init(src, errorHandler)

	for ctoken := lex.NextToken(); ctoken.Type != common.TtEOI; ctoken = lex.NextToken() {
		line := fmt.Sprintf("%-12s %-15s %q\n", ctoken.Pos, ctoken.Type, ctoken.Lexeme)
		outputString(ofile, line)
	}
}

func DoTestLexer() {
	filename := "tests/lexer/test_01.gr" //os.Args[1]
	runBuiltInLexer(filename, "out/output_b.txt")
	runOwnLexer(filename, "out/output_o.txt")
}
