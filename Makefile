.PHONY: build test run clean

build:
	go build -o bin/app .

test:
	go test ./...

run:
	go run .

clean:
	rm -rf bin/
	go clean