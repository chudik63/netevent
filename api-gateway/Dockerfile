FROM golang:1.23 AS build

ARG NETRC

ARG GIT_URL
ARG GIT_SSH_URL

ARG SSH_HOST
ARG SSH_USER
ARG SSH_HOSTNAME
ARG SSH_IDENTITY_FILE

ARG SSH_PRIVATE_KEY
ARG SSH_KEY

RUN cd /root/ && mkdir .ssh && touch .netrc && touch .gitconfig && echo ${NETRC} >> .netrc \
        && echo "[url \${GIT_SSH_URL}\]" > /root/.gitconfig && \
        echo "    insteadOf = ${GIT_URL}" >> /root/.gitconfig && \
        touch /root/.ssh/config && touch /root/.ssh/lms_gitlab && touch /root/.ssh/lms_gitlab.pub && \
        echo "Host ${SSH_HOST}" > /root/.ssh/config && \
        echo "  User ${SSH_USER}" >> /root/.ssh/config && \
        echo "  HostName ${SSH_HOSTNAME}" >> /root/.ssh/config && \
        echo "  IdentityFile ${SSH_IDENTITY_FILE}" >> /root/.ssh/config && \
        echo "  IdentitiesOnly yes" >> /root/.ssh/config && \
        chmod 600 /root/.ssh/config && \
        printf "%s\n" "${SSH_PRIVATE_KEY}" > /root/.ssh/lms_gitlab && \
        chmod 600 /root/.ssh/lms_gitlab && echo ${SSH_KEY} > /root/.ssh/lms_gitlab.pub && \
        chmod 600 /root/.ssh/lms_gitlab.pub

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
