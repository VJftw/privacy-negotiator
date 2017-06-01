GIT_VERSION = $(shell git describe --always)

build:
	cd web_app && make build

test:
	echo "HELLO"
