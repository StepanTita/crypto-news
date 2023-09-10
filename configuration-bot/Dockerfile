FROM golang:1.20

WORKDIR /go/src/github.com/crypto-news
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/configuration-bot configuration-bot

###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/configuration-bot /usr/local/bin/configuration-bot
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["configuration-bot"]
