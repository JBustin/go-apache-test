# Docker image version
DK_REGISTRY ?= jbustin1/

DK_VERSION ?= 0.0.1
DK_NAME ?= go-apache-test
DK_IMAGE ?= $(DK_REGISTRY)$(DK_NAME):$(DK_VERSION)
DK_LATEST ?= $(DK_REGISTRY)$(DK_NAME):latest

DK_DEV_NAME ?= go-apache-test-dev
DK_DEV_IMAGE ?= $(DK_REGISTRY)$(DK_DEV_NAME):latest

.DEFAULT_GOAL := help

help: ## Display this help
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN { FS = ":.*?## " }; { printf "\033[36m%-30s\033[0m %s\n", $$1, $$2 }'

####### DEV IMG

build-dev: ## Build dev image
	@docker build -t ${DK_DEV_IMAGE} -f .docker/dev/Dockerfile .

run-dev: ## Run go-apache-test for dev
	@docker run --rm \
	-v ${PWD}:/usr/go-apache-test/app \
	-v ${PWD}/tmp/config.json:/usr/go-apache-test/config.json \
	-v ${PWD}/tmp/tests:/usr/go-apache-test/tests \
	-v ${PWD}/tmp/vhosts:/usr/go-apache-test/vhosts \
	${DK_DEV_IMAGE}


run-dev-interactive: ## Run go-apache-test for dev
	@docker run --rm -it \
	-v ${PWD}:/usr/go-apache-test/app \
	-v ${PWD}/tmp/config.json:/usr/go-apache-test/config.json \
	-v ${PWD}/tmp/tests:/usr/go-apache-test/tests \
	-v ${PWD}/tmp/vhosts:/usr/go-apache-test/vhosts \
	${DK_DEV_IMAGE} bash

####### PRD IMG

build: ## Build go-apache-test
	@docker build -t ${DK_IMAGE} -f .docker/prd/Dockerfile .

build-latest: ## Build latest go-apache-test
	@docker build -t ${DK_LATEST} -f .docker/prd/Dockerfile .

push: ## Push go-apache-test
	@docker push ${DK_IMAGE}

push-latest: ## Push latest go-apache-test
	@docker push $(DK_LATEST)

build-and-push: ## Build and push
	@make build
	@make push

build-and-push-latest: ## Build and push latest
	@make build-latest
	@make push-latest

run: ## Run go-apache-test
	@docker run --rm \
	-v ${PWD}/tmp/tests:/usr/go-apache-test/tests \
	-v ${PWD}/tmp/vhosts:/usr/go-apache-test/vhosts \
	-v ${PWD}/tmp/config.json:/usr/go-apache-test/config.json \
	${DK_IMAGE}

build-and-run: ## Build and run go-apache-test
	@make build
	@make run

# interactive: ## Run go-apache-test in interactive mode
# 	@docker run --rm -it \
# 	-v ${PWD}/tests:/usr/app/tests \
# 	-v ${PWD}/vhosts:/usr/app/vhosts \
# 	-v ${PWD}/.env:/usr/app/.env \
# 	${DK_IMAGE} bash

# dev-interactive: ## Run dev go-apache-test in interactive mode
# 	@docker run --rm -it \
# 	-v ${PWD}/tests:/usr/app/tests \
# 	-v ${PWD}/lib:/usr/app/lib \
# 	-v ${PWD}/vhosts:/usr/app/vhosts \
# 	-v ${PWD}/.env:/usr/app/.env \
# 	${DK_IMAGE} bash