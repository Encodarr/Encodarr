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

# Copy go files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the source code
COPY . .

# Build the application
RUN go build -o /app/transfigurr ./cmd/transfigurr

# Stage 3: Final image
FROM alpine:latest
WORKDIR /app

# Setup non-root user first
RUN addgroup -g 1000 nonroot && \
    adduser -u 1000 -G nonroot -h /app -D nonroot

# Create all required directories
RUN mkdir -p \
    /config/db \
    /config/artwork \
    /movies \
    /series \
    /transcode \
    /app/frontend && \
    touch /config/restart.txt /config/shutdown.txt && \
    chmod 666 /config/restart.txt /config/shutdown.txt && \
    chown -R nonroot:nonroot \
    /app \
    /config \
    /movies \
    /series \
    /transcode



# Install runtime dependencies
RUN apk add --no-cache ffmpeg sqlite-libs

# Create frontend directory
RUN mkdir -p /app/frontend && \
    chown -R nonroot:nonroot /app

# Copy built artifacts from previous stages
COPY --from=frontend /frontend/dist /app/frontend/dist
COPY --from=backend /app/transfigurr /app/

# Copy initialization script
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh && \
    chown -R nonroot:nonroot /app


# Set default environment variables
ENV PUID=1000
ENV PGID=1000
ENV TZ=America/New_York

# Switch to non-root user
USER nonroot:nonroot


EXPOSE 7889

CMD ["/bin/ash", "/docker-entrypoint.sh"]