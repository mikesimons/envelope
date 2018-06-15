version=$(shell git describe --all | sed -e's/.*\///g')

.PHONY: test dev-deps

all: test cmd/envelope/envelope-$(version)-linux-amd64 cmd/envelope/envelope-$(version)-darwin-amd64 cmd/envelope/envelope-$(version)-windows-amd64

cmd/envelope/envelope-$(version)-linux-amd64:
	GOARCH=amd64 GOOS=linux cd cmd/envelope && go build -o $@

cmd/envelope/envelope-$(version)-darwin-amd64:
	GOARCH=amd64 GOOS=darwin cd cmd/envelope && go build -o $@

cmd/envelope/envelope-$(version)-windows-amd64:
	GOARCH=amd64 GOOS=windows cd cmd/envelope && go build -o $@

dev-deps:
	go get github.com/Masterminds/glide
	go get github.com/alecthomas/gometalinter
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega/...

test:
	ginkgo ./...
