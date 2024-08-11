lint:
	golangci-lint run

fmt:
	gofumpt -w -l -extra .

proto:
	cd internal/grpc ; protoc --go_out=. --go_grpc_out=. ./*.proto
