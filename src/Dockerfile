# Stage 1: build
FROM golang:1.18.3-alpine3.16 AS builder

WORKDIR /app
RUN apk add --virtual build-dependencies build-base gcc wget git

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o myapp

# Stage 2: base image
FROM alpine:latest AS production
WORKDIR /app
COPY --from=builder /app/myapp /app/config.yaml /app/
CMD ["/app/myapp"]
