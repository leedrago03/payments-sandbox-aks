FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy workspace definition and shared packages
COPY go.work .
COPY pkg ./pkg

# Copy all services to satisfy go.work
COPY services ./services

# Build the service
WORKDIR /app/services/payment-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/services/payment-service/payment-service .
RUN addgroup -g 1001 -S appuser && adduser -S appuser -u 1001
USER appuser
EXPOSE 8081
CMD ["./payment-service"]