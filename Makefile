lint:
	golangci-lint run ./...
check-format:
	test -z $$(go fmt ./...)
build:
	go build -ldflags '-w -s' -o bin/url-shortener cmd/url-shortener/main.go
.PHONY: lint check-format build

test:
	mkdir -p reports
	go test ./... -coverprofile=reports/coverage.out
coverage:
	go tool cover -func reports/coverage.out | grep "total:" | \
	awk '{print ((int($$3) > 2) != 1) }'
report:
	go tool cover -html=reports/coverage.out -o reports/cover.html
.PHONY: test coverage report

build-compose:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env build
up-compose:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env up
run-compose: build-compose up-compose
down-compose:
	docker compose -f deployments/docker-compose.yaml --env-file configs/url-shortener.env down
.PHONY: build-compose up-compose run-compose down-compose

