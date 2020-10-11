package json

import "errors"

var (
	errInvalidJSONString = errors.New("invalid JSON string")
)

func Lex(input string) (Tokens, error) {
	l := &lexer{
		input: input,
		index: -1,
	}

	for l.next() {
		switch l.value() {
		case Space, Tab, NewLine:
			continue
		case Quotation:
			token, err := l.readString()
			if err != nil {
				return Tokens{}, errInvalidJSONString
			}
			l.output = append(l.output, token)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			token, err := l.readNumber()
			if err != nil {
				return Tokens{}, errInvalidJSONString
			}
			l.output = append(l.output, token)
			l.index--
		case 't':
			l.output = append(l.output, Token{"true", BooleanToken})
			l.index += len("rue")
		case 'f':
			l.output = append(l.output, Token{"false", BooleanToken})
			l.index += len("alse")
		case 'n':
			l.output = append(l.output, Token{"null", NullToken})
			l.index += len("ull")
		default:
			l.output = append(l.output, NewToken(l.value()))
		}
	}
	return l.output, nil
}

type lexer struct {
	index  int
	input  string
	output []Token
}

func (l *lexer) next() bool {
	l.index++
	return l.index < len(l.input)
}

func (l *lexer) value() byte {
	return l.input[l.index]
}

func (l *lexer) readNumber() (Token, error) {
	tokenType := IntegerToken
	buf := []byte{l.value()}
	for l.next() {
		switch l.value() {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			buf = append(buf, l.value())
		case '.':
			tokenType = FloatToken
			buf = append(buf, l.value())
		default:
			return Token{
				Value: string(buf),
				Type:  tokenType,
			}, nil
		}
	}
	return Token{
		Value: string(buf),
		Type:  tokenType,
	}, nil
}

func (l *lexer) readString() (Token, error) {
	var buf []byte
	for l.next() {
		if l.value() == Quotation {
			return Token{
				Value: string(buf),
				Type:  StringToken,
			}, nil
		}
		buf = append(buf, l.value())
	}
	return Token{}, errors.New("unexpected end of file, trying to read string")
}
