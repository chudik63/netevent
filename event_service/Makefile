protoc:
	protoc --go-grpc_out=pkg/api --go_out=pkg/api --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative ./proto/event/event.proto

up:
	docker-compose --env-file ./.env up