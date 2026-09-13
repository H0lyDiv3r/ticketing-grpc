# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build service binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/payment-service ./payment/cmd

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /bin/payment-service /app/payment-service

EXPOSE 50053

ENTRYPOINT ["/app/payment-service"]
