lint:
	golangci-lint run ./...
check-format:
	test -z $$(go fmt ./...)
build:
	go build -ldflags '-w -s' -o bin/url-shortener cmd/url-shortener/main.go
up:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env up
down:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env down


.PHONY: lint check-format build up run down 

