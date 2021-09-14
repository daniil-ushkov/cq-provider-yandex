# syntax=docker/dockerfile:1

FROM golang:1.16
WORKDIR /app

# Copy YC provider
COPY . .

CMD go test -v ./...
