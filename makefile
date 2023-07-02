.PHONY: run
run:
	go run main.go

.PHONY: dev
dev:
	CompileDaemon -build="go build -o main" -command="./main"

.PHONY: test
test:
	go test ./...

.PHONY: test-verbose
test-verbose:
	go test -v ./...

.PHONY: test-bench
test-bench:
	go test -bench=.

.PHONY: test-no-cache
test-no-cache:
	go test -count=1 ./...

.PHONY: test-cover
test-cover:
	go test -cover ./...

.PHONY: test-cover-output
test-cover-output:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: build
build:
	go build -o myapp main.go
