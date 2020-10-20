// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 10/2020

package bench

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/dtgorski/jsonlex"
)

func Benchmark_jsonlex_lexer_2kB(b *testing.B) {
	runLexer(b, "../testdata/2kB.json")
}

func Benchmark_jsonlex_lexer_20kB(b *testing.B) {
	runLexer(b, "../testdata/20kB.json")
}

func Benchmark_jsonlex_lexer_200kB(b *testing.B) {
	runLexer(b, "../testdata/200kB.json")
}

func Benchmark_jsonlex_lexer_2000kB(b *testing.B) {
	runLexer(b, "../testdata/2000kB.json")
}

func Benchmark_jsonlex_cursor_2kB(b *testing.B) {
	runCursor(b, "../testdata/2kB.json")
}

func Benchmark_jsonlex_cursor_20kB(b *testing.B) {
	runCursor(b, "../testdata/20kB.json")
}

func Benchmark_jsonlex_cursor_200kB(b *testing.B) {
	runCursor(b, "../testdata/200kB.json")
}

func Benchmark_jsonlex_cursor_2000kB(b *testing.B) {
	runCursor(b, "../testdata/2000kB.json")
}

func runLexer(b *testing.B, file string) {
	f, _ := os.Open(file)
	defer func() { _ = f.Close() }()
	buf, _ := ioutil.ReadAll(f)

	r := newReader(buf)

	lexer := NewLexer(
		func(kind TokenKind, load []byte, pos uint) bool {
			if kind == TokenERR {
				b.Fatal(kind)
			}
			return true
		},
	)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.Reset()
		lexer.Scan(r)
	}
}

func runCursor(b *testing.B, file string) {
	f, _ := os.Open(file)
	defer func() { _ = f.Close() }()
	buf, _ := ioutil.ReadAll(f)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r := newReader(buf)
		cursor := NewCursor(r, nil)

		for ; ; cursor.Next() {
			if cursor.Curr().Is(TokenERR) {
				b.Errorf("%s", cursor.Curr().Load)
				break
			}
			if cursor.Curr().Is(TokenEOF) {
				break
			}
		}
	}
}

type (
	reader struct {
		buf []byte
		pos int
		len int
	}
)

func newReader(b []byte) *reader {
	return &reader{buf: b, len: len(b)}
}

func (r *reader) Read(b []byte) (n int, err error) {
	if r.pos == r.len {
		return 0, io.EOF
	}
	b[0] = r.buf[r.pos]
	r.pos++
	return 1, nil
}

func (r *reader) Reset() {
	r.pos = 0
}
