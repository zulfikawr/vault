package rules

import (
	"fmt"
	"strconv"
)

// TokenType represents the type of a token
type TokenType int

const (
	TOKEN_EOF TokenType = iota
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER
	TOKEN_BOOL
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_AND
	TOKEN_OR
	TOKEN_EQ
	TOKEN_NEQ
	TOKEN_GT
	TOKEN_LT
	TOKEN_GTE
	TOKEN_LTE
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '(':
		tok = Token{Type: TOKEN_LPAREN, Value: string(l.ch)}
	case ')':
		tok = Token{Type: TOKEN_RPAREN, Value: string(l.ch)}
	case '=':
		tok = Token{Type: TOKEN_EQ, Value: "="}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_NEQ, Value: "!="}
		} else {
			tok = Token{Type: TOKEN_EOF, Value: "ILLEGAL"}
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_GTE, Value: ">="}
		} else {
			tok = Token{Type: TOKEN_GT, Value: ">"}
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_LTE, Value: "<="}
		} else {
			tok = Token{Type: TOKEN_LT, Value: "<"}
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = Token{Type: TOKEN_AND, Value: "&&"}
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = Token{Type: TOKEN_OR, Value: "||"}
		}
	case 0:
		tok = Token{Type: TOKEN_EOF, Value: ""}
	default:
		if isLetter(l.ch) || l.ch == '@' {
			tok.Value = l.readIdentifier()
			if tok.Value == "true" || tok.Value == "false" {
				tok.Type = TOKEN_BOOL
			} else {
				tok.Type = TOKEN_IDENTIFIER
			}
			return tok
		} else if isDigit(l.ch) {
			tok.Type = TOKEN_NUMBER
			tok.Value = l.readNumber()
			return tok
		} else if l.ch == '\'' || l.ch == '"' {
			tok.Type = TOKEN_STRING
			tok.Value = l.readString(l.ch)
			return tok
		} else {
			tok = Token{Type: TOKEN_EOF, Value: "ILLEGAL"}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '.' || l.ch == '_' || l.ch == '@' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString(quote byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == quote || l.ch == 0 {
			break
		}
	}
	str := l.input[position:l.position]
	l.readChar() // consume closing quote
	return str
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// AST Nodes

type Node interface {
	String() string
}

type InfixExpression struct {
	Left     Node
	Operator string
	Right    Node
}

func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) String() string {
	return fmt.Sprintf("'%s'", sl.Value)
}

type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", il.Value)
}

type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) String() string {
	return fmt.Sprintf("%t", bl.Value)
}

type Identifier struct {
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}

// Parser

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() (Node, error) {
	return p.parseExpression(LOWEST)
}

const (
	LOWEST      int = iota
	OR              // ||
	AND             // &&
	EQUALS          // ==, !=
	LESSGREATER     // >, <, >=, <=
)

var precedences = map[TokenType]int{
	TOKEN_EQ:  EQUALS,
	TOKEN_NEQ: EQUALS,
	TOKEN_LT:  LESSGREATER,
	TOKEN_GT:  LESSGREATER,
	TOKEN_LTE: LESSGREATER,
	TOKEN_GTE: LESSGREATER,
	TOKEN_AND: AND,
	TOKEN_OR:  OR,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseExpression(precedence int) (Node, error) {
	left, err := p.parsePrefix()
	if err != nil {
		return nil, err
	}

	for p.peekToken.Type != TOKEN_EOF && precedence < p.peekPrecedence() {
		p.nextToken()
		left, err = p.parseInfix(left)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

func (p *Parser) parsePrefix() (Node, error) {
	switch p.curToken.Type {
	case TOKEN_IDENTIFIER:
		return &Identifier{Value: p.curToken.Value}, nil
	case TOKEN_STRING:
		return &StringLiteral{Value: p.curToken.Value}, nil
	case TOKEN_NUMBER:
		val, err := strconv.ParseInt(p.curToken.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse %q as integer", p.curToken.Value)
		}
		return &IntegerLiteral{Value: val}, nil
	case TOKEN_BOOL:
		return &BooleanLiteral{Value: p.curToken.Value == "true"}, nil
	case TOKEN_LPAREN:
		p.nextToken()
		exp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		if p.peekToken.Type != TOKEN_RPAREN {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		p.nextToken()
		return exp, nil
	default:
		return nil, fmt.Errorf("unexpected token %v", p.curToken)
	}
}

func (p *Parser) parseInfix(left Node) (Node, error) {
	expression := &InfixExpression{
		Left:     left,
		Operator: p.curToken.Value,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}
	expression.Right = right
	return expression, nil
}
