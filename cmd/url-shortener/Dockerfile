## Build
FROM golang:1.25.3-alpine AS build

LABEL maintainer="mrvin v.v.vinogradovv@gmail.com"

RUN apk update && apk add make && apk add tzdata

WORKDIR /app

# Copy the code into the container.
COPY cmd/url-shortener cmd/url-shortener
COPY internal internal
COPY pkg pkg

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/url-shortener/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/usr/share/zoneinfo", "/usr/share/zoneinfo"]
COPY --from=build ["/app/bin/url-shortener", "/usr/local/bin/url-shortener"]

ENV TZ=Europe/Moscow

EXPOSE 8081

ENTRYPOINT ["/usr/local/bin/url-shortener"]
