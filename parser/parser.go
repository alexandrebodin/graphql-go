package parser

// A GraphqlAST is a representation of a parsed graphql schema or query
type GraphqlAST struct {
	tokens []Token
}

// Parse parses a graphql schema or query
func Parse(source string) []Token {
	lexer := CreateLexer(source)
	return parseDocument(lexer)
}

func parseDocument(lexer *Lexer) []Token {
	return lexer.Run()
}
