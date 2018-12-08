FROM golang:alpine
RUN mkdir /app 
ADD . /app/
WORKDIR /app 
RUN apk update && \
    apk upgrade && \
    apk add git && \
    go get github.com/nlopes/slack && \
    go get github.com/Jeffail/gabs && \
    go build -o saltgopher .  && \
    mv config.json.example config.json && \
    mv roles.json.example roles.json && \
    adduser -S -D -H -h /app saltgopher
USER saltgopher
CMD ["./saltgopher"]