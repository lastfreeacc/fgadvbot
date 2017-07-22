FROM golang:alpine

# add user
RUN \
 adduser -D bot && \
 mkdir -p /bot && \
 chown -R bot:bot /bot

# build bot
ADD . /go/src/github.com/lastfreeacc/fgadvbot
RUN \
 apk add --update --no-progress git && \
 cd /go/src/github.com/lastfreeacc/fgadvbot && \
 go get && \
 go build -o /bot/fgadvbot && \
 apk del git && \
 rm -rf /go/src/* && \
 rm -rf /var/cache/apk/*

RUN \
 echo "#!/bin/sh" > /bot/exec.sh && \
 echo "/bot/fgadvbot" >> /bot/exec.sh && \
 chmod +x /bot/exec.sh

USER bot
WORKDIR /bot
ENTRYPOINT ["/bot/exec.sh"]