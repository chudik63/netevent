proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  ./pkg/proto/auth.proto
build:
	go build -o cmd/main cmd/main.go


asdf:
	sadf

.PHONY: proto