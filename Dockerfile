# Stage 1: Build the React frontend
FROM node:lts-alpine AS frontend
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm cache clean --force && npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.24.0-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/transfigurr ./cmd/transfigurr

# Stage 3: Final image
FROM alpine:latest
WORKDIR /app

# Set environment variables
ENV PUID=1000 \
    PGID=1000 \
    TZ=America/New_York

# Setup non-root user and install dependencies in one layer
RUN addgroup -g 1000 nonroot && \
    adduser -u 1000 -G nonroot -h /app -D nonroot && \
    apk add --no-cache ffmpeg sqlite-libs && \
    mkdir -p \
    /config/db \
    /config/artwork \
    /movies \
    /series \
    /transcode \
    /app/frontend && \
    chown -R nonroot:nonroot \
    /app \
    /config \
    /movies \
    /series \
    /transcode

# Copy built artifacts from previous stages and entrypoint script
COPY --chown=nonroot:nonroot --from=frontend /frontend/dist /app/frontend/dist
COPY --chown=nonroot:nonroot --from=backend /app/transfigurr /app/
COPY --chown=nonroot:nonroot docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

USER nonroot:nonroot
EXPOSE 7889

CMD ["/bin/ash", "/docker-entrypoint.sh"]