# ============================================
# Stage 1: Build
# ============================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Copy go.mod dan go.sum dulu untuk cache dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary — CGO disabled untuk static binary, strip debug info untuk ukuran kecil
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/gocatat-api .

# ============================================
# Stage 2: Production (minimal image ~10MB)
# ============================================
FROM alpine:3.21

WORKDIR /app

# Install ca-certificates untuk HTTPS connections (ke DB, external APIs)
RUN apk add --no-cache ca-certificates tzdata

# Set timezone ke WIB (opsional, sesuaikan)
ENV TZ=Asia/Jakarta

# Copy binary dari builder stage
COPY --from=builder /app/gocatat-api .

# EasyPanel akan inject env vars lewat dashboard,
# jadi TIDAK perlu copy .env ke dalam container

# Port yang diexpose — sesuaikan dengan APP_PORT di env
EXPOSE 8080

# Jalankan binary
CMD ["./gocatat-api"]
