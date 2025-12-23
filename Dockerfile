## Build
FROM golang:1.25.4-alpine AS dev

LABEL maintainer="mrvin v.v.vinogradovv@gmail.com"

RUN apk update && apk add make && apk add tzdata && apk add git

WORKDIR /app

# Copy the code into the container.
COPY cmd/url-shortener cmd/url-shortener
COPY internal internal
COPY pkg pkg
COPY Makefile ./
COPY .git .git
COPY api api

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN make build

RUN mkdir /var/log/url-shortener/

ENV TZ=Europe/Moscow

EXPOSE 8080

ENTRYPOINT ["/app/bin/url-shortener"]

## Deploy
FROM scratch AS prod

LABEL maintainer="mrvin v.v.vinogradovv@gmail.com"

WORKDIR /

COPY --from=dev ["/app/api/openapi.yaml", "/app/api/openapi.yaml"]
COPY --from=dev ["/var/log/url-shortener/", "/var/log/url-shortener/"]
COPY --from=dev ["/usr/share/zoneinfo", "/usr/share/zoneinfo"]
COPY --from=dev ["/app/bin/url-shortener", "/usr/local/bin/url-shortener"]

ENV TZ=Europe/Moscow

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/url-shortener"]
