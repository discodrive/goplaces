# Adding @ before commands prevents them from being printed

build:
	@go build -o bin/web-service

run: build
	@./bin/web-service

test:
	@sgo test -v ./...
