package json

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

var (
	Space     byte = ' '
	Tab       byte = '\t'
	NewLine   byte = '\n'
	Quotation byte = '"'
)

type TokenType string

func (t TokenType) IsEnding() bool {
	return t == ClosingBraceToken || t == ClosingCurlyToken ||
		t == ClosingBracketToken
}

const (
	UnknownToken        TokenType = "UNKNOWN"
	IntegerToken        TokenType = "INTEGER"
	FloatToken          TokenType = "FLOAT"
	StringToken         TokenType = "STRING"
	KeyToken            TokenType = "KEY_TOKEN"
	ColonToken          TokenType = ":"
	CommaToken          TokenType = ","
	WhiteSpaceToken     TokenType = "WHITESPACE"
	OpenBraceToken      TokenType = "["
	ClosingBraceToken   TokenType = "]"
	OpenBracketToken    TokenType = "("
	ClosingBracketToken TokenType = ")"
	OpenCurlyToken      TokenType = "{"
	ClosingCurlyToken   TokenType = "}"
	FullStopToken       TokenType = "."
)

var TokenTypes = map[string]TokenType{
	"UNKNOWN":    UnknownToken,
	"INTEGER":    IntegerToken,
	"FLOAT":      FloatToken,
	"STRING":     StringToken,
	"KEY_TOKEN":  KeyToken,
	":":          ColonToken,
	",":          CommaToken,
	"WHITESPACE": WhiteSpaceToken,
	"[":          OpenBraceToken,
	"]":          ClosingBraceToken,
	"(":          OpenBracketToken,
	")":          ClosingBracketToken,
	"{":          OpenCurlyToken,
	"}":          ClosingCurlyToken,
	".":          FullStopToken,
}

type Token struct {
	Value string
	Type  TokenType
}

func NewToken(b byte) Token {
	t, ok := TokenTypes[string(b)]
	if !ok {
		return Token{
			Value: string(b),
			Type:  UnknownToken,
		}
	}
	return Token{
		Value: string(b),
		Type:  t,
	}
}

func (token Token) ToValue() (reflect.Value, error) {
	switch token.Type {
	case StringToken:
		val := reflect.New(reflectTypeString).Elem()
		val.SetString(token.Value)
		return val, nil
	case IntegerToken:
		val := reflect.New(reflectTypeInteger).Elem()
		n, err := strconv.ParseInt(token.Value, 10, 64)
		if err != nil {
			return val, err
		}
		val.SetInt(n)
		return val, err
	case FloatToken:
		val := reflect.New(reflectTypeFloat).Elem()
		f, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return val, err
		}
		val.SetFloat(f)
		return val, err
	default:
		return reflect.New(nil), fmt.Errorf("cannot convert token to value: %v", token)
	}
}

func (token Token) IsEnding() bool {
	return token.Type.IsEnding()
}

func (token Token) String() string {
	return token.Value + ": " + string(token.Type)
}

type Tokens []Token

func (tokens Tokens) Join(sep string) string {
	var buf bytes.Buffer
	for i, token := range tokens {
		buf.WriteString(token.Value)
		if i < len(tokens)-1 {
			buf.WriteString(sep)
		}
	}
	return buf.String()
}
