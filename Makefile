build:
	@go build -o bin/fs

run: build
	@./bin/fs

test:
	@go clean -testcache && go test ./... -v