FROM golang:1.22.3-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o loggerApp /app/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/loggerApp .
EXPOSE 5000
CMD ["./loggerApp"]