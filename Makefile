BINARY_NAME=md2pdf

build:
	go build -o ${BINARY_NAME} .

run:
	./${BINARY_NAME}

dev:
	go build -o ${BINARY_NAME} .
	./${BINARY_NAME} -file sample.md
	xdg-open sample.pdf

build_and_run: build run

clean:
	go clean
	rm -f ${BINARY_NAME}

test:
	go test ./... -v | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

lint:
	golangci-lint run --enable-all --color always --issues-exit-code=0 --disable cyclop,lll,gomnd

cover:
	go test ./... -cover
