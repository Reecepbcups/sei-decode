BIN=juno-decode

install:
	go install .

build:
	go build .

test:
	${BIN} tx decode-file test.json output.json