package stub

import (
	"Grobbit/common"
	"go/ast"
	par "go/parser"
	"go/token"
	"strings"
)

type Parser struct {
}

func (parser *Parser) ParseSrc(src []byte, handler common.ErrorHandler) string {
	// Ignoring handler, needed for interface.
	fset := token.NewFileSet()
	node, err := par.ParseFile(fset, "", src, 2)
	if err != nil {
		panic(err)
	}
	var strBuilder strings.Builder
	err = ast.Fprint(&strBuilder, fset, node, nil)
	if err != nil {
		panic(err)
	}
	return strBuilder.String()
}
