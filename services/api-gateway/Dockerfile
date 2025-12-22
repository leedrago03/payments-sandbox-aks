FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-gateway .
RUN addgroup -g 1001 -S appuser && adduser -S appuser -u 1001
USER appuser
EXPOSE 8000
CMD ["./api-gateway"]