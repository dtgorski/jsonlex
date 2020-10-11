.PHONY: help clean test bench prof tidy sniff .travis

help:                   # Displays this list
	@echo; grep "^[a-z][a-zA-Z0-9_<> -]\+:" Makefile | sed -E "s/:[^#]*?#?(.*)?/\r\t\t\1/" | sed "s/^/ make /"; echo

clean:                  # Removes build/test artifacts
	@find . -type f | grep "\.out$$" | xargs -I{} rm {};
	@find . -type f | grep "\.html$$" | xargs -I{} rm {};
	@find . -type f | grep "\.test$$" | xargs -I{} rm {};
	@find . -type f | grep "\.prof$$" | xargs -I{} rm {};

test: clean             # Runs integrity test with -race
	CGO_ENABLED=1 go test -v -count=1 -race -covermode=atomic -coverprofile=./coverage.out .
	@go tool cover -html=./coverage.out -o ./coverage.html && echo "coverage: <file://$(PWD)/coverage.html>"

bench: clean            # Executes artificial benchmarks
	CGO_ENABLED=0 go test -benchmem -bench=. ./bench

prof-cpu: clean         # Creates CPU profiler output
	CGO_ENABLED=0 go test -cpuprofile=cpu.prof -bench=jsonlex.*2000kB ./bench
	@echo "\nCPU --------------------------------------"
	@go tool pprof -top cpu.prof | head -20 | sed "s/^/    /"
	@go tool pprof -weblist=. ./bench.test cpu.prof &

prof-mem: clean        # Creates memory profiler output
	CGO_ENABLED=0 go test -benchmem -memprofilerate=0 -memprofile=mem.prof -bench=jsonlex.*2000kB ./bench
	@echo "\nMEM --------------------------------------"
	@go tool pprof -top mem.prof | head -20 | sed "s/^/    /"
	@go tool pprof -weblist=. ./bench.test mem.prof &

sniff:                  # Checks format and runs linter (void on success)
	@gofmt -d .
	@2>/dev/null revive -config revive.toml ./... || echo "get a linter first:  go install github.com/mgechev/revive"

tidy:                   # Formats source files, cleans go.mod
	@gofmt -w .
	@go mod tidy

.travis:                # Travis CI (see .travis.yml), runs tests
    ifndef TRAVIS
	    @echo "Fail: requires Travis runtime"
    else
	    @$(MAKE) test --no-print-directory && \
	    goveralls -coverprofile=./coverage.out -service=travis-ci
    endif
