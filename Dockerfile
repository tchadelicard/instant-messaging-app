# Stage 1: Build the binary
FROM golang:1.20 AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags '-extldflags "-static"' -o /build/instant-messaging-app main.go 

# Stage 2: Runtime
FROM alpine:3

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /build/instant-messaging-app .

# Expose the default port (API uses 8080, others can map differently in Compose)
EXPOSE 8080

# Default command; overridden by Docker Compose
CMD ["./instant-messaging-app"]
