# Stage 1: Build the React frontend
FROM node:lts-alpine AS frontend
WORKDIR /frontend
COPY frontend/package*.json ./

# Clean npm cache
RUN npm cache clean --force

# Install dependencies
RUN npm ci

COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Go backend with CGO enabled
FROM golang:1.23-alpine AS backend
WORKDIR /src

# Install necessary dependencies for CGO
RUN apk add --no-cache gcc musl-dev

COPY src /src
RUN go mod download

# Enable CGO and build the Go application
ENV CGO_ENABLED=1
RUN go build -o main .

# Stage 3: Combine frontend and backend
FROM alpine:latest
WORKDIR /

# Install SQLite3 library
RUN apk add --no-cache sqlite-libs

COPY --from=frontend /frontend/dist /frontend/dist
COPY --from=backend /src/main /src/main

# Stage 4: Install ffmpeg
RUN apk add --no-cache --update \
    ffmpeg \
    && rm -rf /var/cache/apk/*

# Stage 5: Copy the init script and execute
WORKDIR /
COPY init /init
RUN chmod +x /init

ENV PUID=1000
ENV PGID=1000
ENV TZ=America/New_York
EXPOSE 7889

# Use ash to run the init script
CMD ["/bin/ash", "/init"]