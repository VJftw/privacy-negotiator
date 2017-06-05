GIT_VERSION = $(shell git describe --always)

install:
	docker run --rm \
	--volume ${CURDIR}:/go/src/github.com/VJftw/privacy-negotiator/api/ \
	--workdir /go/src/github.com/VJftw/privacy-negotiator/api/priv-neg \
 	--env CGO_ENABLED=0 \
 	--env GOOS=linux \
	golang \
	/bin/sh -c "curl https://glide.sh/get | sh && glide install"

	docker build -t vjftw/privacy-negotiator:api-${GIT_VERSION} .

build: install
	docker run --rm \
	--volume ${CURDIR}:/go/src/github.com/VJftw/privacy-negotiator/api/ \
	--workdir /go/src/github.com/VJftw/privacy-negotiator/api/priv-neg \
 	--env CGO_ENABLED=0 \
 	--env GOOS=linux \
	golang \
	/bin/sh -c "go build -a -installsuffix cgo -o dist/priv-neg"

	docker build -t vjftw/privacy-negotiator:api-${GIT_VERSION} .
