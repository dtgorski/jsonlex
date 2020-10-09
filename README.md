[![Build Status](https://travis-ci.org/dtgorski/jsonlex.svg?branch=master)](https://travis-ci.org/dtgorski/jsonlex)
[![Coverage Status](https://coveralls.io/repos/github/dtgorski/jsonlex/badge.svg?branch=master)](https://coveralls.io/github/dtgorski/jsonlex?branch=master)
[![Open Issues](https://img.shields.io/github/issues/dtgorski/jsonlex.svg)](https://github.com/dtgorski/jsonlex/issues)
[![Report Card](https://goreportcard.com/badge/github.com/dtgorski/jsonlex)](https://goreportcard.com/report/github.com/dtgorski/jsonlex)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dtgorski/jsonlex)](https://pkg.go.dev/github.com/dtgorski/jsonlex)

## jsonlex

Fast JSON lexer (tokenizer) with low memory footprint and low garbage collector pressure.

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
    reader := bytes.NewReader(...)

    lexer := jsonlex.NewLexer(
        func(token jsonlex.Token, load []byte) {
            println(token, string(load))
        },
    )

    lexer.Scan(reader)
}
```

### Emitted tokens
| [```jsonlex```](https://pkg.go.dev/github.com/dtgorski/jsonlex) | Representation
| --- | ---
|```TokenEOF``` | signals end of file/stream
|```TokenERR``` | error string
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
|```encoding/json```|```10946 doc/s```|```1179 doc/s```|```128 doc/s```|```14 doc/s```
|```dtgorski/jsonlex```|**```57790 doc/s```**|**```5874 doc/s```**|**```602 doc/s```**|**```62 doc/s```**

```
goos: linux
goarch: amd64
pkg: github.com/dtgorski/jsonlex/bench

Benchmark_encoding_json_2kB-8        10946     109034 ns/op      36528 B/op      1963 allocs/op
Benchmark_encoding_json_20kB-8        1179    1016282 ns/op     318487 B/op     18231 allocs/op
Benchmark_encoding_json_200kB-8        128    9403946 ns/op    2882058 B/op    164401 allocs/op
Benchmark_encoding_json_2000kB-8        14   79604995 ns/op   23655504 B/op   1319127 allocs/op

Benchmark_dtgorski_jsonlex_2kB-8     57790      20546 ns/op        256 B/op        11 allocs/op
Benchmark_dtgorski_jsonlex_20kB-8     5874     204287 ns/op       2370 B/op        86 allocs/op
Benchmark_dtgorski_jsonlex_200kB-8     602    1989911 ns/op      23803 B/op       630 allocs/op
Benchmark_dtgorski_jsonlex_2000kB-8     62   17726113 ns/op     267372 B/op      3977 allocs/op
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
 make prof       Executes artificial benchmarks w/ profile info
 make sniff      Checks format and runs linter (void on success)
 make tidy       Formats source files, cleans go.mod
```

### License
[MIT](https://opensource.org/licenses/MIT) - © dtg [at] lengo [dot] org
