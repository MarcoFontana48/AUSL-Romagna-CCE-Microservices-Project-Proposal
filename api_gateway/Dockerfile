FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY api_gateway/ ./
COPY utils/ ../utils/
RUN go mod tidy
RUN go mod download
RUN go build -o api_gateway main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api_gateway .
CMD ["./api_gateway"]