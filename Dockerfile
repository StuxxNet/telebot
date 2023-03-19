### Build container
FROM --platform=amd64 golang:1.19.3-alpine3.16 AS build

WORKDIR /app
COPY . .

RUN go build -o telebot main.go
### Final container
FROM --platform=amd64 alpine:3.16

ENV TELEGRAM_BOT_TOKEN

COPY --from=build /app/telebot /usr/local/bin/telebot

EXPOSE 80

CMD ["/usr/local/bin/telebot"]