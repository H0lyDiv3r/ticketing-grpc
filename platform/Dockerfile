# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build platform binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/platform-service ./platform/cmd

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /bin/platform-service /app/platform-service

EXPOSE 8080

ENTRYPOINT ["/app/platform-service"]
