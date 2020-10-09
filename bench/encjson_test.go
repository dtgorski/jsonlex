// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 10/2020

package bench

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func Benchmark_encoding_json_2kB(b *testing.B) {
	runDecoder(b, "../testdata/2kB.json")
}

func Benchmark_encoding_json_20kB(b *testing.B) {
	runDecoder(b, "../testdata/20kB.json")
}

func Benchmark_encoding_json_200kB(b *testing.B) {
	runDecoder(b, "../testdata/200kB.json")
}

func Benchmark_encoding_json_2000kB(b *testing.B) {
	runDecoder(b, "../testdata/2000kB.json")
}

func runDecoder(t *testing.B, file string) {
	f, _ := os.Open(file)
	defer func() { _ = f.Close() }()
	b, _ := ioutil.ReadAll(f)

	for n := 0; n < t.N; n++ {
		d := json.NewDecoder(bytes.NewReader(b))
		for {
			_, err := d.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Error(err)
			}
		}
	}
}
