GIT_VERSION = $(shell git describe --always)

install:
	docker run --rm \
	--volume ${CURDIR}:/go/src/github.com/VJftw/privacy-negotiator/backend/ \
	--workdir /go/src/github.com/VJftw/privacy-negotiator/backend/priv-neg \
 	--env CGO_ENABLED=0 \
 	--env GOOS=linux \
	golang \
	/bin/sh -c "curl https://glide.sh/get | sh && glide install"

build: install
	docker run --rm \
	--volume ${CURDIR}:/go/src/github.com/VJftw/privacy-negotiator/backend/ \
	--workdir /go/src/github.com/VJftw/privacy-negotiator/backend/priv-neg \
 	--env CGO_ENABLED=0 \
 	--env GOOS=linux \
	golang \
	/bin/sh -c "go build -a -installsuffix cgo -o dist/priv-neg"

	docker build -t vjftw/privacy-negotiator:backend-${GIT_VERSION} .
	docker tag vjftw/privacy-negotiator:backend-${GIT_VERSION} vjftw/privacy-negotiator:backend-latest

build-dev:
	docker build -t vjftw/privacy-negotiator:backend-dev -f dev.Dockerfile .

test: build-dev
	$(eval NOVENDOR=$(shell docker run --rm \
	--volume ${CURDIR}:/go/src/github.com/VJftw/privacy-negotiator/backend/ \
	--workdir /go/src/github.com/VJftw/privacy-negotiator/backend/priv-neg \
	vjftw/privacy-negotiator:backend-dev \
	glide novendor))

	docker run --rm \
	--volume ${CURDIR}:/go/src/github.com/VJftw/privacy-negotiator/backend/ \
	--workdir /go/src/github.com/VJftw/privacy-negotiator/backend/priv-neg \
	vjftw/privacy-negotiator:backend-dev \
	golint -set_exit_status ${NOVENDOR}
