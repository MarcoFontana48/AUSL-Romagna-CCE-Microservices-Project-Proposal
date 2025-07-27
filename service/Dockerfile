FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY service/go.mod ./
COPY utils/ ../utils/
COPY service/main.go .
RUN go build -o service main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/service .
CMD ["./service"]