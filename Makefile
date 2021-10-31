BINARY_NAME=md2pdf

build:
	go build -o ${BINARY_NAME} .

run:
	./${BINARY_NAME}

dev:
	go run .

build_and_run: build run

clean:
	go clean
	rm -f ${BINARY_NAME}
	rm -f *.pdf

test:
	go test ./...

lint:
	golangci-lint run --enable-all
