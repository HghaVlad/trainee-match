FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go

FROM alpine:3.23

WORKDIR /src

COPY --from=builder /app/auth-service .
COPY --from=builder /app/internal/infra/db/postgres/migrations /src/internal/infra/db/postgres/migrations

EXPOSE 8000

CMD ["./auth-service"]
