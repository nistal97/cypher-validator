package parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type cypherListener struct {
	*BaseCypherListener
}

func (s *cypherListener) VisitErrorNode(node antlr.ErrorNode) {
	fmt.Printf("Error occured when Parsing %s..", node.GetText())
}

func ValidateCypher(exp string, t string) {
	// Setup the input
	is := antlr.NewInputStream(exp)

	// Create the Lexer
	lexer := NewCypherLexer(is)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	// Create the Parser
	p := NewCypherParser(stream)

	// Finally parse the expression (by walking the tree)
	var listener cypherListener
	switch t {
	case "create":
		antlr.ParseTreeWalkerDefault.Walk(&listener, p.OC_Create())
	case "merge":
		antlr.ParseTreeWalkerDefault.Walk(&listener, p.OC_Merge())
	}
}
