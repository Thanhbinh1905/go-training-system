FROM golang:1.25rc2-alpine3.22 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/migration

RUN go mod tidy

RUN go build -o migration .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/migration .

CMD ["./migration"]