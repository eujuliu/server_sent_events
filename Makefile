.PHONY: build run watch watch_down test test_with_cache coverage clean check pre-commit

BINARY_NAME=sse
OUTPUT_DIR=bin
MAIN_FILE=./cmd/app/main.go
COVERAGE_DIR=coverage

build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_FILE)

docker_build:
	docker build -t eujuliu/server_sent_events .

dev:
	air

watch:
	cat .env* > .env.docker; \
	docker compose --env-file ./.env.docker -f ./docker-compose.development.yml -p sse --profile all watch --prune; \
	make watch_down

watch_down:
	docker stop $$(docker ps -a -q --filter "label=com.docker.compose.project=sse"); \
	docker system prune -a --filter "label=com.docker.compose.project=sse" -f; \
	rm -f .env.docker

run:
	make build
	./$(OUTPUT_DIR)/scheduler

test:
	go test -tags=unit -count=1 -short -v ./...

test_w_cache:
	go test -tags=unit -short -v ./...

coverage:
	make clean; \
	mkdir $(COVERAGE_DIR); \
	go test -coverprofile=coverage/cover.out ./...; \
	go tool cover -html=coverage/cover.out -o coverage/cover.html

check:
	golangci-lint fmt ./...
	golangci-lint run ./...

clean:
	rm -rf $(OUTPUT_DIR); \
	rm -rf $(COVERAGE_DIR)

pre-commit:
	pre-commit install
