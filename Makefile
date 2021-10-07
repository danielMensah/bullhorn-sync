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