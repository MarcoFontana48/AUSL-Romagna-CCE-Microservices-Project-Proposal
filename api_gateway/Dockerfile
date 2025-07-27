FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY api_gateway/go.mod ./
COPY utils/ ../utils/
COPY api_gateway/main.go .
RUN go build -o api_gateway main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api_gateway .
CMD ["./api_gateway"]