ARG GO_VERSION=1.26

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS app-build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" \
    -o /out/d8r ./cmd/d8r

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS migrate-build

WORKDIR /src

RUN go mod init migrate-build \
    && go get github.com/golang-migrate/migrate/v4/cmd/migrate@v4.20.1

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -tags postgres \
    -o /out/migrate \
    github.com/golang-migrate/migrate/v4/cmd/migrate

FROM alpine:3.22

RUN apk add --no-cache bash ca-certificates \
    && mkdir -p /app/bin /app/config /app/files \
    && chown -R 65532:65532 /app

WORKDIR /app

COPY --from=app-build /out/d8r /app/d8r
COPY --from=migrate-build /out/migrate /usr/local/bin/migrate
COPY --chmod=755 bin/migrate /app/bin/migrate
COPY migrations/ /app/migrations/
COPY --chown=65532:65532 config/config.yaml /app/config/config.yaml

USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/app/d8r"]
CMD ["server", "--configuration", "/app/config/config.yaml"]