PKG = github.com/sapcc/limesctl
ifeq ($(shell uname -s),Darwin)
	PREFIX := /usr/local
else
	PREFIX := /usr
endif

VERSION := $(shell scripts/find_version.sh)

GO            := GOBIN=$(CURDIR)/build go
GO_BUILDFLAGS :=
GO_LDFLAGS    := -s -w -X main.version=$(VERSION)

# Which packages to test with `go test`?
GO_TESTPKGS   := $(shell $(GO) list $(GO_BUILDFLAGS) -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)
# Which packages to measure test coverage for?
GO_COVERPKGS  := $(shell $(GO) list $(GO_BUILDFLAGS) $(PKG)/internal/...)
# Output files from `go test`.
GO_COVERFILES := $(patsubst %,build/%.cover.out,$(subst /,_,$(GO_TESTPKGS)))

# This is needed for substituting spaces with commas.
space := $(null) $(null)
comma := ,

all: build/limesctl

build/limesctl: FORCE
	$(GO) install -v $(GO_BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' '$(PKG)'

install: FORCE build/limesctl
	install -D -m 0755 build/limesctl "$(DESTDIR)$(PREFIX)/bin/limesctl"

lint: FORCE
	@printf "\e[1;34m>> golangci-lint\e[0m\n"
	@command -v golangci-lint >/dev/null 2>&1 || { echo >&2 "Error: golangci-lint is not installed. See: https://golangci-lint.run/usage/install/"; exit 1; }
	golangci-lint run

# Run unit tests
test: FORCE build/limesctl
	@printf "\e[1;34m>> go test\e[0m\n"
	$(GO) test $(GO_BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' ./...

# Run all checks
check: FORCE build/limesctl lint build/cover.html
	@printf "\e[1;32m>> All checks successful\e[0m\n"

# Test with coverage
build/%.cover.out: FORCE
	@printf "\e[1;34m>> go test $(subst _,/,$*)\e[0m\n"
	$(GO) test $(GO_BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' -coverprofile=$@ -covermode=count -coverpkg=$(subst $(space),$(comma),$(GO_COVERPKGS)) $(subst _,/,$*)
build/cover.out: $(GO_COVERFILES)
	$(GO) run $(GO_BUILDFLAGS) tools/gocovcat.go $(GO_COVERFILES) > $@
build/cover.html: build/cover.out
	$(GO) tool cover -html $< -o $@

clean: FORCE
	rm -rf -- build/*

tidy-deps: FORCE
	$(GO) mod tidy -v
	$(GO) mod verify

.PHONY: FORCE
