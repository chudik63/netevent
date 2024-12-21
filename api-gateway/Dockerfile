FROM alpine:latest AS setup

RUN apk --no-cache add \
    openssh \
    git \
    bash \
    ca-certificates \
    && mkdir -p /root/.ssh

ARG NETRC
ARG GIT_URL
ARG GIT_SSH_URL
ARG SSH_HOST
ARG SSH_USER
ARG SSH_HOSTNAME
ARG SSH_IDENTITY_FILE
ARG SSH_PRIVATE_KEY
ARG SSH_KEY
ARG KNOWN_HOSTS

RUN echo "${NETRC}" > /root/.netrc && \
    echo "${KNOWN_HOSTS}" > /root/.ssh/known_hosts && \
    touch /root/.ssh/config && \
    echo "Host ${SSH_HOST}" > /root/.ssh/config && \
    echo "  User ${SSH_USER}" >> /root/.ssh/config && \
    echo "  HostName ${SSH_HOSTNAME}" >> /root/.ssh/config && \
    echo "  IdentityFile /root/.ssh/lms_gitlab" >> /root/.ssh/config && \
    echo "  IdentitiesOnly yes" >> /root/.ssh/config && \
    chmod 600 /root/.ssh/config && \
    chmod 600 /root/.ssh/known_hosts

# Записываем SSH-ключи
RUN printf "%s\n" "${SSH_PRIVATE_KEY}" > /root/.ssh/lms_gitlab && \
    echo "${SSH_KEY}" > /root/.ssh/lms_gitlab.pub && \
    chmod 600 /root/.ssh/lms_gitlab /root/.ssh/lms_gitlab.pub && \
    chmod 700 /root/.ssh

FROM golang:1.23 AS build

COPY --from=setup /root/.ssh /root/.ssh
COPY --from=setup /root/.netrc /root/.netrc

ENV GOPRIVATE="gitlab.crja72.ru/*"

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
