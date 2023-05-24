FROM golang:1.18 as base
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
    --mount=type=cache,target=/root/.cache/golangci-lint \
      golangci-lint run --timeout 10m0s ./...

FROM base AS unit-test
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    go test -v -count=1 -p=1 ./...

FROM scratch as service
COPY --from=build /out/ports .
CMD ["./ports"]
