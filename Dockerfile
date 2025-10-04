FROM golang:1.25 AS builder

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux make build

FROM gcr.io/distroless/base-debian11 AS runner

WORKDIR /

COPY --from=builder /app/bin/sse /app

USER nonroot:nonroot

ENTRYPOINT [ "/app" ]
