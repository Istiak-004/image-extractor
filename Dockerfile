FROM golang:1.24.2-alpine
WORKDIR /app
COPY . .
RUN go build -o bin/server ./cmd/api
CMD ["./bin/server"]
