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
	CGO_ENABLED=0 GOMAXPROCS=1 GOGC=off go test -run=^$$ -bench=. ./bench

prof-cpu: clean         # Creates CPU profiler output
	CGO_ENABLED=0 GOMAXPROCS=1 GOGC=off go test -cpuprofile=cpu.prof -bench=cursor.*2000kB ./bench
	@echo "\nCPU --------------------------------------"
	go tool pprof -top cpu.prof | sed "s/^/    /"

prof-mem: clean        # Creates memory profiler output
	CGO_ENABLED=0 GOMAXPROCS=1 GOGC=off go test -memprofile=mem.prof -bench=cursor.*2000kB ./bench
	@echo "\nMEM --------------------------------------"
	go tool pprof -top mem.prof | sed "s/^/    /"

sniff:                  # Checks format and runs linter (void on success)
	@find . -type f -not -path "*/\.*" -name "*.go" | xargs -I{} gofmt -d {}
	@go vet ./... || true
	@>/dev/null which revive || (echo "Missing a linter, install with:  go install github.com/mgechev/revive" && false)
	@revive -config .revive.toml ./...

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
