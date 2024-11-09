# Use a Golang image for building the application
FROM golang:1.20 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application from the cmd directory
RUN cd cmd && go build -o ausf

# Use a base image with the required GLIBC version
FROM debian:bookworm-slim

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder
COPY --from=builder /app/cmd/ausf .

# Expose necessary ports
EXPOSE 8000

# Set the entrypoint
ENTRYPOINT ["./ausf"]