build:
	@go build -o bin/golang-json-server

run: build
	@./bin/golang-json-server

test:
	@go test -v ./...

postgres:
	@docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=12345 -d postgres:16-alpine