FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/api_gateway
RUN go mod tidy
RUN go mod download
RUN go build -o api_gateway main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api_gateway .
CMD ["./api_gateway"]