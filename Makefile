build:
	@go build -o bin/golang-json-server

run: build
	@./bin/golang-json-server

test:
	@go test -v ./...