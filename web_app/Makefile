GIT_VERSION = $(shell git describe --always)

install:
	docker run --rm -v ${CURDIR}/priv-neg:/app -w /app node:8 npm install

build: install
	docker run --rm -v ${CURDIR}/priv-neg:/app -w /app node:8 node_modules/.bin/ng build --environment prod
	docker build -t vjftw/priv-neg:${GIT_VERSION} .
