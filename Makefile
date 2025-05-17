PROJECT_NAME=short-url
TEMP_FILE=/tmp/shortener_storage.txt
TEST_SERVER_PORT=\$$(random unused-port)

SHELL := /bin/bash
CURRENT_UID := $(shell id -u)
CURRENT_GID := $(shell id -g)

# внешний порт.
LOCAL_PORT=5001
# порт сервера внутри докера.
SERVER_PORT=5001
# адрес, который слушает сервер.
SERVER_ADDRESS=:$(SERVER_PORT)
# урл сервера.
BASE_URL=http://localhost:$(SERVER_PORT)

COMPOSE := export PROJECT_NAME=$(PROJECT_NAME) CURRENT_UID=$(CURRENT_UID) \
 		   CURRENT_GID=$(CURRENT_GID) SERVER_ADDRESS=$(SERVER_ADDRESS) \
 		   BASE_URL=$(BASE_URL) SERVER_PORT=$(SERVER_PORT) \
 		   LOCAL_PORT=$(LOCAL_PORT) && cd docker &&

DOCKER_RUN := cd docker && docker run --rm -v ${PWD}:/app -it $(PROJECT_NAME)-app

build-local:
	go build -o ./cmd/shortener/shortener ./cmd/shortener/*.go

build:
	$(COMPOSE) docker compose -f docker-compose.build.yml build app

go-build:
	cd docker && docker run --rm -v ${PWD}:/app -it $(PROJECT_NAME)-app \
	go build -gcflags "all=-N -l" -o /app/cmd/shortener/shortener -tags dynamic /app/cmd/shortener/

run-test: \
	go-build \
	run-test-a \
	run-test-u \
	run-test-s \
	run-lint

run-test-u:
	$(DOCKER_RUN) sh -c "go test ./..."

run-test-s:
	$(DOCKER_RUN) sh -c "go vet -vettool=\$$(which statictest) ./..."

run-lint:
	$(DOCKER_RUN) golangci-lint run

run-test-a: \
	run-test-a1 \
	run-test-a2 \
	run-test-a3 \
	run-test-a4 \
	run-test-a5 \
	run-test-a6 \
	run-test-a7 \
	run-test-a8 \

run-test-a1:
	$(DOCKER_RUN) sh -c "shortenertestbeta -test.v -test.run=^TestIteration1$$ -binary-path=cmd/shortener/shortener"
run-test-a2:
	$(DOCKER_RUN) sh -c "shortenertestbeta -test.v -test.run=^TestIteration2$$ -source-path=."
run-test-a3:
	$(DOCKER_RUN) sh -c "shortenertestbeta -test.v -test.run=^TestIteration3$$ -source-path=."
run-test-a4:
	$(DOCKER_RUN) sh -c "shortenertestbeta -test.v -test.run=^TestIteration4$$ -binary-path=cmd/shortener/shortener -server-port=$(TEST_SERVER_PORT)"
run-test-a5:
	$(DOCKER_RUN) sh -c "shortenertestbeta -test.v -test.run=^TestIteration5$$ -binary-path=cmd/shortener/shortener -server-port=$(TEST_SERVER_PORT)"
run-test-a6:
	$(DOCKER_RUN) sh -c "shortenertestbeta -test.v -test.run=^TestIteration6$$ -source-path=."
run-test-a7:
	$(DOCKER_RUN) sh -c "shortenertestbeta -test.v -test.run=^TestIteration7$$ -binary-path=cmd/shortener/shortener -source-path=."
run-test-a8:
	$(DOCKER_RUN) sh -c "shortenertestbeta -test.v -test.run=^TestIteration8$$ -binary-path=cmd/shortener/shortener"
run-test-a9:
	$(DOCKER_RUN) sh -c "shortenertestbeta -test.v -test.run=^TestIteration9$$ -binary-path=cmd/shortener/shortener -source-path=. -file-storage-path=$(TEMP_FILE)"

update-tpl:
	# git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
	git fetch template && git checkout template/main .github

up: \
	go-build
	$(COMPOSE) docker compose up -d container

down:
	$(COMPOSE) docker compose down container