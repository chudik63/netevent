FROM golang:1.23-alpine

WORKDIR /

COPY . .

RUN go mod download
RUN go build -o notification-server ./cmd/notification/main.go

CMD ["./notification-server"]