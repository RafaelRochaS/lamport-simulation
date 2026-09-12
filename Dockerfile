# Stage 1: Build the Go binary
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o simulation main.go

# Stage 2: Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/simulation .

# Create the logs folder inside the container
RUN mkdir logs

CMD ["./simulation"]
