FROM golang:1.20

WORKDIR /go/src/github.com/crypto-news
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/twitter-bot twitter-bot

###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/twitter-bot /usr/local/bin/twitter-bot
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["twitter-bot"]
