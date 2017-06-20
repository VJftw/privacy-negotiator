GIT_VERSION = $(shell git describe --always)
ENVIRONMENT ?= beta

install:
	docker run --rm -v ${CURDIR}/priv-neg:/app -w /app node:8 npm install --silent

build: install
	# TODO: Better environment setup
	cp priv-neg/src/environments/environment.${ENVIRONMENT}.tpl.ts priv-neg/src/environments/environment.${ENVIRONMENT}.ts
	sed -i "s#gold#${GIT_VERSION}#" priv-neg/src/environments/environment.${ENVIRONMENT}.ts
	# Build environment
	docker run --rm -v ${CURDIR}/priv-neg:/app -e GIT_VERSION=${GIT_VERSION} -w /app node:8 node_modules/.bin/ng build --environment ${ENVIRONMENT}

test: install
	docker run --rm -v ${CURDIR}/priv-neg:/app -e GIT_VERSION=${GIT_VERSION} -w /app node:8 node_modules/.bin/tslint src/**/*.ts{,x}
