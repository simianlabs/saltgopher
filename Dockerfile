FROM golang:alpine
RUN mkdir /app 
ADD . /app/
WORKDIR /app 
RUN go get github.com/nlopes/slack && \
    go get github.com/Jeffail/gabs
RUN go build -o saltgopher .
RUN adduser -S -D -H -h /app saltgopher
USER saltgopher
CMD ["./saltgopher"]