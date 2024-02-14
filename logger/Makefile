server:
	go run main.go --runserver

test:
	go test -v -cover ./...

mock:
	mockgen -package mock -destination internal/repo/mock/logger.go logger/internal/repo  LoggerRepoImply