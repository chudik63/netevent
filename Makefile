proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  ./pkg/proto/auth.proto
build:
	go build -o cmd/main cmd/main.go
send:
	scp -P443 -r -B . root@94.159.99.214:/root/src/go/temp/lms/02/netevent
tarsend:
	tar -zcf 02.gzip ../auth/
	scp -P443 02.gzip root@94.159.99.214:/root/src/go/temp/lms/02/netevent
update:
	docker cp auth/cmd/main auth_app:/app
.PHONY: proto build send