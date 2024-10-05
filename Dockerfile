# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.1
ARG NODE_VERSION=20.17.0

FROM golang:${GO_VERSION} AS core

ARG ENTRY_DIRECTORY

ENV ENTRY_DIRECTORY=${ENTRY_DIRECTORY}

WORKDIR /app

COPY . .

WORKDIR /app/${ENTRY_DIRECTORY}

RUN --mount=type=cache,target=/go/pkg/mod  \
    CGO_ENABLED=0 GOOS=linux go build -o /entry

FROM node:${NODE_VERSION}-alpine AS builder

COPY --from=core /app /app

WORKDIR /app/web

ENV YARN_CACHE_FOLDER=/root/.yarn \
    PATH="$PATH:/app/web/node_modules/.bin"

VOLUME /root/.yarn

RUN --mount=type=cache,target=/root/.yarn yarn install --frozen-lockfile && yarn build

ENTRYPOINT ["/app/web/entrypoint.sh"]

CMD [ "yarn", "build" ]

FROM scratch AS runnable

LABEL org.opencontainers.image.source=https://github.com/senorihl/yopa
LABEL org.opencontainers.image.licenses=UNLICENSED

COPY --from=core /entry /entry
COPY --from=builder /app/web/dist /app/web/dist

ENTRYPOINT ["/entry"]

FROM core AS watchable

COPY --from=core /app /app

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod  \
    go install github.com/mitranim/gow@latest

CMD gow -v run "./${ENTRY_DIRECTORY}"

FROM node:${NODE_VERSION}-alpine as demo

COPY --from=builder /app /app

WORKDIR /app/demo

EXPOSE 80

ENTRYPOINT ["/app/demo/entrypoint.sh"]

CMD [ "yarn", "serve" ]