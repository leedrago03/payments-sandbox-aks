FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o reconciliation-service ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/reconciliation-service .
RUN addgroup -g 1001 -S appuser && adduser -S appuser -u 1001
USER appuser
EXPOSE 3007
CMD ["./reconciliation-service"]
