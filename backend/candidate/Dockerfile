FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY candidate/go.mod /src/candidate/go.mod
COPY candidate/go.sum /src/candidate/go.sum

COPY contracts /src/contracts

WORKDIR /src/candidate

RUN go mod download

COPY candidate /src/candidate
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/candidate-service ./cmd/main.go

FROM alpine:3.23

WORKDIR /src
COPY --from=builder /bin/candidate-service .
COPY --from=builder /src/candidate/internal/infrastructure/db/postgres/migrations /src/internal/infrastructure/db/postgres/migrations

EXPOSE 8081
EXPOSE 8085

CMD ["./candidate-service"]