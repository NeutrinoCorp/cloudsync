.PHONY: help run
help:
	go run ./cmd/uploader/main.go -h

run:
	go run ./cmd/uploader/main.go -d $(directory)

test:
	go test ./... -v --cover

test-cov:
	go test ./... -coverprofile coverage.out . && go tool cover -html=coverage.out
