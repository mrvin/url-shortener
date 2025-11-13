lint:
	golangci-lint run ./...
check-format:
	test -z $$(go fmt ./...)
build:
	go build -ldflags '-w -s' -o bin/url-shortener cmd/url-shortener/main.go
build-compose:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env build
up-compose:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env up
run-compose: build-compose up-compose
down-compose:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env down


.PHONY: lint check-format build build-compose up-compose run-compose down-compose

