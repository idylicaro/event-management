# Build stage: Compile the Go application
FROM golang:1.23-bullseye AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules to the container and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code into the container
COPY . ./

# Build the Go binary
RUN go build -o main .

# Production stage: Create a lightweight runtime image for production
FROM debian:bullseye-slim AS production

# Install necessary libraries (e.g., libc6) to run the Go binary
RUN apt-get update && \
    apt-get install -y libc6

# Set the working directory for the application in the container
WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main .

# Set the environment variable for production
ENV ENVIRONMENT=production

# Expose port 8080 for the application
EXPOSE 8080

# Make the Go binary executable
RUN chmod +x /app/main

# Set the entry point to the Go binary for production
ENTRYPOINT ["/app/main"]

# Development stage: Create a lightweight runtime image for development
FROM debian:bullseye-slim AS development

# Install necessary libraries (e.g., libc6) to run the Go binary
RUN apt-get update && \
    apt-get install -y libc6

# Set the working directory for the application in the container
WORKDIR /app

# Copy the environment configuration file for development
COPY .env .env

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main .

# Set the environment variable for development
ENV ENVIRONMENT=development

# Expose port 8080 for the application
EXPOSE 8080

# Make the Go binary executable
RUN chmod +x /app/main

# Set the entry point to the Go binary for development
ENTRYPOINT ["/app/main"]
