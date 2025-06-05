set positional-arguments
set export

#
# - Task prefix "@": don't echo commands
# - Task prefix "_": don't advertise task in the user list
# - Tasks run each line in a new shell; use a shebang to override
# - Below default shows users the task list in file order
#
@_default:
	just --list --unsorted

#
# Everyday work
#

# Run once
@run *args='': gen
	go run -v . "$@"

# Live reload
@live *args='': gen
	eval `go env`; $GOPATH/bin/reflex -d none -s -- just run

@test:
	# Specifying a test count overrides test caching
	# go test -count 1 ./test/...

@lint:
	eval `go env`; $GOPATH/bin/golangci-lint run -c .github/.golangci.yaml

@fmt:
	gofmt -w .

# Pre-commit sanity checks
pre: gen fmt lint test

# Update generated code
@gen:
	# go generate ./ent

# Browse local documentation
@doc:
	@echo "\nBrowse to http://localhost:6060\n"
	eval `go env`; $GOPATH/bin/pkgsite -http localhost:6060

# ORM CLI
@ent *args='':
	# Important: match against go.mod; avoids drift
	go run -mod=mod entgo.io/ent/cmd/ent "$@"

# Developer setup
setup: _precache gen

#
# Util
#

# For convenience + memory: update all libraries
_update:
	go get -u
	go mod tidy

# Warm compilation & module cache. Useful for containers, CI, etc
_precache:
	#!/usr/bin/env sh

	# Third-party dependencies that should be precompiled, with blacklist
	# Note double open-brackets to escape Justfile syntax
	precache_dep_filter='{{{{if not (or .Indirect .Main)}}{{{{.Path}}{{{{end}}'
	precache_skip_modules="golang.org/x/crypto|go.opentelemetry.io/otel/sdk"

	go mod download
	go build -v std
	go list -m -f "$precache_dep_filter" all | (grep -Ev "$precache_skip_modules" || true) | xargs -r -- go build -v
	echo {{linter}} {{pkgsite}} {{gopls}} {{gotestsum}} {{reflex}} | xargs -n 1 -- go install -v

# Tool versions
linter    := "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6"
pkgsite   := "golang.org/x/pkgsite/cmd/pkgsite@latest"
gopls     := "golang.org/x/tools/gopls@latest"
gotestsum := "gotest.tools/gotestsum@latest"
reflex    := "github.com/cespare/reflex@latest"

#
# CI
#

# CI target
_ci: _ci-check-fmt lint _ci-test

# Confirm code is formatted
_ci-check-fmt:
	test -z "$(gofmt -l .)" || (echo "One or more lines need formatting"; exit 1)

# Test with junit XML
_ci-test:
	eval `go env`; $GOPATH/bin/gotestsum --junitfile report.xml --format testname

# Cut binary. Strip is faster & simpler than UPX
_release: gen
	rm -f "$BinaryName"
	go build -v -o  .
	strip --strip-debug --strip-unneeded "$BinaryName"
