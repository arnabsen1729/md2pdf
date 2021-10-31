BINARY_NAME=md2pdf

build:
	go build -o ${BINARY_NAME} main.go

run:
	./${BINARY_NAME}

dev:
	go run main.go

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}

test:
	go test ./...

lint:
	golangci-lint run --enable-all
