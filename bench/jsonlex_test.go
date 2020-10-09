// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 10/2020

package bench

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/dtgorski/jsonlex"
)

func Benchmark_dtgorski_jsonlex_2kB(b *testing.B) {
	runLexer(b, "../testdata/2kB.json")
}

func Benchmark_dtgorski_jsonlex_20kB(b *testing.B) {
	runLexer(b, "../testdata/20kB.json")
}

func Benchmark_dtgorski_jsonlex_200kB(b *testing.B) {
	runLexer(b, "../testdata/200kB.json")
}

func Benchmark_dtgorski_jsonlex_2000kB(b *testing.B) {
	runLexer(b, "../testdata/2000kB.json")
}

func runLexer(b *testing.B, file string) {
	f, _ := os.Open(file)
	defer func() { _ = f.Close() }()
	buf, _ := ioutil.ReadAll(f)
	r := newReader(buf)

	i := 0
	lexer := jsonlex.NewLexer(
		func(kind jsonlex.Token, load []byte) { i &= 0 },
	)
	for n := 0; n < b.N; n++ {
		r.Reset()
		lexer.Scan(r)
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
