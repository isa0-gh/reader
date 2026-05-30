# Build Stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go

# Production Stage
FROM alpine:latest

WORKDIR /app

# Copy binary and configuration
COPY --from=builder /app/server .
COPY .env.example .env

EXPOSE 8080

CMD ["./server"]
