FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o module2app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/module2app .
RUN chmod +x ./module2app
CMD ["./module2app"]