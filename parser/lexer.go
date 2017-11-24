package parser

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

// A Lexer handles tokenization of the source code
type Lexer struct {
	source  string
	currPos int
}

// A Token is a lexed token
type Token struct {
	kind  tokenType
	value string
	pos   int
}

func (t Token) String() string {
	return fmt.Sprintf("%v(\"%s\")", t.kind, t.value)
}

type tokenType int

const (
	SOF tokenType = iota
	EOF
	BANG
	DOLLAR
	PAREN_L
	PAREN_R
	SPREAD
	COLON
	EQUALS
	AT
	BRACKET_L
	BRACKET_R
	BRACE_L
	PIPE
	BRACE_R
	NAME
	INT
	FLOAT
	STRING
	COMMENT
	ERROR
)

func (t tokenType) String() string {
	switch t {
	case SOF:
		return "SOF"
	case EOF:
		return "EOF"
	case BANG:
		return "BANG"
	case DOLLAR:
		return "DOLLAR"
	case PAREN_L:
		return "PAREN_L"
	case PAREN_R:
		return "PAREN_R"
	case SPREAD:
		return "SPREAD"
	case COLON:
		return "COLON"
	case EQUALS:
		return "EQUALS"
	case AT:
		return "AT"
	case BRACKET_L:
		return "BRACKET_L"
	case BRACKET_R:
		return "BRACKET_R"
	case BRACE_L:
		return "BRACE_L"
	case PIPE:
		return "PIPE"
	case BRACE_R:
		return "BRACE_R"
	case NAME:
		return "NAME"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case COMMENT:
		return "COMMENT"
	case ERROR:
		return "ERROR"
	}
	return ""
}

type lexError string

func (err lexError) Error() string {
	return string(err)
}

// CreateLexer returns a new lexer instance
func CreateLexer(source string) *Lexer {
	return &Lexer{source: source, currPos: 0}
}

// Run will start the lexing process
func (l *Lexer) Run() []Token {
	var tokens []Token

	for {
		token, err := l.readToken()
		if err != nil {
			fmt.Println(err)
			break
		}

		tokens = append(tokens, *token)

		if token.kind == EOF {
			break
		}
	}

	return tokens
}

func (l *Lexer) readToken() (*Token, error) {

	positionAfterWhitespace(l)

	if l.currPos >= len(l.source) {
		return &Token{EOF, "", l.currPos}, nil
	}

	startPos := l.currPos
	firstChar, size := utf8.DecodeRuneInString(l.source[l.currPos:])
	l.currPos += size

	if firstChar < 0x0020 && firstChar != 0x0009 && firstChar != 0x000A && firstChar != 0x000D {
		log.Fatalf("Invalid character %U", firstChar)
	}

	switch {
	case firstChar == '#':
		for l.currPos < len(l.source) {
			ch, size := utf8.DecodeRuneInString(l.source[l.currPos:])
			if ch >= 0x0020 || ch == 0x0009 {
				l.currPos += size
			} else {
				break //end of comment
			}
		}
		return &Token{kind: COMMENT, value: l.source[startPos:l.currPos], pos: startPos}, nil
	case firstChar == '{':
		return &Token{BRACE_L, "{", startPos}, nil
	case firstChar == '}':
		return &Token{BRACE_R, "}", startPos}, nil
	case firstChar == '!':
		return &Token{BANG, "!", startPos}, nil
	case firstChar == '$':
		return &Token{DOLLAR, "$", startPos}, nil
	case firstChar == '(':
		return &Token{PAREN_L, "(", startPos}, nil
	case firstChar == ')':
		return &Token{PAREN_R, ")", startPos}, nil
	case strings.HasPrefix(l.source[startPos:], "..."):
		l.currPos = startPos + 3
		return &Token{SPREAD, "...", startPos}, nil
	case firstChar == ':':
		return &Token{COLON, ":", startPos}, nil
	case firstChar == '=':
		return &Token{EQUALS, "=", startPos}, nil
	case firstChar == '@':
		return &Token{AT, "@", startPos}, nil
	case firstChar == '[':
		return &Token{BRACKET_L, "[", startPos}, nil
	case firstChar == ']':
		return &Token{BRACKET_R, "]", startPos}, nil
	case firstChar == '|':
		return &Token{PIPE, "|", startPos}, nil
	case ('_' == firstChar) || (firstChar >= 'A' && firstChar <= 'Z') || (firstChar >= 'a' && firstChar <= 'z'):
		for l.currPos < len(l.source) {
			ch, size := utf8.DecodeRuneInString(l.source[l.currPos:])
			if ch == '_' || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
				l.currPos += size
			} else {
				break //end of comment
			}
		}
		return &Token{kind: NAME, value: l.source[startPos:l.currPos], pos: startPos}, nil
	case (firstChar >= '0' && firstChar <= '9') || firstChar == '-':
		isFloat := false
		var code = firstChar
		var size int

		if code == '-' {
			code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
			l.currPos += size
		}

		if code == '0' {
			code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
			l.currPos += size
			if code >= '0' && code <= '9' {
				return &Token{}, lexError(fmt.Sprintf("Invalid number, unexpected digit after 0: %c", code))
			}
		} else {
			for l.currPos < len(l.source) {
				code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
				if code >= '0' && code <= '9' {
					l.currPos += size
				} else {
					break
				}
			}
			code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
		}

		if code == '.' {
			l.currPos += size
			isFloat = true

			for l.currPos < len(l.source) {
				code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
				if code >= '0' && code <= '9' {
					l.currPos += size
				} else {
					break
				}
			}
		}

		if code == 'E' || code == 'e' {
			l.currPos += size
			isFloat = true
			code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
			l.currPos += size
			if code == '-' || code == '+' {
				for l.currPos < len(l.source) {
					code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
					if code >= '0' && code <= '9' {
						l.currPos += size
					} else {
						break
					}
				}
			}
		}

		// read number
		if isFloat {
			return &Token{kind: FLOAT, value: l.source[startPos:l.currPos], pos: startPos}, nil
		}
		return &Token{kind: INT, value: l.source[startPos:l.currPos], pos: startPos}, nil
	case firstChar == '"':
		var code rune
		var size int
		val := ""

		for l.currPos < len(l.source) {
			code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
			if ((code >= 0x0020 && code <= 0xFFFF) || code == 0x0009) && code != '"' && code != '\\' && code != 0x000A && code != 0x000D {
				l.currPos += size
				val += string(code)
				continue
			}

			if code == '\\' {
				l.currPos += size
				code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
				switch code {
				case '"':
					val += "\""
					l.currPos += size
				case '/':
					val += "/"
					l.currPos += size
				case '\\':
					val += "\\"
					l.currPos += size
				case 'b':
					val += "\b"
					l.currPos += size
				case 'f':
					val += "\f"
					l.currPos += size
				case 'n':
					val += "\n"
					l.currPos += size
				case 'r':
					val += "\r"
					l.currPos += size
				case 't':
					val += "\t"
					l.currPos += size
				case 'u': // escape
					v, err := unicodeToChar(l.source[l.currPos+1:])

					if err != nil {
						return nil, lexError("Invalid escape sequence")
					}

					val += string(v)
					l.currPos += 5
				}
				continue
			}

			break
		}

		if code == '"' {
			l.currPos += size
		} else {
			return nil, lexError(fmt.Sprintf("Invalid end of string '%c'", code))
		}

		return &Token{kind: STRING, value: val, pos: startPos}, nil
	}

	return nil, lexError(fmt.Sprintf("Invalid charactère code '%c'", firstChar))
}

func positionAfterWhitespace(l *Lexer) {
	for l.currPos < len(l.source) {
		ch, size := utf8.DecodeRuneInString(l.source[l.currPos:])

		// ignore horizontal tab || space || new line || carriage return || comma | BOM
		if ch == 0x0009 || ch == 0x0020 || ch == 0x000A || ch == 0x000D || ch == ',' || ch == 0xFEFF {
			l.currPos += size
		} else {
			break
		}
	}
}

func unhex(b byte) (v rune, ok bool) {
	c := rune(b)
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return
}

func unicodeToChar(s string) (string, error) {
	var v rune
	for j := 0; j < 4; j++ {
		x, ok := unhex(s[j])
		if !ok {
			return "", errors.New("Invalid syntax")
		}
		v = v<<4 | x
	}

	return string(v), nil
}
