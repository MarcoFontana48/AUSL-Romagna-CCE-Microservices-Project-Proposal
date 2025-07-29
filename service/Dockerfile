FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/service
RUN go mod tidy
RUN go mod download
RUN go build -o service main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/service/service .
CMD ["./service"]