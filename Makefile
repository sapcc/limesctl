ifeq ($(shell uname -s),Darwin)
	PREFIX  := /usr/local
else
	PREFIX  := /usr
endif
PKG      = github.com/sapcc/limesctl
VERSION := $(shell util/find_version.sh)

GO          := GOPATH=$(CURDIR)/.gopath GOBIN=$(CURDIR)/build go
BUILD_FLAGS :=
LD_FLAGS    := -s -w -X main.version=$(VERSION)

ifndef GOOS
	GOOS := $(word 1, $(subst /, " ", $(word 4, $(shell go version))))
endif
BINARY64  := limesctl-$(GOOS)_amd64
RELEASE64 := limesctl-$(VERSION)-$(GOOS)_amd64

################################################################################

all: build/limesctl

# This target uses the incremental rebuild capabilities of the Go compiler to speed things up.
# If no source files have changed, `go install` exits quickly without doing anything.
build/limesctl: FORCE
	$(GO) install $(BUILD_FLAGS) -ldflags '$(LD_FLAGS)' '$(PKG)'

install: FORCE all
	install -d -m 0755 "$(DESTDIR)$(PREFIX)/bin"
	install -m 0755 build/limesctl "$(DESTDIR)$(PREFIX)/bin/limesctl"

ifeq ($(GOOS),windows)
release: FORCE release/$(BINARY64)
	cd release && cp -f $(BINARY64) limesctl.exe && zip $(RELEASE64).zip limesctl.exe
	cd release && rm -f limesctl.exe
else
release: FORCE release/$(BINARY64)
	cd release && cp -f $(BINARY64) limesctl && tar -czf $(RELEASE64).tgz limesctl
	cd release && rm -f limesctl
endif

release-all: FORCE clean
	GOOS=darwin make release
	GOOS=linux  make release

release/$(BINARY64): FORCE
	GOARCH=amd64 $(GO) build $(BUILD_FLAGS) -o $@ -ldflags '$(LD_FLAGS)' '$(PKG)'

################################################################################

# which packages to test with static checkers?
GO_ALLPKGS := $(PKG) $(shell go list $(PKG)/pkg/...)
# which packages to test with `go test`?
GO_TESTPKGS := $(shell go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' $(PKG)/pkg/...)
# which packages to measure coverage for?
GO_COVERPKGS := $(shell go list $(PKG)/pkg/...)
# output files from `go test`
GO_COVERFILES := $(patsubst %,build/%.cover.out,$(subst /,_,$(GO_TESTPKGS)))

# down below, I need to substitute spaces with commas; because of the syntax,
# I have to get these separators from variables
space := $(null) $(null)
comma := ,

check: all static-check build/cover.html FORCE
	@echo -e "\e[1;32m>> All tests successful.\e[0m"
static-check: FORCE
	@if ! hash golint 2>/dev/null; then echo ">> Installing golint..."; go get -u golang.org/x/lint/golint; fi
	@echo '>> gofmt'
	@if s="$$(gofmt -s -l *.go pkg 2>/dev/null)" && test -n "$$s"; then printf ' => %s\n%s\n' gofmt  "$$s"; false; fi
	@echo '>> golint'
	@if s="$$(golint . && find pkg -type d ! -name dbdata -exec golint {} \; 2>/dev/null)" && test -n "$$s"; then printf ' => %s\n%s\n' golint "$$s"; false; fi
	@echo '>> go vet'
	@$(GO) vet -composites=false $(GO_ALLPKGS)
build/%.cover.out: FORCE
	@echo '>> go test $*'
	$(GO) test $(GO_BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' -coverprofile=$@ -covermode=count -coverpkg=$(subst $(space),$(comma),$(GO_COVERPKGS)) $(subst _,/,$*)
build/cover.out: $(GO_COVERFILES)
	util/gocovcat/main.go $(GO_COVERFILES) > $@
build/cover.html: build/cover.out
	$(GO) tool cover -html $< -o $@

clean: FORCE
	rm -rf build release

.PHONY: FORCE
