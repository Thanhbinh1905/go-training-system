FROM golang:1.25rc2-alpine3.22

WORKDIR /app

COPY . .

WORKDIR /app/services/asset-service

RUN go install github.com/air-verse/air@latest

CMD ["air"]