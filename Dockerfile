FROM golang:latest

LABEL maintainer="PeterMBenjamin@gmail.com"

WORKDIR /go/src/github.com/petermbenjamin/slackr

COPY . /go/src/github.com/petermbenjamin/slackr

CMD ["/bin/bash", "-c"]

