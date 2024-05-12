# Stage 1: Build the React frontend
FROM node:alpine as frontend
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the FastAPI backend
FROM golang:alpine as backend
WORKDIR /src
COPY src /src
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


# Stage 3: Combine frontend and backend
FROM alpine:latest
WORKDIR /
COPY --from=frontend /frontend/dist /frontend/dist
COPY --from=backend /src /src

# Stage 4: Install ffmpeg
RUN apk add --no-cache --update \
    ffmpeg \
    && rm -rf /var/cache/apk/*

# Stage 5: Copy the init script and execute
WORKDIR /
COPY init /init
COPY /startup /startup

RUN chmod +x /init
ENV PUID=1000
ENV PGID=1000
ENV TZ=America/New_York
EXPOSE 7889
CMD ["/init"]