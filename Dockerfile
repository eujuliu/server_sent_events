FROM golang:1.25 AS builder

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux make build

FROM debian:bullseye-slim AS runner

RUN apt-get update && apt-get install -y curl

WORKDIR /

RUN addgroup --system --gid 1001 app
RUN adduser --system --uid 1001 app

USER app

COPY --from=builder --chown=app:app /app/bin/sse /app

ENTRYPOINT [ "/app" ]
