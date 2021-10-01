ARG go_version=1.16

#User
FROM alpine:3.13.1 as user
ARG uid=10001
ARG gid=10001
RUN echo "scratchuser:x:${uid}:${gid}::/home/scratchuser:/bin/sh" > /scratchpasswd

#Certs
FROM alpine:3.13.2 as certs
RUN apk add -U --no-cache ca-certificates

# BUILDER
FROM golang:${go_version}-alpine as build
WORKDIR /code/
ENV GO111MODULE=on

COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/
COPY ./go.mod/ ./go.sum/ ./
RUN go mod download

RUN set -eux && GOOS=linux CGO_ENABLED=0 GOGC=off GOARCH=amd64 go build -o ./bin/sync ./cmd/sync

# API
FROM scratch as sync
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=user /scratchpasswd /etc/passwd
COPY --from=build /code/bin/sync .
USER scratchuser
EXPOSE 8001
ENTRYPOINT ["./sync"]