FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o candidate-service ./cmd/main.go

FROM alpine:3.23

WORKDIR /root/

COPY --from=builder /app .

EXPOSE 8081

CMD ["./candidate-service"]