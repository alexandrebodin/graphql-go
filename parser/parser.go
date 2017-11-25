package parser

import (
	"fmt"

	"github.com/alexandrebodin/graphql-go/lexer"
)

// A GraphqlAST is a representation of a parsed graphql schema or query
type GraphqlAST struct {
	tokens []lexer.Token
}

// Parse parses a graphql schema or query
func Parse(source string) []lexer.Token {
	lexer := lexer.New(source)
	return parseDocument(lexer)
}

func parseDocument(l *lexer.Lexer) []lexer.Token {
	var tokens []lexer.Token

	for {
		token, err := l.ReadToken()
		if err != nil {
			fmt.Println(err)
			break
		}

		tokens = append(tokens, *token)

		if token.Kind == lexer.EOF {
			break
		}
	}

	return tokens
}
