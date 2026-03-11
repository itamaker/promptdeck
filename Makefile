BINARY := promptdeck

.PHONY: build test example release-check snapshot

build:
	mkdir -p dist
	go build -o dist/$(BINARY) .

test:
	go test ./...

example:
	go run . matrix -template examples/review.tmpl -matrix examples/matrix.json

release-check:
	goreleaser check

snapshot:
	goreleaser release --snapshot --clean

