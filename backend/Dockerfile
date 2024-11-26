# Stage 1: Build the Go application
FROM golang:1.23-alpine as builder

# Set the current working directory in the container
WORKDIR /app

# Copy the go modules and source code
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# Stage 2: Create a lightweight image to run the app
FROM alpine:latest

# Install necessary certificates (for SSL, etc.)
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /app/myapp .

# Run the binary
CMD ["./myapp"]
