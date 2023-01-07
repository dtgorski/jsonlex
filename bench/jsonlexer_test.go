/* MIT license - Nicholas Hancock 1/2023 */

package bench

import (
	"bytes"
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/brimdata/zed/pkg/jsonlexer"
)

func Benchmark_jsonlexer_lexer_2kB(b *testing.B) {
	runLexer(b, "../testdata/2kB.json")
}

func Benchmark_jsonlexer_lexer_20kB(b *testing.B) {
	runLexer(b, "../testdata/20kB.json")
}

func Benchmark_jsonlexer_lexer_200kB(b *testing.B) {
	runLexer(b, "../testdata/200kB.json")
}

func Benchmark_jsonlexer_lexer_2000kB(b *testing.B) {
	runLexer(b, "../testdata/2000kB.json")
}

func runJsonLexer(b *testing.B, file string) {
	b.ReportAllocs()

	f, _ := os.Open(file)
	defer func() { _ = f.Close() }()
	buf, _ := ioutil.ReadAll(f)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		lexer := New(bufio.NewReader(bytes.NewReader(buf)))
		for ;; {
			t := lexer.Token()
			if t == TokenErr {
				if err := lexer.Err() ; err != io.EOF {
					b.Errorf("%s", err)
				}
				break
			}
		}
	}
}
