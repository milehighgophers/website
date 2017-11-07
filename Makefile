run:
	go run main.go

test:
	go test ./...

fmt:
	find . -type f -name '*.go' | xargs gofmt -w

