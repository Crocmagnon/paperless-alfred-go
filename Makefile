.PHONY: build workflow clean
build:
	go mod tidy
	GOOS=darwin GOARCH=amd64 go build -o ./dist/ppl-go-darwin-amd64 ./
	GOOS=darwin GOARCH=arm64 go build -o ./dist/ppl-go-darwin-arm64 ./
	lipo -create -output dist/ppl-go dist/ppl-go-darwin-amd64 dist/ppl-go-darwin-arm64

workflow: build
	./scripts/make-workflow

install: workflow
	open dist/paperless-alfred-go.alfredworkflow

clean:
	rm -rf dist
