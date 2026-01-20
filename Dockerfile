# Stage 1: Build
FROM golang:tip-alpine3.23 AS builder

RUN apk add --no-cache gcc musl-dev git ca-certificates \
    && update-ca-certificates
# Install gcc for sqlite3
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary file. 
# CGO_ENABLED=1 need for sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o url-shortener ./cmd/url-shortener/main.go \
    && go clean -modcache

# Stage 2: Run
FROM alpine:3.23

RUN apk --no-cache add ca-certificates \
    && update-ca-certificates

WORKDIR /root/

# Copy only the assembled file from the first stage
COPY --from=builder /app/url-shortener .
# Copy the folder with configs (templates)
COPY --from=builder /app/config ./config

# Create a folder for the database
RUN mkdir -p ./storage

# Run application
CMD ["./url-shortener"]