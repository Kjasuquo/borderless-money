# Build Stage
FROM golang:1.23-alpine AS build

WORKDIR /app
COPY . .

# Install necessary dependencies (g++ for any CGO requirements)
RUN apk update && apk add --no-cache g++

# Build the binary for Linux AMD64 (for GKE)
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main cmd/main.go

# Deploy Stage
FROM alpine:3.13
WORKDIR /app
COPY --from=build /app/main .

EXPOSE 8080
CMD ["/app/main"]



