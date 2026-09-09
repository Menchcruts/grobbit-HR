package handmade

import (
	"Grobbit/common"
	"Grobbit/frontend/stub"
	"fmt"
	"os"
)

func outputString(file *os.File, line string) {
	fmt.Print(line)
	_, err := file.WriteString(line)
	if err != nil {
		panic(err)
	}
}

func runLexer(lex common.LexerInterface, txt string, filename string, outFilename string) {
	fmt.Printf("\nLEXER %s ==>\n", txt)
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
	lex.Init(src, errorHandler)
	for ctoken := lex.NextToken(); ctoken.Type != common.TtEOI; ctoken = lex.NextToken() {
		line := fmt.Sprintf("%-12s %-15s %q\n", ctoken.Pos, ctoken.Type, ctoken.Lexeme)
		outputString(ofile, line)
	}
}

func DoTestLexer() {
	var lex Lexer
	var lexBuiltin stub.Lexer
	filename := "tests/lexer/test_01.gr" //os.Args[1]
	runLexer(&lexBuiltin, "BUILTIN", filename, "out/output_b.txt")
	runLexer(&lex, "OWN", filename, "out/output_o.txt")
}
