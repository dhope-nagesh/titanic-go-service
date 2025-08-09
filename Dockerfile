# Stage 1: Build the application on a standard Debian-based image
FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Copy dependency files and download them
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary. CGO is enabled by default.
RUN go build -v -o /app/server ./cmd/server


# Stage 2: Create the final production image
# Use a slim Debian image that is compatible with the build environment
FROM debian:bullseye-slim

WORKDIR /root/

# Copy the pre-built binary from the 'builder' stage
COPY --from=builder /app/server .

# Expose the port the app will run on
EXPOSE 8080

# The command to run when the container starts
CMD ["./server"]
