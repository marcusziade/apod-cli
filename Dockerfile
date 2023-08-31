# Use the official Go image
FROM golang:1.19 AS build

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Use a lightweight image for deployment
FROM alpine:latest AS runtime

# Install necessary libraries
RUN apk --no-cache add ca-certificates

# Copy binary
COPY --from=build /app/main /app/

# Set the working directory
WORKDIR /app

# Run the application
CMD [". --download-only"]
