BINARY_NAME=md2pdf

build:
	go build -o ${BINARY_NAME} .

run:
	./${BINARY_NAME}

dev:
	go build -o ${BINARY_NAME} .
	./${BINARY_NAME} -file test.md
	xdg-open test.pdf

build_and_run: build run

clean:
	go clean
	rm -f ${BINARY_NAME}
	rm -f *.pdf

test:
	go test ./...

lint:
	golangci-lint run --enable-all --color always --issues-exit-code=0 --disable cyclop,lll
