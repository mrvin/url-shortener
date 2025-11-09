build:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env build
up:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env up
run: build up
down:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env down
lint:
	golangci-lint run ./...

.PHONY: build up run down lint

