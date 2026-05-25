DOCKER_COMPOSE_FILE := docker-compose.local.yaml
DOCKER_USER := propolisss

BACKEND_IMAGE := ${DOCKER_USER}/problum-backend
FRONTEND_IMAGE := ${DOCKER_USER}/problum-frontend
WORKER_IMAGE := ${DOCKER_USER}/problum-worker
PLATFORMS := linux/amd64,linux/arm64

TESTER_IMAGE := problum-integration-tester
NETWORK_NAME := problum_network
BACKEND_DIR := apps/backend
GOLANGCI_CONFIG := ../../.golangci.yml
COVERAGE_OUT := coverage.out
COVERAGE_HTML := coverage.html

.PHONY: up
up:
	docker compose -f ${DOCKER_COMPOSE_FILE} up --build --scale worker=3 -d

.PHONY: down
down:
	docker compose -f ${DOCKER_COMPOSE_FILE} down

.PHONY: generate
generate:
	cd ${BACKEND_DIR} && go generate ./internal/...

.PHONY: fmt
fmt:
	cd ${BACKEND_DIR} && go fmt ./...

.PHONY: tests
tests:
	cd ${BACKEND_DIR} && go test -count=1 -race ./internal/...

.PHONY: tests-coverage
tests-coverage:
	cd ${BACKEND_DIR} && \
	COVERPKG=$$(go list ./internal/... | grep -v '/mocks' | paste -sd, -) && \
	go test -count=1 -race -covermode=atomic -coverpkg="$$COVERPKG" -coverprofile=${COVERAGE_OUT} ./internal/... && \
	grep -v -E "mock|_mock.go|/mocks/" ${COVERAGE_OUT} > ${COVERAGE_OUT}.tmp && \
	mv ${COVERAGE_OUT}.tmp ${COVERAGE_OUT} && \
	go tool cover -func=${COVERAGE_OUT} | tail -n 1

.PHONY: tests-coverage-open
tests-coverage-open: tests-coverage
	cd ${BACKEND_DIR} && go tool cover -html=${COVERAGE_OUT}

.PHONY: clean
clean:
	find . -type f \( -name '*.out' -o -name '${COVERAGE_HTML}' \) -delete

.PHONY: lint
lint:
	cd ${BACKEND_DIR} && golangci-lint run --config ${GOLANGCI_CONFIG} ./...

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

.PHONY: test-integration
test-integration:
	@printf -- "             Запускаем все сервисы   "
	@(trap 'tput cnorm' EXIT; tput civis; \
	  while true; do \
	    for x in . .. ...; do printf -- "\r             Запускаем все сервисы %-3s" "$$x"; sleep 0.5; done; \
	  done) & pid=$$!; \
	  make up > /dev/null 2>&1; \
	  kill $$pid; wait $$pid 2>/dev/null; \
	  tput cnorm; echo "\r             Запускаем все сервисы [OK]   "

	@printf -- "             Собираем образ тестера   "
	@(trap 'tput cnorm' EXIT; tput civis; \
	  while true; do \
	    for x in . .. ...; do printf -- "\r             Собираем образ тестера %-3s" "$$x"; sleep 0.5; done; \
	  done) & pid=$$!; \
	  docker build -t ${TESTER_IMAGE} tests/integration > /dev/null 2>&1; \
	  kill $$pid; wait $$pid 2>/dev/null; \
	  tput cnorm; echo "\r             Собираем образ тестера [OK]   "

	@echo "             Начинаем тестировать"
	@docker run --rm \
		--network $(NETWORK_NAME) \
		-e API_URL=http://backend:8080 \
		-e DB_HOST=postgres \
		-e DB_USER=problum \
		-e DB_PASS=problum \
		-e DB_NAME=problum \
		${TESTER_IMAGE} && \
		(echo "--- TESTS PASSED ---"; make down > /dev/null 2>&1) || \
		(echo "--- TESTS FAILED ---"; make down > /dev/null 2>&1; exit 1)
