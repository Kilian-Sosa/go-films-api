FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server/main.go
RUN go build -o migrate ./cmd/migrate/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["./server"]
    