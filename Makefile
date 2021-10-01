install:
	go get -v ./...

local: down
	docker compose up --remove-orphans --build

down:
	docker compose down --remove-orphans

test.fmt:
	go fmt ./...;

test.vet:
	go vet ./...;

test.lint:
	golangci-lint run ./...;

test.testfmt:
	test -z $(gofmt -s -l -w .);

test%: export DATABASE_HOST = 0.0.0.0
test%: export DATABASE_PORT = 5432
test%: export DATABASE_NAME = postgres
test%: export DATABASE_USER = postgres
test%: export DATABASE_PASSWORD = password

test.postgres: down
	docker run --env POSTGRES_DB=$(DATABASE_NAME) --env POSTGRES_USER=$(DATABASE_USER) --env POSTGRES_PASSWORD=$(DATABASE_PASSWORD) -p ${DATABASE_PORT}:${DATABASE_PORT} --detach --name db postgres:13.2

test.tests: test.postgres
	go test -cover -race -coverprofile=c.out ./...;

test.coverage: test.tests
	go tool cover -html=c.out -o coverage.html;

test: test.fmt test.vet test.lint test.testfmt test.tests test.coverage

down.db:
	docker rm -f db

down: down.db
	docker compose down --remove-orphans