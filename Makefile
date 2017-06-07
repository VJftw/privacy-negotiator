GIT_VERSION = $(shell git describe --always)
AWS_DEFAULT_REGION ?= eu-west-1

build:
	cd web_app && make build

test:
	echo "HELLO"

tf-fmt:
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/cluster \
	hashicorp/terraform:0.9.7 \
	fmt

cluster-init: tf-fmt
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/cluster \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	hashicorp/terraform:0.9.7 \
	init

cluster-plan: cluster-init
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/cluster \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	hashicorp/terraform:0.9.7 \
	plan

cluster-apply: cluster-plan
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/cluster \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	hashicorp/terraform:0.9.7 \
	apply
