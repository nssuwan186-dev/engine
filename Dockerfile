# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /cli ./cmd/cli
RUN CGO_ENABLED=0 GOOS=linux go build -o /backdoor ./cmd/backdoor

# Frontend stage
FROM node:18-alpine AS frontend-builder

WORKDIR /app
COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

# Production stage
FROM alpine:latest

WORKDIR /app

# Copy binaries
COPY --from=builder /api /app/api
COPY --from=builder /cli /app/cli
COPY --from=builder /backdoor /app/backdoor

# Copy frontend
COPY --from=frontend-builder /app/dist /app/frontend/dist

# Copy static files
COPY --from=builder /app/internal/web/static /app/web/static

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

# Create directories
RUN mkdir -p /app/data /app/logs

# Environment
ENV PORT=8080
ENV DATABASE_PATH=/app/data/hotel.db
ENV LOG_PATH=/app/logs

# Expose
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Start
CMD ["/app/api"]
