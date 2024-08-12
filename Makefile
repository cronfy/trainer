lint:
	golangci-lint run

fmt:
	gofumpt -w -l -extra .

test:
	go test ./...

