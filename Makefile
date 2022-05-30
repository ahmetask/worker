lint:
	golangci-lint run -v --exclude-use-default=false --timeout 2m0s

test:
	go test -cover ./...