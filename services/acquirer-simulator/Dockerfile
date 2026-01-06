FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy workspace definition and shared packages
COPY go.work .
COPY pkg ./pkg

# Copy all services to satisfy go.work
COPY services ./services

# Build the service
WORKDIR /app/services/acquirer-simulator
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o acquirer-simulator ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/services/acquirer-simulator/acquirer-simulator .
RUN addgroup -g 1001 -S appuser && adduser -S appuser -u 1001
USER appuser
EXPOSE 3004
CMD ["./acquirer-simulator"]
