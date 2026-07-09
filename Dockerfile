FROM golang:1.25 AS builder

RUN apt-get update && apt-get install -y git make && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
ENV GOPROXY=direct
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -a -installsuffix cgo \
    -o /app/bin/subscription-service ./cmd/main.go



FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app


COPY --from=builder /app/bin/subscription-service /app/subscription-service

COPY migrations /app/migrations

COPY --from=builder /app/docs /app/docs

COPY .env /app/.env

RUN chmod +x /app/subscription-service

EXPOSE 8080

CMD ["/app/subscription-service"]