GIT_VERSION = $(shell git describe --always)
AWS_DEFAULT_REGION ?= eu-west-1
ENVIRONMENT ?= beta
DOMAIN ?= beta.privacy-negotiator.vjpatel.me

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

cluster-destroy: cluster-init
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/cluster \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	hashicorp/terraform:0.9.7 \
	destroy --force

deploy-init: tf-fmt
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/env/${ENVIRONMENT} \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	hashicorp/terraform:0.9.7 \
	init

deploy-plan: deploy-init
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/env/${ENVIRONMENT} \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	hashicorp/terraform:0.9.7 \
	plan

deploy-apply: build deploy-plan
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/env/${ENVIRONMENT} \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	hashicorp/terraform:0.9.7 \
	apply

	# Upload static Web to Cloudfront
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/web_app/priv-neg \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	anigeo/awscli:latest \
	s3 cp dist/. s3://${DOMAIN} --acl public-read --recursive --cache-control max-age=120
