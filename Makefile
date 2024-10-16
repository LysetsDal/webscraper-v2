app := bin/webscraper-v2

build:
	@go build -o $(app) cmd/*.go

run: build
	@./$(app)
