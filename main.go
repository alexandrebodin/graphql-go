package main

import "github.com/alexandrebodin/graphql-parser/parser"
import "fmt"

func main() {
	toks := parser.Parse("#azd zad   \n{ hello(id: -45, tap: 12345E-12, thrid: \"My super string\") { id, firstname, lastname, ...fra } } fragment fra on User { id }")
	fmt.Println(toks)
}
