package handmade

import (
	"Grobbit/common"
	"Grobbit/frontend/stub"
	"fmt"
	"os"
)

func runParser(par common.ParserInterface, txt string, filename string, outFilename string) {
	fmt.Printf("\nPARSER %s ==>\n", txt)
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
	str := par.ParseSrc(src, errorHandler)
	outputString(ofile, str)
}

func DoTestParser() {
	filename := "tests/parser/test_01.gr" //os.Args[1]
	var parBuiltin stub.Parser
	var par Parser
	runParser(&parBuiltin, "BUILTIN", filename, "out/output_b.txt")
	runParser(&par, "OWN", filename, "out/output_o.txt")
}
