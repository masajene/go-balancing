.PHONY: build

build-arm:
	env GOOS=linux GOARCH=arm64 go build -o main