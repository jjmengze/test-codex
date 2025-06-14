# --- Stage 1: Builder ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

ENV GOPROXY=direct
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/log-receiver ./cmd/server

## --- Stage 2: Final Image ---
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/log-receiver /usr/local/bin/log-receiver
COPY --from=builder /app/config /app/config

EXPOSE 8080

ENTRYPOINT ["log-receiver"]
CMD ["--port=80"]
