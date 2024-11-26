# Step 1: Build the Go application
FROM golang:1.20-alpine as builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN GOOS=linux GOARCH=amd64 go build -o main .

# Step 2: Create a small image to run the application
FROM alpine:latest

WORKDIR /root/

# Install necessary dependencies (like ca-certificates)
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port (use the PORT environment variable in Railway)
EXPOSE 8080

# Start the app
CMD ["./main"]
