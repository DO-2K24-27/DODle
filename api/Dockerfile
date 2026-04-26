# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add curl

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/data ./data

# Expose port
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "curl -f http://localhost:8080 || exit 1" ]

# Run the application
CMD ["./main"] 