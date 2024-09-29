# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.1

FROM golang:${GO_VERSION} as core

ARG ENTRY_DIRECTORY

ENV ENTRY_DIRECTORY=${ENTRY_DIRECTORY}

WORKDIR /app

COPY . .

WORKDIR /app/${ENTRY_DIRECTORY}

RUN --mount=type=cache,target=/go/pkg/mod  \
    CGO_ENABLED=0 GOOS=linux go build -o /entry

FROM scratch AS runnable

COPY --from=core /entry /entry

ENTRYPOINT ["/entry"]

FROM core AS watchable

COPY --from=core /app /app

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod  \
    go install github.com/mitranim/gow@latest

CMD gow -v run "./${ENTRY_DIRECTORY}"