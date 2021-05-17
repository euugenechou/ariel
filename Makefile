build:
	@go generate

repl: build
	@go run main.go --repl=true

debug: build
	@go run main.go --debug=true --repl=true

test: build
	@go run main.go tests/lrc.arl

clean:
	@rm -f parser/parser.go
