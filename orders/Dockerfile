# Build environment
# Use a Go base image to build the application
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.work and go.sum files to cache dependencies
COPY go.work go.work.sum ./

# Copy all go.mod files first for better caching
COPY users/go.mod users/go.sum ./users/
COPY products/go.mod products/go.sum ./products/
COPY orders/go.mod orders/go.sum ./orders/
COPY utils/go.mod utils/go.sum ./utils/

# Download all dependencies
RUN go work sync
RUN go mod download

# Copy the rest of the source code
COPY users/ ./users/
COPY products/ ./products/
COPY orders/ ./orders/
COPY utils/ ./utils/

# Build the Go app as a static binary for small size
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./orders

# Production Environment
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Set working directory
COPY --from=builder /app/main .

# Expose REST port
EXPOSE 8282

# The command to run the application
CMD ["./main"]