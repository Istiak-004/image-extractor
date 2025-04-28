FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/bin/server ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin/server /app/server
RUN chmod +x /app/server

CMD ["/app/server"]
