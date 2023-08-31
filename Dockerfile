# Use the official Golang image to create a build artifact.
FROM golang:1.19 as builder

# Set the working directory
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /app/main .

# Use a minimal image for deployment
FROM alpine:latest

# Copy the binary
COPY --from=builder /app/main /app/main

# This is the entry point, specify your flags here
ENTRYPOINT [ "/app/main", "--download-only" ]
