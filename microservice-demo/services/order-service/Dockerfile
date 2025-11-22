# Stage 1: Build
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Copy go.mod + go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy the shared pkg folder for local module replacement
COPY ../../pkg ./pkg


# Copy service source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 go build -o /auth-service ./cmd

# Stage 2: Run
FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /auth-service .
EXPOSE 8001
CMD ["./order-service"]
