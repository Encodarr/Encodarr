# Stage 1: Build the React frontend
FROM node:lts-alpine AS frontend
WORKDIR /frontend
COPY frontend/package*.json ./

# Clean npm cache and install dependencies
RUN npm cache clean --force && npm ci

# Copy frontend source and build
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.23.4-alpine AS backend
WORKDIR /app

# Install necessary dependencies for CGO
RUN apk add --no-cache gcc musl-dev

# Copy go files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the source code
COPY . .

# Build for target platform
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=1
ENV GOOS=$TARGETOS 
ENV GOARCH=$TARGETARCH

# Install cross-compiler for ARM64 if needed
RUN if [ "$TARGETARCH" = "arm64" ]; then \
    apt-get update && apt-get install -y gcc-aarch64-linux-gnu \
    && rm -rf /var/lib/apt/lists/* ; \
    export CC=aarch64-linux-gnu-gcc; \
    fi

# Build the application
RUN go build -o /app/transfigurr ./cmd/transfigurr

# Stage 3: Final image
FROM alpine:latest
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ffmpeg sqlite-libs

# Create frontend directory
RUN mkdir -p /app/frontend

# Copy built artifacts from previous stages
COPY --from=frontend /frontend/dist /app/frontend/dist
COPY --from=backend /app/transfigurr /app/

# Copy initialization script
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# Set default environment variables
ENV PUID=1000
ENV PGID=1000
ENV TZ=America/New_York

# Switch to non-root user
USER nonroot

EXPOSE 7889

CMD ["/bin/ash", "/docker-entrypoint.sh"]