PROTO_DIR = internal/proto

UNAME := $(shell uname -s)
OS = macos
PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')

proto:
	protoc -I${PROTO_DIR} --go_opt=module=${PACKAGE} --go_out=. ${PROTO_DIR}/*.proto

clean:
	rm ${PROTO_DIR}/*.pb.go

install:
	go get -v ./...

test.fmt:
	go fmt ./...;

test.vet:
	go vet ./...;

test.lint:
	golangci-lint run ./...;

test.testfmt:
	test -z $(gofmt -s -l -w .);

test.tests:
	go test -cover -race -coverprofile=c.out ./...;

test.coverage: test.tests
	go tool cover -html=c.out -o coverage.html;

test: test.fmt test.vet test.lint test.testfmt test.tests test.coverage

local: down
	docker compose up --remove-orphans --build

down.db:
	docker rm -f db

down: down.db
	docker compose down --remove-orphans
