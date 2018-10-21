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
release: release/$(BINARY64)
	cd release && cp -f $(BINARY64) limesctl.exe && zip $(RELEASE64).zip limesctl.exe
	cd release && rm -f limesctl.exe
else 
release: release/$(BINARY64)
	cd release && cp -f $(BINARY64) limesctl && tar -czf $(RELEASE64).tgz limesctl
	cd release && cp -f $(BINARY64) limesctl && tar -czf $(RELEASE64).tgz limesctl
	cd release && rm -f limesctl
endif

release-all: FORCE clean 
	GOOS=darwin  make release
	GOOS=linux   make release

release/$(BINARY64): FORCE
	GOARCH=amd64 $(GO) build $(BUILD_FLAGS) -o $@ -ldflags '$(LD_FLAGS)' '$(PKG)'

clean: FORCE
	rm -rf build release

.PHONY: FORCE
