GIT_VERSION = $(shell git describe --always)

install:
	docker run --rm -v ${CURDIR}/priv-neg:/app -w /app node:8 npm install

build: install
	# TODO: Better prod environment setup
	cp priv-neg/src/environments/environment.prod.tpl.ts priv-neg/src/environments/environment.prod.ts
	sed -i "s#gold#${GIT_VERSION}#" priv-neg/src/environments/environment.prod.ts
	# Build Production environment
	docker run --rm -v ${CURDIR}/priv-neg:/app -e GIT_VERSION=${GIT_VERSION} -w /app node:8 node_modules/.bin/ng build --environment prod
	# Build Dockerfile
	docker build -t vjftw/priv-neg:${GIT_VERSION} .

test: install
	docker run --rm -v ${CURDIR}/priv-neg:/app -e GIT_VERSION=${GIT_VERSION} -w /app node:8 node_modules/.bin/tslint src/**/*.ts{,x}
