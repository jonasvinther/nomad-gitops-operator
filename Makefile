.PHONY: test build compile

test:
	go test -v ./...

build:
	go build -o nomoperator

compile:
	echo "Compiling for every OS and Platform"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/nomoperator-linux-amd64 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o bin/nomoperator-linux-arm main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/nomoperator-linux-arm64 main.go
	CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -o bin/nomoperator-freebsd-386 main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o bin/nomoperator.exe main.go

all: test build