package parser

import (
	"errors"
	"fmt"
	"log"

	"github.com/alexandrebodin/graphql-go/lexer"
)

// DefinitionType represents an AST Node Kind
type DefinitionType string

const (
	Query        DefinitionType = "Query"
	Mutation                    = "Mutation"
	Subscription                = "Subscription"
	Fragment                    = "Fragment"
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
	Kind                DefinitionType
	SelectionSet        GraphQLSelectionSet
	VariableDefinitions GraphQLVariableDefinitions
	Name                GraphQLName
}

type GraphQLVariableDefinitions []GraphQLVariableDefinition

type GraphQLVariableDefinition struct {
}

// GraphQLDefinitions a list of GraphQLDefinition
type GraphQLDefinitions []GraphQLDefinition

// GraphQLSelectionSet represent the information that an operation requests
type GraphQLSelectionSet struct {
}

// GraphQLName can be an operation name or sth else
type GraphQLName struct {
	value string
}

func (n GraphQLName) String() string {
	return fmt.Sprintf("GraphqlName{value: %s}", n.value)
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
		def, err := parseDefinition(l)
		if err != nil {
			log.Fatal(err)
		}
		definitions = append(definitions, def)
	}

	return GraphQLDocument{definitions}
}

func parseDefinition(l *lexer.Lexer) (GraphQLDefinition, error) {
	if peek(l, lexer.BRACE_L) {
		return parserOperatonDefinition(l), nil
	}

	if peek(l, lexer.NAME) {
		switch l.Token.Value {
		case "query":
			fallthrough
		case "mutation":
			fallthrough
		case "subscription":
			fmt.Printf("Parsing %v\n", l.Token.Value)
			return parserOperatonDefinition(l), nil
		case "fragment":
			fmt.Println("Parsing fragment")
			return GraphQLDefinition{Kind: Fragment}, nil
		}
	}

	return GraphQLDefinition{}, errors.New("Invalid operation")
}

func parserOperatonDefinition(l *lexer.Lexer) GraphQLDefinition {

	if peek(l, lexer.BRACE_L) {
		// shorthand query
		return GraphQLDefinition{
			Kind:         Query,
			SelectionSet: parseSelectionSet(l),
		}
	}

	t := expect(l, lexer.NAME)
	var opType DefinitionType
	switch t.Value {
	case "query":
		opType = Query
	case "mutation":
		opType = Mutation
	case "subscription":
		opType = Subscription
	}

	name := GraphQLName{}
	if peek(l, lexer.NAME) {
		nameToken := expect(l, lexer.NAME)
		name = GraphQLName{nameToken.Value}
	}

	return GraphQLDefinition{
		Kind:                opType,
		Name:                name,
		VariableDefinitions: parseVariableDefinitions(l),
		SelectionSet:        parseSelectionSet(l),
	}
}

func parseVariableDefinitions(l *lexer.Lexer) GraphQLVariableDefinitions {
	if peek(l, lexer.PAREN_L) {
		expect(l, lexer.PAREN_L)
		definitions := GraphQLVariableDefinitions{}

		for !skip(l, lexer.PAREN_R) {
			definitions = append(definitions, parseVariableDefinition(l))
		}

		return definitions
	}

	return GraphQLVariableDefinitions{}
}

func parseVariableDefinition(l *lexer.Lexer) GraphQLVariableDefinition {
	return GraphQLVariableDefinition{}
}

func parseSelectionSet(l *lexer.Lexer) GraphQLSelectionSet {
	for !skip(l, lexer.EOF) {
		l.Next()
	}
	return GraphQLSelectionSet{}
}

func expect(l *lexer.Lexer, kind lexer.TokenType) *lexer.Token {
	t := l.Token
	if t.Kind == kind {
		l.Next()
		return t
	}

	err := fmt.Errorf("Unexpected token %s, expected %s", t.Kind, kind)
	panic(err)
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
