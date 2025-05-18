# syntax=docker/dockerfile:1.3

########################################
# 👉 builder: собираем Go-бинарник
########################################
FROM golang:1.23-alpine AS builder
WORKDIR /app

# 1) Копируем только манифесты — кэшируем зависимости
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# 2) Копируем весь код и собираем бинарник, снова пользуясь кешем
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o svc ./cmd/api/main.go

########################################
# 👉 runtime: минимальный Alpine с нашим svc
########################################
FROM alpine:3.18
WORKDIR /app

# Копируем собранный сервис
COPY --from=builder /app/svc .

# Порт, поправьте под свой
EXPOSE 50051

# Запуск
ENTRYPOINT ["./svc"]