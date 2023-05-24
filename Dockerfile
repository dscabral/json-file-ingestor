FROM golang:1.20 as base
ENV CGO_ENABLED=0
WORKDIR /opt/app
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM base as build
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /out/ports src/cmd/main.go

FROM golangci/golangci-lint:v1.50.1-alpine as lint-base

FROM base as lint
COPY --from=lint-base /usr/bin/golangci-lint /usr/bin/golangci-lint
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    golangci-lint run --timeout 10m0s ./...

FROM base AS unit-test
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    go test -v ./...

FROM scratch as webservice
EXPOSE 3000
COPY --from=build /out/ports .
COPY ports.json .
CMD ["./ports"]
