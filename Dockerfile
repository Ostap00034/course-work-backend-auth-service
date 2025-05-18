# ─── Stage 1: build ───────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder
WORKDIR /app

# cache modules
COPY go.mod go.sum ./
RUN go mod download

# build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o svc ./cmd/api/main.go

# ─── Stage 2: final ───────────────────────────────────────────────────────────────
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/svc .
EXPOSE 50051
ENTRYPOINT ["./svc"]