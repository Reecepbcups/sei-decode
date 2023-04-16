BIN=juno-decode

install:
	go install .

build:
	go build .

test: install
	${BIN} tx decode-file test.json output.json