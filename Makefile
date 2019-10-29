PKGS := $(shell go list ./... | grep -v /vendor)
VERSION := $(shell git describe --tags --always --dirty)
BINARY := tempo
BUILD_OS = $(shell go env GOHOSTOS)
BUILD_ARCH = $(shell go env GOARCH)
GOBIN = $(shell go env GOPATH)/bin

.PHONY: test
test:
	go test $(PKGS)

build: release
	echo $(GOPATH)
	cp release/$(BINARY)-$(VERSION)-$(BUILD_OS)-$(BUILD_ARCH) $(GOBIN)/$(BINARY)
	chmod +x $(GOPATH)/bin/$(BINARY)
	$(GOPATH)/bin/$(BINARY)


.PHONY: windows
windows:
	mkdir -p release
	GOOS=windows GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-windows-amd64

.PHONY: linux
linux:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-linux-amd64
.PHONY: darwin
darwin:
	mkdir -p release
	GOOS=darwin GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-darwin-amd64

.PHONY: release
release: windows linux darwin
