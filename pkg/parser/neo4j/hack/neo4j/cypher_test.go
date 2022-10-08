package parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"testing"
)

func TestCypherParser(t *testing.T) {
	exp := "create (n:label1 {prop:1} )"

	// Setup the input
	is := antlr.NewInputStream(exp)

	// Create the Lexer
	lexer := NewCypherLexer(is)
	// Read all tokens
	for {
		t := lexer.NextToken()
		if t.GetTokenType() == antlr.TokenEOF {
			break
		}
		fmt.Printf("%s (%q)\n",
			lexer.SymbolicNames[t.GetTokenType()], t.GetText())
	}

	// ReCreate the Lexer
	is = antlr.NewInputStream(exp)
	lexer = NewCypherLexer(is)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	// Create the Parser
	p := NewCypherParser(stream)

	// Finally parse the expression (by walking the tree)
	var listener cypherListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, p.OC_Create())
}
