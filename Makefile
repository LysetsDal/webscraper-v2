app := bin/webscraper-v2

.PHONY: test
# ^ Otherwise make thinks `test` is a file and optimizes

build:
	@go build -o $(app) cmd/*.go

run: build
	@./$(app)

test:
	@go test -v ./test
