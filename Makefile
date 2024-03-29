GIT_VERSION = $(shell git describe --always)
AWS_DEFAULT_REGION ?= eu-west-1
ENVIRONMENT ?= beta
DOMAIN ?= beta.privacymanager.social

build:
	cd web_app && make build
	cd backend && make build

tf-fmt:
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/cluster \
	hashicorp/terraform:0.9.9 \
	fmt

cluster-init: tf-fmt
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/cluster \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	hashicorp/terraform:0.9.9 \
	init

cluster-plan: cluster-init
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/cluster \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	hashicorp/terraform:0.9.9 \
	plan

cluster-apply: cluster-plan
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/cluster \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	hashicorp/terraform:0.9.9 \
	apply

cluster-destroy: cluster-init
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/cluster \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	hashicorp/terraform:0.9.9 \
	destroy --force

deploy-init: tf-fmt
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/env/${ENVIRONMENT} \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	anigeo/awscli:latest \
	s3 cp s3://privneg-terraform/vars/${ENVIRONMENT}.terraform.tfvars terraform.tfvars

	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/env/${ENVIRONMENT} \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	hashicorp/terraform:0.9.9 \
	init

deploy-plan: deploy-init
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/env/${ENVIRONMENT} \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	--env TF_VAR_version=${GIT_VERSION} \
	hashicorp/terraform:0.9.9 \
	plan

deploy-apply: build deploy-plan
	# Push API images
	docker push vjftw/privacy-negotiator:backend-${GIT_VERSION}
	docker push vjftw/privacy-negotiator:backend-latest

	# Terraform Apply
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/env/${ENVIRONMENT} \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	--env TF_VAR_version=${GIT_VERSION} \
	hashicorp/terraform:0.9.9 \
	apply

	# Upload static Web to S3
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/web_app/priv-neg \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	anigeo/awscli:latest \
	s3 cp dist/. s3://${DOMAIN} --acl public-read --recursive --cache-control max-age=120

deploy-destroy: deploy-init
	# Remove static Web from S3
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/web_app/priv-neg \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	anigeo/awscli:latest \
	s3 rm s3://${DOMAIN} --recursive

	# Terraform Destroy
	docker run --rm \
	--volume ${CURDIR}:/app \
	--workdir /app/infrastructure/env/${ENVIRONMENT} \
	--env AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
	--env AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
	--env AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
	--env TF_VAR_version=${GIT_VERSION} \
	hashicorp/terraform:0.9.9 \
	destroy --force
