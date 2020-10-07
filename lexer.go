// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 10/2020

package jsonlex

import (
	"fmt"
	"io"
)

type (
	// Lexer splits JSON byte stream into tokens.
	Lexer struct {
		yield func(token Token, load []byte)
	}

	// Token denotes the type of token.
	Token uint8
)

// Kinds of tokens emitted by the lexer.
const (
	TokenEOF Token = iota // signals end of file/stream
	TokenERR              // error string
	TokenLIT              // literal (true, false, null)
	TokenNUM              // float number
	TokenSTR              // "...\"..."
	TokenCOL              // : colon
	TokenCOM              // , comma
	TokenLSB              // [ left square bracket
	TokenRSB              // ] right square bracket
	TokenLCB              // { left curly brace
	TokenRCB              // } right curly brace

	sniffing
)

// NewLexer takes a callback (yield) function as parameter.
// This yield function will be invoked for each token
// consumed from the byte stream by Scan().
//
// Example:
//
//    reader := bytes.NewReader(...)
//
//    lexer := jsonlex.NewLexer(
//        func(token jsonlex.Token, load []byte) {
//            println(token, string(load))
//        },
//    )
//
//    lexer.Scan(reader)
//
func NewLexer(yield func(kind Token, load []byte)) *Lexer {
	return &Lexer{yield}
}

// Scan reads and tokenizes the byte stream.
// The 'yield' function is invoked for each token found.
//
// After (and only after) emitting a jsonlex.TokenEOF
// or jsonlex.TokenERR, the Scan() function terminates.
func (l *Lexer) Scan(r io.Reader) {
	var (
		b    byte   // byte under scrutiny
		load []byte // token payload
		tok  Token  // current token
		pos  uint   // byte position in stream

		hold bool // whether to advance reader
		frac bool // number fraction mode
		expo bool // number exponent mode
		sign bool // exponent sign

		esc bool  // string escaping mode
		eot bool  // end of delimited text
		err error // ordinary error holder
		buf = []byte{0}
	)

nextToken:
	frac, expo, sign = false, false, false
	eot, esc = false, false

	load = load[:0]
	tok = sniffing

nextByte:
	if hold {
		hold = false
	} else {
		_, err = r.Read(buf)
		pos++
	}

	if err != nil {
		if err == io.EOF && len(load) > 0 && tok == TokenNUM {
			goto emitNumToken
		}
		if err == io.EOF && len(load) > 0 && tok == TokenLIT {
			goto emitLitToken
		}
		if err == io.EOF && len(load) > 0 {
			goto emitToken
		}
		if err == io.EOF {
			goto emitEofToken
		}
		goto emitErrToken
	}

	if b = buf[0]; tok != sniffing {
		switch tok {
		case TokenSTR:
			goto scanStr
		case TokenNUM:
			goto scanNum
		case TokenLIT:
			goto scanLit
		default:
			goto emitToken
		}
	}

	if b == 0x20 || b == '\n' || b == '\r' || b == '\t' {
		goto nextByte
	}
	if b > 0x7F || b < 0x20 {
		goto emitUnexpErrToken
	}
	if s := states[b]; s != 0 {
		tok = s
		goto consume
	}

emitUnexpErrToken:
	err = fmt.Errorf(unexpected, b, b)

emitErrToken:
	l.yield(TokenERR, []byte(fmt.Sprintf(atposition, err.Error(), pos-1)))
	return

emitEofToken:
	l.yield(TokenEOF, nil)
	return

emitStrToken:
	load = load[1 : len(load)-1]
	goto emitToken

emitNumToken:
	if t := load[len(load)-1]; t == '.' || t == '-' || t == 'e' || t == 'E' {
		goto emitUnexpErrToken
	}
	if size := len(load); size == 2 && load[0] == '-' && load[1] == '.' ||
		size >= 3 && load[0] == '-' && load[1] == '.' && load[2] == 'e' {
		goto emitUnexpErrToken
	}
	goto emitToken

emitLitToken:
	if s := string(load); s != "null" && s != "true" && s != "false" {
		goto emitUnexpErrToken
	}

emitToken:
	hold = true
	l.yield(tok, load)
	goto nextToken

scanStr:
	if esc {
		esc = false
		goto consume
	}
	if b == '\\' {
		esc = true
		goto consume
	}
	if b != '"' && !eot {
		goto consume
	}
	if !eot {
		eot = true
		goto consume
	}
	goto emitStrToken

scanNum:
	if b >= '0' && b <= '9' {
		sign = false
		goto consume
	}
	if !frac && b == '.' {
		frac = true
		goto consume
	}
	if !expo && (b == 'e' || b == 'E') {
		frac, expo, sign = true, true, true
		goto consume
	}
	if sign && (b == '+' || b == '-') {
		sign = false
		goto consume
	}
	goto emitNumToken

scanLit:
	if b >= 'a' && b <= 'z' {
		goto consume
	}
	goto emitLitToken

consume:
	load = append(load, b)
	goto nextByte
}

const unexpected = "unexpected %q (0x%X)"
const atposition = "%s at position %d"

var states = [0x80]Token{
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
