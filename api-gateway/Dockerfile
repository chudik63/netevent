FROM golang:1.23 AS build

COPY volume/.ssh /root/.ssh
COPY volume/.netrc /root/.netrc
COPY volume/.gitconfig /root/.gitconfig

WORKDIR /build

COPY go.mod go.sum ./     

RUN go mod download      

COPY . .       

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway ./cmd/main

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=build /build/gateway .

CMD ["./gateway"]
