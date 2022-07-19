.PHONY: help run
help:
	go run ./cmd/uploader/main.go -h

run:
	go run ./cmd/uploader/main.go -d $(directory)
