# Stage 1: Build stage
FROM golang:1.22.3-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all the files to the container
COPY . .

# Perform go mod tidy and go mod vendor
RUN go mod tidy
RUN go mod vendor

# Build the binary
RUN go build -o myapp .

####################################################################
# Stage 2: Production stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/myapp .

# This is the port that our application will be listening on.
EXPOSE 4100

# This is the command that will be executed when the container is started.
ENTRYPOINT ["./myapp"]