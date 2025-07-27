FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o api_gateway main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api_gateway .
RUN chmod +x ./api_gateway
CMD ["./api_gateway"]