# ── Stage 1: Build ──────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o urlanalyzer .


# ── Stage 2: Run ────────────────────────────────────────────────────────────
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/urlanalyzer .

COPY --from=builder /app/webPages ./webPages

EXPOSE 8080

CMD ["./urlanalyzer"]
