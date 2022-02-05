// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 10/2020
//               Gregor Noczinski · gregor [at] noczinski [dot] eu · 12/2021

package jsonlex

import (
	"errors"
	"fmt"
	"io"
)

type (
	// Lexer splits JSON byte stream into tokens.
	Lexer struct {
		yield Yield   // callback function
		area  []byte  // pre-allocated space
		buff  [1]byte // read-in buffer
		bpos  uint    // byte position in stream
		tpos  uint    // token position in stream
		hold  bool    // whether to advance reader
		frac  bool    // number fraction mode
		expo  bool    // number exponent mode
		sign  bool    // exponent sign
		esc   bool    // string escaping mode
		burd  bool    // is true if buffer was unread
		burde bool    // unread feature (if supported) enabled

	}

	// Yield is a callback function. It will be invoked
	// by the Scan() function for each token to be found.
	// The callback must return whether the scan process
	// should continue (true) or stop (false).
	Yield func(kind TokenKind, load []byte, pos uint) bool
)

// NewLexer takes a callback (yield) function as parameter.
// This yield function will be invoked for each token
// consumed from the byte stream by Scan().
func NewLexer(yield Yield, opts ...lexerOpt) *Lexer {
	l := &Lexer{
		yield: yield,
		area:  make([]byte, 0, 1024),
		buff:  [1]byte{0},
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

type lexerOpt func(*Lexer)

var (
	// LexerOptEnableUnreadBuffer enables if the given io.Reader
	// does implement UnreadableReader; it's UnreadableReader.UnreadByte
	// will be called if the Lexer reads one byte more to ensure that
	// a literal or number was ended. This ensures this Lexer never reads
	// more bytes than it is currently processing.
	LexerOptEnableUnreadBuffer lexerOpt = func(l *Lexer) {
		l.burde = true
	}
)

// Scan reads and tokenizes the byte stream.
// The yield function is invoked for each token found.
//
// The Scan() function terminates in following cases:
//   a) when the yield function return false
//   b) after emitting a jsonlex.TokenEOF or jsonlex.TokenERR
//
// Important: The Scan() function is reentrant, subsequent invocations will
// continue to consume the available byte stream as long as you provide
// a reader that implements an UnreadByte() interface, and you configure
// the Lexer with the LexerOptEnableUnreadBuffer option activated.
func (l *Lexer) Scan(r io.Reader) {
	var (
		b   byte      // byte under scrutiny
		n   int       // number of bytes read
		t   TokenKind // current token or state
		err error     // ordinary error holder
	)

nextToken:
	l.esc, l.frac = false, false
	l.expo, l.sign = false, false
	load := l.area[:0]
	t = scanning

	if l.burd {
		if _, ok := r.(UnreadableReader); ok {
			expb := l.buff[0]
			n, err = r.Read(l.buff[:])
			if n == 0 || err == io.EOF {
				err = io.ErrUnexpectedEOF
				goto emitErrToken
			}
			if err != nil {
				goto emitErrToken
			}
			if expb != l.buff[0] {
				err = errUnexpectedByte
			}
			l.burd = false
		}
	}

nextByte:
	if l.hold {
		l.hold = false
	} else {
		n, err = r.Read(l.buff[:])
		l.burd = false
		l.bpos += uint(n)
	}

	if err != nil {
		if err == io.EOF && len(load) > 0 && t.Is(TokenNUM) {
			goto emitNumToken
		}
		if err == io.EOF && len(load) > 0 && t.Is(TokenLIT) {
			goto emitLitToken
		}
		if err == io.EOF && len(load) > 0 {
			goto emitToken
		}
		if err == io.EOF {
			l.tpos = l.bpos
			goto emitEofToken
		}
		goto emitErrToken
	}

	if b = l.buff[0]; t != scanning {
		if t.Is(TokenSTR) {
			goto scanStr
		}
		if t.Is(TokenNUM) {
			goto scanNum
		}
		if t.Is(TokenLIT) {
			goto scanLit
		}
		goto emitTokenHold
	}

	if b == 0x20 || b == '\n' || b == '\r' || b == '\t' {
		goto nextByte
	}
	if b > 0x7F || b < 0x20 {
		goto emitUnexpErrToken
	}

	l.tpos = l.bpos - 1
	if s := states[b]; s != 0 {
		t = s
		if b == '"' {
			goto nextByte
		}
		goto consume
	}

emitUnexpErrToken:
	if m := fmt.Sprintf("unexpected %q (0x%X)", b, b); true {
		l.yield(TokenERR, []byte(m), l.tpos)
	}
	return

emitErrToken:
	l.yield(TokenERR, []byte(err.Error()), l.tpos)
	return

emitEofToken:
	l.yield(TokenEOF, nil, l.tpos)
	return

emitNumToken:
	if b := load[len(load)-1]; b == '.' || b == '-' ||
		b == 'e' || b == 'E' {
		goto emitUnexpErrToken
	}
	if len(load) >= 3 {
		if s := string(load[:3]); s == "-.e" || s == "-.E" {
			goto emitUnexpErrToken
		}
	}
	goto emitTokenHold

emitLitToken:
	if s := string(load); true {
		if s != "null" && s != "true" && s != "false" {
			goto emitUnexpErrToken
		}
	}

emitTokenHold:
	if err == nil {
		l.hold = true
	}

emitToken:
	if l.yield(t, load, l.tpos) {
		goto nextToken
	}
	return

scanStr:
	if l.esc {
		l.esc = false
		goto consume
	} else if b == '\\' {
		l.esc = true
	}
	if !l.esc && b == '"' {
		goto emitToken
	}
	goto consume

scanNum:
	if b >= '0' && b <= '9' {
		l.sign = false
		goto consume
	}
	if !l.frac && b == '.' {
		l.frac = true
		goto consume
	}
	if !l.expo && (b == 'e' || b == 'E') {
		l.frac, l.expo, l.sign = true, true, true
		goto consume
	}
	if l.sign && (b == '+' || b == '-') {
		l.sign = false
		goto consume
	}

	if l.burde {
		if ur, ok := r.(UnreadableReader); ok {
			if err = ur.UnreadByte(); err != nil {
				goto emitErrToken
			}
			l.burd = true
		}
	}

	goto emitNumToken

scanLit:
	if b >= 'a' && b <= 'z' {
		goto consume
	}

	if l.burde {
		if ur, ok := r.(UnreadableReader); ok {
			if err = ur.UnreadByte(); err != nil {
				goto emitErrToken
			}
			l.burd = true
		}
	}

	goto emitLitToken

consume:
	load = append(load, b)
	switch t {
	case TokenLSB, TokenRSB,
		TokenLCB, TokenRCB,
		TokenCOL, TokenCOM:
		goto emitToken
	}
	goto nextByte
}

var states = [0x80]TokenKind{
	' ':  0,        // 0x20 space
	'!':  0,        // 0x21 exclamation mark
	'"':  TokenSTR, // 0x22 quotation mark
	'#':  0,        // 0x23 number sign
	'$':  0,        // 0x24 dollar sign
	'%':  0,        // 0x25 percent sign
	'&':  0,        // 0x26 ampersand
	'\'': 0,        // 0x27 apostrophe
	'(':  0,        // 0x28 left parenthesis
	')':  0,        // 0x29 right parenthesis
	'*':  0,        // 0x2A asterisk
	'+':  0,        // 0x2B plus sign
	',':  TokenCOM, // 0x2C comma
	'-':  TokenNUM, // 0x2D minus sign
	'.':  0,        // 0x2E full stop
	'/':  0,        // 0x2F forward slash
	'0':  TokenNUM, // 0x30 digit
	'1':  TokenNUM, // 0x31 digit
	'2':  TokenNUM, // 0x32 digit
	'3':  TokenNUM, // 0x33 digit
	'4':  TokenNUM, // 0x34 digit
	'5':  TokenNUM, // 0x35 digit
	'6':  TokenNUM, // 0x36 digit
	'7':  TokenNUM, // 0x37 digit
	'8':  TokenNUM, // 0x38 digit
	'9':  TokenNUM, // 0x39 digit
	':':  TokenCOL, // 0x3A colon
	';':  0,        // 0x3B semicolon
	'<':  0,        // 0x3C less-than sign
	'=':  0,        // 0x3D equals sign
	'>':  0,        // 0x3E greater-than sign
	'?':  0,        // 0x3F question mark
	'@':  0,        // 0x40 commercial at
	'A':  0,        // 0x41 capital letter
	'B':  0,        // 0x42 capital letter
	'C':  0,        // 0x43 capital letter
	'D':  0,        // 0x44 capital letter
	'E':  0,        // 0x45 capital letter
	'F':  0,        // 0x46 capital letter
	'G':  0,        // 0x47 capital letter
	'H':  0,        // 0x48 capital letter
	'I':  0,        // 0x49 capital letter
	'J':  0,        // 0x4A capital letter
	'K':  0,        // 0x4B capital letter
	'L':  0,        // 0x4C capital letter
	'M':  0,        // 0x4D capital letter
	'N':  0,        // 0x4E capital letter
	'O':  0,        // 0x4F capital letter
	'P':  0,        // 0x50 capital letter
	'Q':  0,        // 0x51 capital letter
	'R':  0,        // 0x52 capital letter
	'S':  0,        // 0x53 capital letter
	'T':  0,        // 0x54 capital letter
	'U':  0,        // 0x55 capital letter
	'V':  0,        // 0x56 capital letter
	'W':  0,        // 0x57 capital letter
	'X':  0,        // 0x58 capital letter
	'Y':  0,        // 0x59 capital letter
	'Z':  0,        // 0x5A capital letter
	'[':  TokenLSB, // 0x5B left square bracket
	'\\': 0,        // 0x5C reverse slash
	']':  TokenRSB, // 0x5D right square bracket
	'^':  0,        // 0x5E circumflex accent
	'_':  0,        // 0x5F low line
	'`':  0,        // 0x60 grave accent
	'a':  0,        // 0x61 small letter
	'b':  0,        // 0x62 small letter
	'c':  0,        // 0x63 small letter
	'd':  0,        // 0x64 small letter
	'e':  0,        // 0x65 small letter
	'f':  TokenLIT, // 0x66 small letter
	'g':  0,        // 0x67 small letter
	'h':  0,        // 0x68 small letter
	'i':  0,        // 0x69 small letter
	'j':  0,        // 0x6A small letter
	'k':  0,        // 0x6B small letter
	'l':  0,        // 0x6C small letter
	'm':  0,        // 0x6D small letter
	'n':  TokenLIT, // 0x6E small letter
	'o':  0,        // 0x6F small letter
	'p':  0,        // 0x70 small letter
	'q':  0,        // 0x71 small letter
	'r':  0,        // 0x72 small letter
	's':  0,        // 0x73 small letter
	't':  TokenLIT, // 0x74 small letter
	'u':  0,        // 0x75 small letter
	'v':  0,        // 0x76 small letter
	'w':  0,        // 0x77 small letter
	'x':  0,        // 0x78 small letter
	'y':  0,        // 0x79 small letter
	'z':  0,        // 0x7A small letter
	'{':  TokenLCB, // 0x7B left curly brace
	'|':  0,        // 0x7C vertical line
	'}':  TokenRCB, // 0x7D right curly brace
	'~':  0,        // 0x7E tilde
	0x7F: 0,        // 0x7F unexpected character
}

// UnreadableReader provides the UnreadByte
type UnreadableReader interface {
	io.Reader

	// UnreadByte unreads the last read byte by this reader.
	UnreadByte() error
}

var errUnexpectedByte = errors.New("unexpected byte")
