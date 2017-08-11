FROM golang:latest

LABEL maintainer="PeterMBenjamin@gmail.com"

WORKDIR /go/src/github.com/petermbenjamin/slackr
COPY . /go/src/github.com/petermbenjamin/slackr
RUN go get github.com/golang/dep/cmd/dep
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o slackr .

FROM alpine:latest
WORKDIR /app/
COPY --from=0 /go/src/github.com/petermbenjamin/slackr/slackr .
ENTRYPOINT ["./slackr"]
