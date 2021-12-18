// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 10/2020

package jsonlex

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// expect EOF
func TestLexer_Scan_1(t *testing.T) {
	s := ``
	i := 0
	y := func(kind TokenKind, load []byte, pos uint) bool {
		i++
		if !kind.Is(TokenEOF) {
			t.Errorf("unexpected %q", load)
		}
		return true
	}
	l := NewLexer(y)
	r := bytes.NewReader([]byte(s))
	l.Scan(r)

	if i != 1 {
		t.Error("unexpected")
	}
}

// expect error, unexpected input
func TestLexer_Scan_2(t *testing.T) {
	s := ` * `
	i := 0
	y := func(kind TokenKind, load []byte, pos uint) bool {
		i++
		if !kind.Is(TokenERR) {
			t.Errorf("unexpected %q", load)
		}
		return true
	}
	l := NewLexer(y)
	r := bytes.NewReader([]byte(s))
	l.Scan(r)

	if i != 1 {
		t.Error("unexpected")
	}
}

// expect standard functionality
func TestLexer_Scan_3(t *testing.T) {
	s := ` { "foo": "bar", "b\"az": [ null, true, false, -42, "false" ] } `

	e := []struct {
		kind TokenKind
		load []byte
	}{
		{kind: TokenLCB, load: []byte(`{`)},
		{kind: TokenSTR, load: []byte(`foo`)},
		{kind: TokenCOL, load: []byte(`:`)},
		{kind: TokenSTR, load: []byte(`bar`)},
		{kind: TokenCOM, load: []byte(`,`)},
		{kind: TokenSTR, load: []byte(`b\"az`)},
		{kind: TokenCOL, load: []byte(`:`)},
		{kind: TokenLSB, load: []byte(`[`)},
		{kind: TokenLIT, load: []byte(`null`)},
		{kind: TokenCOM, load: []byte(`,`)},
		{kind: TokenLIT, load: []byte(`true`)},
		{kind: TokenCOM, load: []byte(`,`)},
		{kind: TokenLIT, load: []byte(`false`)},
		{kind: TokenCOM, load: []byte(`,`)},
		{kind: TokenNUM, load: []byte(`-42`)},
		{kind: TokenCOM, load: []byte(`,`)},
		{kind: TokenSTR, load: []byte(`false`)},
		{kind: TokenRSB, load: []byte(`]`)},
		{kind: TokenRCB, load: []byte(`}`)},
		{kind: TokenEOF, load: nil},
	}

	i := 0
	y := func(kind TokenKind, load []byte, pos uint) bool {
		if !e[i].kind.Is(kind) {
			t.Errorf("unexpected %q", kind)
		}
		if !bytes.Equal(e[i].load, load) {
			t.Errorf("unexpected %q", load)
		}
		i++
		return true
	}
	l := NewLexer(y)
	r := bytes.NewReader([]byte(s))
	l.Scan(r)
}

// expect no errors while tokenizing floats and other valid literals
func TestLexer_Scan_4(t *testing.T) {
	s := []string{
		"-0",
		"-1",
		"0.1e-20",
		"1.e+5",
		"1.0",
		"1e+1",
		"-.0E+0",
		"1E-0",
		"1E-1",
		":",
		"true",
		"false",
		"null",
	}

	i := 0
	y := func(kind TokenKind, load []byte, pos uint) bool {
		i++
		if kind.Is(TokenERR) {
			t.Errorf("unexpected %q", load)
			return false
		}
		return true
	}
	for _, v := range s {
		i = 0
		l := NewLexer(y)
		r := bytes.NewReader([]byte(v))
		l.Scan(r)
		if i != 2 {
			t.Error("unexpected")
		}
	}
}

// expect errors while tokenizing broken floats
func TestLexer_Scan_5(t *testing.T) {
	s := []string{
		"-",
		"--",
		"+1",
		".",
		"-0.",
		"-E",
		"-e",
		".E",
		".e",
		"-.E",
		"-.e",
		"1e",
		"-.e0",
		".e0",
		"1E-+0",
		"1e.",
	}

	i := 0
	y := func(kind TokenKind, load []byte, pos uint) bool {
		i++
		if !kind.Is(TokenERR) {
			t.Errorf("unexpected %q %q", kind, load)
		}
		return true
	}
	for _, v := range s {
		i = 0
		l := NewLexer(y)
		r := bytes.NewReader([]byte(v))
		l.Scan(r)
		if i != 1 {
			t.Error("unexpected")
		}
	}
	for _, v := range s {
		i = 0
		l := NewLexer(y)
		v += " " // ws after token
		r := bytes.NewReader([]byte(v))
		l.Scan(r)
		if i != 1 {
			t.Error("unexpected")
		}
	}
}

// expect error when byte stream contains illegal values
func TestLexer_Scan_6(t *testing.T) {
	s := []byte{0x05, 0x7F, 0x80}

	i := 0
	y := func(kind TokenKind, load []byte, pos uint) bool {
		i++
		if !kind.Is(TokenERR) {
			t.Errorf("unexpected %q", load)
		}
		return true
	}
	for _, v := range s {
		i = 0
		l := NewLexer(y)
		r := bytes.NewReader([]byte{v})
		l.Scan(r)
		if i != 1 {
			t.Error("unexpected")
		}
	}
}

// expect error when malformed tokens found
func TestLexer_Scan_7(t *testing.T) {
	s := `frue nalse tull`

	i := 0
	y := func(kind TokenKind, load []byte, pos uint) bool {
		i++
		if !kind.Is(TokenERR) {
			t.Errorf("unexpected %d %q", kind, load)
		}
		return true
	}
	for _, v := range strings.Split(s, " ") {
		i = 0
		l := NewLexer(y)
		r := bytes.NewReader([]byte(v))
		l.Scan(r)
		if i != 1 {
			t.Error("unexpected")
		}
	}
	for _, v := range strings.Split(s, " ") {
		i = 0
		l := NewLexer(y)
		v += " " // ws after token
		r := bytes.NewReader([]byte(v))
		l.Scan(r)
		if i != 1 {
			t.Error("unexpected")
		}
	}
}

// re-entrance
func TestLexer_Scan_8(t *testing.T) {
	s := ` { } `
	i := 0
	y := func(kind TokenKind, load []byte, pos uint) bool {
		i++
		if i == 1 && !kind.Is(TokenLCB) {
			t.Errorf("unexpected %q", load)
		}
		if i == 2 && !kind.Is(TokenRCB) {
			t.Errorf("unexpected %q", load)
		}
		if i == 3 && !kind.Is(TokenEOF) {
			t.Errorf("unexpected %q", load)
		}
		if i == 4 && !kind.Is(TokenEOF) {
			t.Errorf("unexpected %q", load)
		}
		return false
	}
	l := NewLexer(y)
	r := bytes.NewReader([]byte(s))

	l.Scan(r)
	l.Scan(r)
	l.Scan(r)
	l.Scan(r)
}

// expect error when reader fails with io.ErrUnexpectedEOF
func TestLexer_Scan_9(t *testing.T) {
	y := func(kind TokenKind, load []byte, pos uint) bool {
		if !kind.Is(TokenERR) {
			t.Errorf("unexpected %q", load)
		}
		return true
	}
	l := NewLexer(y)
	r := &FaultyReader{}
	l.Scan(r)
}

type FaultyReader struct{}

func (*FaultyReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

// ensure LexerOptEnableUnreadBuffer is working with objects
func TestLexer_Scan_10(t *testing.T) {
	s := []byte(`{"hello":   "world", "a": 1, "b": false}`)
	r := bytes.NewBuffer(s)

	steps := []struct {
		kind TokenKind
		load string
		left string
	}{{
		kind: TokenLCB, //0
		load: `{`,
		left: `"hello":   "world", "a": 1, "b": false}`,
	}, {
		kind: TokenSTR, //1
		load: `hello`,
		left: `:   "world", "a": 1, "b": false}`,
	}, {
		kind: TokenCOL, //2
		load: `:`,
		left: `   "world", "a": 1, "b": false}`,
	}, {
		kind: TokenSTR, //3
		load: `world`,
		left: `, "a": 1, "b": false}`,
	}, {
		kind: TokenCOM, //4
		load: `,`,
		left: ` "a": 1, "b": false}`,
	}, {
		kind: TokenSTR, //5
		load: `a`,
		left: `: 1, "b": false}`,
	}, {
		kind: TokenCOL, //6
		load: `:`,
		left: ` 1, "b": false}`,
	}, {
		kind: TokenNUM, //7
		load: `1`,
		left: `, "b": false}`,
	}, {
		kind: TokenCOM, //8
		load: `,`,
		left: ` "b": false}`,
	}, {
		kind: TokenSTR, //9
		load: `b`,
		left: `: false}`,
	}, {
		kind: TokenCOL, //10
		load: `:`,
		left: ` false}`,
	}, {
		kind: TokenLIT, //11
		load: `false`,
		left: `}`,
	}, {
		kind: TokenRCB, //12
		load: `}`,
		left: ``,
	}}

	var i int
	y := func(kind TokenKind, load []byte, pos uint) bool {
		if i >= len(steps) {
			panic("once too often called")
		}

		step := steps[i]
		if !kind.Is(step.kind) {
			t.Errorf("%d: unexpected token %q", i, load)
		} else if string(load) != step.load {
			t.Errorf("%d: uexpected load %q", i, load)
		} else if r.String() != step.left {
			t.Errorf("%d: uexpected left content to parse (%d != %d) %q", i, r.Len(), len(step.left), r.String())
		}
		i++
		return false
	}

	l := NewLexer(y, LexerOptEnableUnreadBuffer)

	for range steps {
		l.Scan(r)
	}
}

// ensure LexerOptEnableUnreadBuffer is working with arrays
func TestLexer_Scan_11(t *testing.T) {
	s := []byte(`["hello", 1, false]`)
	r := bytes.NewBuffer(s)

	steps := []struct {
		kind TokenKind
		load string
		left string
	}{{
		kind: TokenLSB, //0
		load: `[`,
		left: `"hello", 1, false]`,
	}, {
		kind: TokenSTR, //1
		load: `hello`,
		left: `, 1, false]`,
	}, {
		kind: TokenCOM, //2
		load: `,`,
		left: ` 1, false]`,
	}, {
		kind: TokenNUM, //3
		load: `1`,
		left: `, false]`,
	}, {
		kind: TokenCOM, //4
		load: `,`,
		left: ` false]`,
	}, {
		kind: TokenLIT, //5
		load: `false`,
		left: `]`,
	}, {
		kind: TokenRSB, //6
		load: `]`,
		left: ``,
	}}

	var i int
	y := func(kind TokenKind, load []byte, pos uint) bool {
		if i >= len(steps) {
			panic("once too often called")
		}

		step := steps[i]
		if !kind.Is(step.kind) {
			t.Errorf("%d: unexpected token %q", i, load)
		} else if string(load) != step.load {
			t.Errorf("%d: uexpected load %q", i, load)
		} else if r.String() != step.left {
			t.Errorf("%d: uexpected left content to parse (%d != %d) %q", i, r.Len(), len(step.left), r.String())
		}
		i++
		return false
	}

	l := NewLexer(y, LexerOptEnableUnreadBuffer)

	for range steps {
		l.Scan(r)
	}
}
