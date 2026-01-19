FROM golang:1.25.1-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/

FROM gcr.io/distroless/base-debian11 AS release
WORKDIR /
COPY --from=build app/main .
USER nonroot:nonroot

CMD ["./main"]