doc:
	docker compose up
bapp:
	go build -o auth/cmd/main auth/cmd/main.go
update:
	docker cp auth/cmd/main auth_app:/app
.PHONY: doc