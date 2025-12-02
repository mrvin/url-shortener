lint:
	golangci-lint run ./...
check-format:
	test -z $$(go fmt ./...)

TAG := $(shell git describe --abbrev=0 --tags --always)
HASH := $(shell git rev-parse HEAD)
DATE := $(shell date +%Y-%m-%d.%H:%M:%S)
LDFLAGS := -w -s \
				-X github.com/mrvin/url-shortener/internal/httpserver/handlers.hash=$(HASH) \
				-X github.com/mrvin/url-shortener/internal/httpserver/handlers.tag=$(TAG) \
				-X github.com/mrvin/url-shortener/internal/httpserver/handlers.date=$(DATE)
build:
	go build -ldflags "$(LDFLAGS)" -o bin/url-shortener cmd/url-shortener/main.go
.PHONY: lint check-format build

test:
	mkdir -p reports
	go test ./... -coverprofile=reports/coverage.out
coverage:
	go tool cover -func reports/coverage.out | grep "total:" | \
	awk '{print ((int($$3) > 23) != 1) }'
report:
	go tool cover -html=reports/coverage.out -o reports/cover.html
.PHONY: test coverage report

build-compose:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env --profile prod build
up-compose:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env --profile prod up
run-compose: build-compose up-compose
down-compose:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env --profile prod down
.PHONY: build-compose up-compose run-compose down-compose

