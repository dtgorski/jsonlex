// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 10/2020

package jsonlex

import (
	"io"
)

type (
	// Cursor allows traversing the token stream.
	Cursor struct {
		reader  io.Reader
		filter  Filter
		lexer   *Lexer
		lastTok Token
		currTok Token
		nextTok Token
	}

	// Token is a container for token information.
	Token struct {
		Kind TokenKind
		Load []byte
		Pos  uint
	}

	// TokenKind denotes the type of token.
	TokenKind uint8

	// Filter is a callback function. It will be invoked when
	// the cursor is advanced. The callback must return whether
	// the token is accepted (true) or should dropped (false).
	// After a token is dropped, the scan for a next token continues.
	Filter func(kind TokenKind, load []byte) bool
)

// NewCursor creates and prepares a Cursor.
func NewCursor(r io.Reader, f Filter) *Cursor {
	c := &Cursor{
		reader: r,
		filter: f,
	}

	yield := func(kind TokenKind, load []byte, pos uint) bool {
		if c.currTok.Is(TokenERR) {
			return false
		}
		if c.filter != nil && !c.filter(kind, load) {
			return true
		}

		val := make([]byte, len(load))
		copy(val, load)

		c.lastTok = c.currTok
		c.currTok = c.nextTok
		c.nextTok = Token{kind, val, pos}

		return false
	}

	c.lexer = NewLexer(yield)
	c.Next()
	c.Next()

	return c
}

// Last returns the previous Token in stream.
// The underlying scanner position is not modified.
func (c *Cursor) Last() Token {
	return c.lastTok
}

// Curr function returns the current Token in stream.
// The underlying scanner position is not modified.
func (c *Cursor) Curr() Token {
	return c.currTok
}

// Peek returns the next Token in stream.
// The underlying scanner position is not modified.
func (c *Cursor) Peek() Token {
	return c.nextTok
}

// Next returns the next Token in stream. In contrast to
// the other methods, the underlying scanner position is
// modified.
func (c *Cursor) Next() Token {
	c.lexer.Scan(c.reader)
	return c.currTok
}

// Is is a convenience function.
func (t Token) Is(kind TokenKind) bool {
	return t.Kind == kind
}

func (t Token) String() string {
	return string(t.Load)
}

// Kinds of tokens emitted by the lexer.
const (
	TokenEOF TokenKind = iota // signals end of file/stream
	TokenERR                  // error string (other than EOF)
	TokenLIT                  // literal (true, false, null)
	TokenNUM                  // float number
	TokenSTR                  // "...\"..."
	TokenCOL                  // : colon
	TokenCOM                  // , comma
	TokenLSB                  // [ left square bracket
	TokenRSB                  // ] right square bracket
	TokenLCB                  // { left curly brace
	TokenRCB                  // } right curly brace

	scanning
)

// Is is a convenience function.
func (k TokenKind) Is(kind TokenKind) bool {
	return k == kind
}
