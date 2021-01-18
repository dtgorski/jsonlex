[![Build Status](https://travis-ci.org/dtgorski/jsonlex.svg?branch=master)](https://travis-ci.org/dtgorski/jsonlex)
[![Coverage Status](https://coveralls.io/repos/github/dtgorski/jsonlex/badge.svg?branch=master)](https://coveralls.io/github/dtgorski/jsonlex?branch=master)
[![Open Issues](https://img.shields.io/github/issues/dtgorski/jsonlex.svg)](https://github.com/dtgorski/jsonlex/issues)
[![Report Card](https://goreportcard.com/badge/github.com/dtgorski/jsonlex)](https://goreportcard.com/report/github.com/dtgorski/jsonlex)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dtgorski/jsonlex)](https://pkg.go.dev/github.com/dtgorski/jsonlex)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdtgorski%2Fjsonlex.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdtgorski%2Fjsonlex?ref=badge_shield)

## jsonlex

Fast JSON lexer (tokenizer) with no memory footprint and no garbage collector pressure (zero heap allocations).

### Installation
```
go get -u github.com/dtgorski/jsonlex
```

### Important
Using an ```io.Reader``` that directly uses system calls (e.g. ```os.File```) will result in poor performance. Wrap your input reader with ```bufio.Reader``` or better ```bytes.Reader``` to achieve best results.

### Usage A - iterating behaviour (Cursor)
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

    cursor := jsonlex.NewCursor(reader, nil)

    println(cursor.Curr().String())
    println(cursor.Next().String())

    if !cursor.Next().Is(jsonlex.TokenEOF) {
        println("there is more ...")
    }
}
```

### Usage B - emitting behaviour (Yield)
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
        func(kind jsonlex.TokenKind, load []byte, pos uint) bool {

            save := make([]byte, len(load))
            copy(save, load)

            println(pos, kind, string(save))
            return true
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

Each benchmark consists of complete tokenization of a JSON document of a given size (2kB, 20kB, 200kB and 2000kB) using one CPU core. The unit ```doc/s``` means _tokenized documents per second_, so more is better. 
The comparison candidate is Go's [encoding/json.Decoder.Token()](https://golang.org/pkg/encoding/json/#Decoder.Token) implementation.

| |2kB|20kB|200kb|2000kB
| --- | --- | --- | --- | ---
|```encoding/json```|```9910 doc/s```|```1152 doc/s```|```126 doc/s```|```14 doc/s```
|```dtgorski/jsonlex```|**```71880 doc/s```**|**```7341 doc/s```**|**```753 doc/s```**|**```85 doc/s```**

```
cpus: 1 core (~8000 BogoMIPS)
goos: linux
goarch: amd64
pkg: github.com/dtgorski/jsonlex/bench

Benchmark_encjson_2kB              9910     120475 ns/op      36528 B/op      1963 allocs/op
Benchmark_encjson_20kB             1152    1040771 ns/op     318432 B/op     18231 allocs/op
Benchmark_encjson_200kB             126    9494534 ns/op    2877968 B/op    164401 allocs/op
Benchmark_encjson_2000kB             14   77593586 ns/op   23355856 B/op   1319126 allocs/op

Benchmark_jsonlex_lexer_2kB       71880      16691 ns/op          0 B/op         0 allocs/op
Benchmark_jsonlex_lexer_20kB       7341     163210 ns/op          0 B/op         0 allocs/op
Benchmark_jsonlex_lexer_200kB       753    1594025 ns/op          0 B/op         0 allocs/op
Benchmark_jsonlex_lexer_2000kB       85   14107866 ns/op          0 B/op         0 allocs/op

Benchmark_jsonlex_cursor_2kB      38002      31776 ns/op       3680 B/op       592 allocs/op
Benchmark_jsonlex_cursor_20kB      4058     300490 ns/op      25168 B/op      5446 allocs/op
Benchmark_jsonlex_cursor_200kB      422    2777058 ns/op     248816 B/op     49141 allocs/op
Benchmark_jsonlex_cursor_2000kB      50   23559879 ns/op    2254896 B/op    396298 allocs/op
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


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdtgorski%2Fjsonlex.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdtgorski%2Fjsonlex?ref=badge_large)