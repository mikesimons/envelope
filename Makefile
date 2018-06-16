repo=github.com/mikesimons/envelope
version=$(shell git describe --all --dirty --long | awk -F"-|/" '/^heads/ {print $$2 "@" substr($$4, 2) "!" $$5}; /^tags/ { print $$2 }')
build_args=-ldflags "-X main.envelope_version_string=$(version)" $(repo)/cmd/envelope

.PHONY: test dev-deps

all: test build

build: build-linux build-darwin build-windows 

build-linux: build/envelope-$(version)-linux-amd64
build/envelope-$(version)-linux-amd64:
	GOARCH=amd64 GOOS=linux go build -o $@ $(build_args)

build-darwin: build/envelope-$(version)-darwin-amd64
build/envelope-$(version)-darwin-amd64:
	GOARCH=amd64 GOOS=darwin go build -o $@ $(build_args)

build-windows: build/envelope-$(version)-windows-amd64
build/envelope-$(version)-windows-amd64:
	GOARCH=amd64 GOOS=windows go build -o $@ $(build_args)

dev-deps:
	go get github.com/Masterminds/glide
	go get github.com/alecthomas/gometalinter
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega/...
	glide install

test:
	ginkgo ./...
