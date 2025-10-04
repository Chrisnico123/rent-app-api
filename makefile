# Makefile for Docker build and push

# Configuration
DOCKER_USERNAME := chrisnico
APP_NAME := rent-application
VERSION := 1.1.6
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# Docker image configuration
IMAGE := ${DOCKER_USERNAME}/${APP_NAME}
TAG := ${VERSION}

.PHONY: help build build-local tag push push-latest version clean

help:  ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

version:  ## Show current version
	@echo "Version: ${VERSION}"
	@echo "Git Commit: ${GIT_COMMIT}"
	@echo "Build Date: ${BUILD_DATE}"

build-staging:  ## Build the Docker image
	docker build \
		--build-arg APP_ENV=staging \
		--build-arg BUILD_DATE="${BUILD_DATE}" \
		--build-arg VERSION="${VERSION}" \
		--build-arg GIT_COMMIT="${GIT_COMMIT}" \
		-t ${IMAGE}:${TAG} \
		-t ${IMAGE}:latest \
		.
run-dev:
	APP_ENV=development go run cmd/main.go

build-local:  ## Build for local development
	docker build \
		-t ${IMAGE}:local \
		.

tag:  ## Tag the Docker image
	docker tag ${IMAGE}:${TAG} ${IMAGE}:${TAG}
	docker tag ${IMAGE}:${TAG} ${IMAGE}:latest

push: tag  ## Push the Docker image to registry
	docker push ${IMAGE}:${TAG}
	docker push ${IMAGE}:latest

push-latest:  ## Push only the latest tag
	docker push ${IMAGE}:latest

clean:  ## Remove local Docker images
	docker rmi ${IMAGE}:${TAG} ${IMAGE}:latest || true

# Version targets
patch:  ## Create and push a patch version (0.0.1 → 0.0.2)
	$(eval NEW_VERSION := $(shell semver bump patch ${VERSION}))
	@git tag -a v${NEW_VERSION} -m "Version ${NEW_VERSION}"
	@git push origin v${NEW_VERSION}

minor:  ## Create and push a minor version (0.1.0 → 0.2.0)
	$(eval NEW_VERSION := $(shell semver bump minor ${VERSION}))
	@git tag -a v${NEW_VERSION} -m "Version ${NEW_VERSION}"
	@git push origin v${NEW_VERSION}

major:  ## Create and push a major version (1.0.0 → 2.0.0)
	$(eval NEW_VERSION := $(shell semver bump major ${VERSION}))
	@git tag -a v${NEW_VERSION} -m "Version ${NEW_VERSION}"
	@git push origin v${NEW_VERSION}