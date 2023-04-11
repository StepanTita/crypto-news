FROM golang:1.20

WORKDIR /go/src/crypto-news
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/parser parser

###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/parser /usr/local/bin/parser
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["parser"]
