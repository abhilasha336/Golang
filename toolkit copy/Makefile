:PHONY

test:
	go clean -testcache
	go test ./... 
	golangci-lint run