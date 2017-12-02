package parser

import (
	"fmt"

	"github.com/alexandrebodin/graphql-go/lexer"
)

// DefinitionType represents an AST Node Kind
type DefinitionType string

const (
	Query    DefinitionType = "Query"
	Mutation                = "Mutation"
	Fragment                = "Fragment"
)

// GraphQLDocument is the representation of a parsed graphql schema or query
type GraphQLDocument struct {
	Definitions GraphQLDefinitions
}

func (d GraphQLDocument) String() string {
	return fmt.Sprintf("GraphqlDocument{Defintions: %+v}", d.Definitions)
}

// GraphQLDefinition represent a graphql definition operation or fragment
type GraphQLDefinition struct {
	Kind         DefinitionType
	SelectionSet GraphqlSelectionSet
}

// GraphQLDefinitions a list of GraphQLDefinition
type GraphQLDefinitions []GraphQLDefinition

// GraphqlSelectionSet represent the information that an operation requests
type GraphqlSelectionSet struct {
}

// Parse parses a graphql schema or query
func Parse(source string) GraphQLDocument {
	lexer := lexer.New(source)
	return parseDocument(lexer)
}

func parseDocument(l *lexer.Lexer) GraphQLDocument {
	definitions := GraphQLDefinitions{}

	expect(l, lexer.SOF)
	for !skip(l, lexer.EOF) {
		definitions = append(definitions, parseDefinition(l))
	}

	return GraphQLDocument{definitions}
}

func parseDefinition(l *lexer.Lexer) GraphQLDefinition {
	if peek(l, lexer.BRACE_L) {
		return parserOperatonDefinition(l)
	}

	if peek(l, lexer.NAME) {
		switch l.Token.Value {
		case "query":
			fallthrough
		case "mutation":
		case "subscription":
			fmt.Printf("Parsing %v\n", l.Token.Value)
			return parserOperatonDefinition(l)
		case "fragment":
			fmt.Println("Parsing fragment")
			return GraphQLDefinition{Kind: Fragment}
		}
	}

	return GraphQLDefinition{}
}

func parserOperatonDefinition(l *lexer.Lexer) GraphQLDefinition {
	if peek(l, lexer.BRACE_L) {
		// shorthand query
		return GraphQLDefinition{
			Kind:         Query,
			SelectionSet: parseSelectionSet(l),
		}
	}
	return GraphQLDefinition{}
}

func parseSelectionSet(l *lexer.Lexer) GraphqlSelectionSet {
	for !skip(l, lexer.EOF) {
		l.Next()
	}
	return GraphqlSelectionSet{}
}

func expect(l *lexer.Lexer, kind lexer.TokenType) *lexer.Token {
	t := l.Token
	if t.Kind == kind {
		l.Next()
		return t
	}

	panic("Oh my GAAAAD the unexpected happened")
}

func skip(l *lexer.Lexer, kind lexer.TokenType) bool {
	match := l.Token.Kind == kind
	if match {
		l.Next()
	}

	return match
}

func peek(l *lexer.Lexer, kind lexer.TokenType) bool {
	return l.Token.Kind == kind
}
