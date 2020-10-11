[![Build Status](https://travis-ci.org/dtgorski/jsonlex.svg?branch=master)](https://travis-ci.org/dtgorski/jsonlex)
[![Coverage Status](https://coveralls.io/repos/github/dtgorski/jsonlex/badge.svg?branch=master)](https://coveralls.io/github/dtgorski/jsonlex?branch=master)
[![Open Issues](https://img.shields.io/github/issues/dtgorski/jsonlex.svg)](https://github.com/dtgorski/jsonlex/issues)
[![Report Card](https://goreportcard.com/badge/github.com/dtgorski/jsonlex)](https://goreportcard.com/report/github.com/dtgorski/jsonlex)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dtgorski/jsonlex)](https://pkg.go.dev/github.com/dtgorski/jsonlex)

## jsonlex

Fast JSON lexer (tokenizer) with no memory footprint and no garbage collector pressure (zero-alloc).

### Installation
```
go get -u github.com/dtgorski/jsonlex
```

### Usage
```
package main

import (
    "bytes"
    "github.com/dtgorski/jsonlex"
)

func main() {
    reader := bytes.NewReader(
        []byte(`{ "foo": "bar", "baz": 42 }`),
    )

    lexer := jsonlex.NewLexer(
        func(token jsonlex.Token, load []byte, pos uint) {

            save := make([]byte, len(load))
            copy(save, load)

            println(pos, token, string(save))
        },
    )

    lexer.Scan(reader)
}
```

### Emitted tokens
| [```jsonlex```](https://pkg.go.dev/github.com/dtgorski/jsonlex) | Representation
| --- | ---
|```TokenEOF``` | signals end of file/stream
|```TokenERR``` | error string (other than EOF)
|```TokenLIT``` | literal (```true```, ```false```, ```null```)
|```TokenNUM``` | float number
|```TokenSTR``` | "...\\"..."
|```TokenCOL``` | : colon
|```TokenCOM``` | , comma
|```TokenLSB``` | [ left square bracket
|```TokenRSB``` | ] right square bracket
|```TokenLCB``` | { left curly brace
|```TokenRCB``` | } right curly brace

### Artificial benchmarks

Each benchmark consists of complete tokenization of a JSON document of a given size (2kB, 20kB, 200kB and 2000kB). The unit ```doc/s``` means _tokenized documents per second_, so more is better. 
The comparison candidate is Go's [encoding/json.Decoder.Token()](https://golang.org/pkg/encoding/json/#Decoder.Token) implementation.

| |2kB|20kB|200kb|2000kB
| --- | --- | --- | --- | ---
|```encoding/json```|```10987 doc/s```|```1184 doc/s```|```128 doc/s```|```13 doc/s```
|```dtgorski/jsonlex```|**```59346 doc/s```**|**```6021 doc/s```**|**```615 doc/s```**|**```68 doc/s```**

```
goos: linux
goarch: amd64
pkg: github.com/dtgorski/jsonlex/bench

Benchmark_encoding_json_2kB-8        10987     109031 ns/op      36528 B/op      1963 allocs/op
Benchmark_encoding_json_20kB-8        1184    1025208 ns/op     318434 B/op     18231 allocs/op
Benchmark_encoding_json_200kB-8        128    9484296 ns/op    2877981 B/op    164401 allocs/op
Benchmark_encoding_json_2000kB-8        13   78722997 ns/op   23356024 B/op   1319126 allocs/op

Benchmark_dtgorski_jsonlex_2kB-8     59346      20237 ns/op          0 B/op         0 allocs/op
Benchmark_dtgorski_jsonlex_20kB-8     6021     199091 ns/op          0 B/op         0 allocs/op
Benchmark_dtgorski_jsonlex_200kB-8     615    1944173 ns/op          0 B/op         0 allocs/op
Benchmark_dtgorski_jsonlex_2000kB-8     68   17415371 ns/op          0 B/op         0 allocs/op
```

### Disclaimer
The implementation and features of ```jsonlex``` follow the [YAGNI](https://en.wikipedia.org/wiki/You_aren%27t_gonna_need_it) principle.
There is no claim for completeness or reliability.

### @dev
Try ```make```:
```
$ make

 make help       Displays this list
 make clean      Removes build/test artifacts
 make test       Runs integrity test with -race
 make bench      Executes artificial benchmarks
 make prof-cpu   Creates CPU profiler output
 make prof-mem   Creates memory profiler output
 make sniff      Checks format and runs linter (void on success)
 make tidy       Formats source files, cleans go.mod
```

### License
[MIT](https://opensource.org/licenses/MIT) - © dtg [at] lengo [dot] org
