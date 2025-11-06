DOCKER_COMPOSE_FILE := docker-compose.yaml

.PHONY: up
up:
	docker compose -f ${DOCKER_COMPOSE_FILE} up --build --scale worker=3 -d

.PHONY: down
down:
	docker compose -f ${DOCKER_COMPOSE_FILE} down
