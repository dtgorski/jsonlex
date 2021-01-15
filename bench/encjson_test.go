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

func Benchmark_encjson_2kB(b *testing.B) {
	runDecoder(b, "../testdata/2kB.json")
}

func Benchmark_encjson_20kB(b *testing.B) {
	runDecoder(b, "../testdata/20kB.json")
}

func Benchmark_encjson_200kB(b *testing.B) {
	runDecoder(b, "../testdata/200kB.json")
}

func Benchmark_encjson_2000kB(b *testing.B) {
	runDecoder(b, "../testdata/2000kB.json")
}

func runDecoder(b *testing.B, file string) {
	b.ReportAllocs()

	f, _ := os.Open(file)
	defer func() { _ = f.Close() }()
	buf, _ := ioutil.ReadAll(f)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		d := json.NewDecoder(bytes.NewReader(buf))
		for {
			_, err := d.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Error(err)
			}
		}
	}
}
