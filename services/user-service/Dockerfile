FROM golang:1.25rc2-alpine3.22

# 1. Set WORKDIR tại gốc project (để go.work hoạt động)
WORKDIR /app

# 2. Copy toàn bộ monorepo vào container
COPY . .

# 3. Làm việc tại team-service
WORKDIR /app/services/user-service

# 4. Cài air
RUN go install github.com/air-verse/air@latest

CMD ["air"]