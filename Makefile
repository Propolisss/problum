DOCKER_COMPOSE_FILE := docker-compose.local.yaml
DOCKER_USER := propolisss

BACKEND_IMAGE := ${DOCKER_USER}/problum-backend
FRONTEND_IMAGE := ${DOCKER_USER}/problum-frontend
WORKER_IMAGE := ${DOCKER_USER}/problum-worker
PLATFORMS := linux/amd64,linux/arm64

.PHONY: up
up:
	docker compose -f ${DOCKER_COMPOSE_FILE} up --build --scale worker=3 -d

.PHONY: down
down:
	docker compose -f ${DOCKER_COMPOSE_FILE} down

.PHONY: push
push:
	docker buildx build --platform ${PLATFORMS} \
	-t ${BACKEND_IMAGE}:latest \
	-f apps/backend/Dockerfile.app apps/backend/ \
	--push

	docker buildx build --platform ${PLATFORMS} \
	-t ${WORKER_IMAGE}:latest \
	-f apps/backend/Dockerfile.worker apps/backend/ \
	--push

	docker buildx build --platform ${PLATFORMS} \
	-t ${FRONTEND_IMAGE}:latest \
	-f apps/frontend/Dockerfile apps/frontend/ \
	--push
