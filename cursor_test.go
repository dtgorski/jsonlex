// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 10/2020

package jsonlex

import (
	"bytes"
	"io"
	"testing"
)

func TestCursor_1(t *testing.T) {
	s := `{ "foo": -1 }`
	r := bytes.NewReader([]byte(s))
	c := NewCursor(r, nil)

	if n := c.Curr(); !n.Is(TokenLCB) {
		t.Errorf("unexpected")
	}
	if n := c.Peek(); !n.Is(TokenSTR) {
		t.Errorf("unexpected")
	}
	if n := c.Next(); !n.Is(TokenSTR) {
		t.Errorf("unexpected")
	}
	if n := c.Next(); !n.Is(TokenCOL) {
		t.Errorf("unexpected")
	}
	if n := c.Next(); !n.Is(TokenNUM) {
		t.Errorf("unexpected")
	}
	if n := c.Next(); !n.Is(TokenRCB) {
		t.Errorf("unexpected")
	}
	if n := c.Last(); !n.Is(TokenNUM) {
		t.Errorf("unexpected")
	}
	if n := c.Next(); !n.Is(TokenEOF) {
		t.Errorf("unexpected")
	}
	if n := c.Next(); !n.Is(TokenEOF) {
		t.Errorf("unexpected")
	}
}

func TestCursor_2(t *testing.T) {
	s := `{ "foo": -1 }`
	r := bytes.NewReader([]byte(s))
	f := func(k TokenKind, l []byte) bool {
		return !k.Is(TokenLCB) && !k.Is(TokenRCB) && !k.Is(TokenCOL)
	}
	c := NewCursor(r, f)

	if n := c.Curr(); !n.Is(TokenSTR) {
		t.Errorf("unexpected")
	}
	if n := c.Next(); !n.Is(TokenNUM) {
		t.Errorf("unexpected")
	}
	if n := c.Next(); !n.Is(TokenEOF) {
		t.Errorf("unexpected")
	}
}

func TestCursor_3(t *testing.T) {
	r := &FaultyReader{}
	c := NewCursor(r, nil)

	if n := c.Next(); !n.Is(TokenERR) {
		t.Errorf("unexpected")
	}
	if c.Next().String() != io.ErrUnexpectedEOF.Error() {
		t.Errorf("unexpected")
	}
}
