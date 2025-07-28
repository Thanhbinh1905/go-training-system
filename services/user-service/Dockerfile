FROM golang:1.25rc2-alpine3.22 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/services/user-service

RUN go mod tidy

RUN go build -o user-service ./cmd

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/services/user-service/user-service .

CMD ["./user-service"]