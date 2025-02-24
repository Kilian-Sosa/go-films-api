FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# For debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go build -gcflags="all=-N -l" -o server ./cmd/server/main.go

# For production
# RUN go build -o server ./cmd/server/main.go

RUN go build -o migrate ./cmd/migrate/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/migrations /app/migrations

COPY docs ./docs

# For debugging
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv

EXPOSE 8080 2345

CMD ["./server"]
    