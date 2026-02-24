# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/server ./cmd/server/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/bin/server .
COPY --from=builder /app/config ./config
COPY --from=builder /app/public ./public

EXPOSE 8080

CMD ["./server"]