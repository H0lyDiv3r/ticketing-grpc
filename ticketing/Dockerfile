# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Copy source code including local vendor directory
COPY . .

# Build service binary offline using vendored dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-w -s" -o /bin/ticketing-service ./ticketing/cmd

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /bin/ticketing-service /app/ticketing-service

EXPOSE 50052

ENTRYPOINT ["/app/ticketing-service"]
