package lexer

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
	Token   *Token
}

// A Token is a lexed token
type Token struct {
	Kind  TokenType
	Value string
	Pos   int
}

func newToken(typ TokenType, pos int, val ...string) *Token {
	if len(val) > 0 {
		return &Token{typ, val[0], pos}
	}
	return &Token{Kind: typ, Pos: pos}
}

func (t Token) String() string {
	if len(t.Value) > 0 {
		return fmt.Sprintf("%v(\"%s\")", t.Kind, t.Value)
	}

	return fmt.Sprintf("%v", t.Kind)
}

// TokenType represent the different token types
type TokenType int

const (
	SOF TokenType = iota
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
	LEX_ERROR
)

func (t TokenType) String() string {
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
	case LEX_ERROR:
		return "LEX_ERROR"
	}
	return ""
}

// New returns a new lexer instance
func New(source string) *Lexer {
	return &Lexer{source: source, currPos: 0, Token: &Token{Kind: SOF, Pos: 0}}
}

// Next returns the next valid token ignoring comments
func (l *Lexer) Next() *Token {
	var token *Token
	if l.Token.Kind != EOF {
		for {
			token = l.readToken()
			l.Token = token
			if token.Kind != COMMENT {
				break
			}
		}
	}

	return token
}

// ReadToken returns the next token
func (l *Lexer) readToken() *Token {
	positionAfterWhitespace(l)

	if l.currPos >= len(l.source) {
		return newToken(EOF, l.currPos)
	}

	start := l.currPos
	code, size := utf8.DecodeRuneInString(l.source[l.currPos:])
	l.currPos += size

	if code < 0x0020 && code != 0x0009 && code != 0x000A && code != 0x000D {
		log.Fatalf("Invalid character %U", code)
	}

	switch {
	case code == '#':
		return readComment(l, start)
	case code == '{':
		return newToken(BRACE_L, start)
	case code == '}':
		return newToken(BRACE_R, start)
	case code == '!':
		return newToken(BANG, start)
	case code == '$':
		return newToken(DOLLAR, start)
	case code == '(':
		return newToken(PAREN_L, start)
	case code == ')':
		return newToken(PAREN_R, start)
	case strings.HasPrefix(l.source[start:], "..."):
		l.currPos = start + 3
		return newToken(SPREAD, start)
	case code == ':':
		return newToken(COLON, start)
	case code == '=':
		return newToken(EQUALS, start)
	case code == '@':
		return newToken(AT, start)
	case code == '[':
		return newToken(BRACKET_L, start)
	case code == ']':
		return newToken(BRACKET_R, start)
	case code == '|':
		return newToken(PIPE, start)
	case ('_' == code) || (code >= 'A' && code <= 'Z') || (code >= 'a' && code <= 'z'):
		return readName(l, start)
	case (code >= '0' && code <= '9') || code == '-':
		return readNumber(l, start, code)
	case code == '"':
		return readString(l, start)
	}

	return newToken(LEX_ERROR, l.currPos, fmt.Sprintf("Invalid charactère code '%c'", code))
}

func positionAfterWhitespace(l *Lexer) {
	for {
		ch, size := utf8.DecodeRuneInString(l.source[l.currPos:])
		if ch == 0x0009 || ch == 0x0020 || ch == 0x000A || ch == 0x000D || ch == ',' || ch == 0xFEFF {
			l.currPos += size
			continue
		}
		break
	}
}

func readComment(l *Lexer, start int) *Token {
	for {
		ch, size := utf8.DecodeRuneInString(l.source[l.currPos:])
		if ch >= 0x0020 || ch == 0x0009 {
			l.currPos += size
			continue
		}
		break
	}
	return newToken(COMMENT, start, l.source[start:l.currPos])
}

func readName(l *Lexer, start int) *Token {
	for {
		ch, size := utf8.DecodeRuneInString(l.source[l.currPos:])
		if ch == '_' || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			l.currPos += size
			continue
		}
		break
	}
	return newToken(NAME, start, l.source[start:l.currPos])
}

func readNumber(l *Lexer, start int, code rune) *Token {
	isFloat := false
	var size int

	if code == '-' {
		code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
		l.currPos += size
	}

	if code == '0' {
		code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
		l.currPos += size
		if code >= '0' && code <= '9' {
			return newToken(LEX_ERROR, l.currPos, fmt.Sprintf("Invalid number, unexpected digit after 0: %c", code))
		}
	} else {
		for {
			code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
			if code >= '0' && code <= '9' {
				l.currPos += size
				continue
			}
			break
		}
		code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
	}

	if code == '.' {
		l.currPos += size
		isFloat = true

		for {
			code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
			if code >= '0' && code <= '9' {
				l.currPos += size
				continue
			}
			break
		}
	}

	if code == 'E' || code == 'e' {
		l.currPos += size
		isFloat = true
		code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
		l.currPos += size
		if code == '-' || code == '+' {
			for {
				code, size = utf8.DecodeRuneInString(l.source[l.currPos:])
				if code >= '0' && code <= '9' {
					l.currPos += size
					continue
				}
				break
			}
		}
	}

	if isFloat {
		return newToken(FLOAT, start, l.source[start:l.currPos])
	}
	return newToken(INT, start, l.source[start:l.currPos])
}

func readString(l *Lexer, start int) *Token {
	var code rune
	var size int
	val := ""

	for {
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
					return newToken(LEX_ERROR, l.currPos, "Invalid escape sequence")
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
		return newToken(LEX_ERROR, l.currPos, fmt.Sprintf("Invalid end of string %U", code))
	}

	return newToken(STRING, start, val)
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
