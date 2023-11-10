# Образ для сборки
FROM golang:1.21-alpine AS builder

WORKDIR /go/src/app

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go/bin/telegram-bot ./cmd


# Образ для запуска
FROM alpine AS runner

COPY --from=builder /go/bin/telegram-bot /
COPY config/local.yaml /config/local.yaml

ENV CONFIG_PATH=/config/local.yaml

RUN mkdir /logs && touch /logs/logs.out
RUN mkdir /storage && touch /storage/storage.db

CMD ["/telegram-bot"]
