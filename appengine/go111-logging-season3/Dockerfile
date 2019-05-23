FROM golang:1.12 as builder

WORKDIR /go/src/github.com/vvakame/til/appengine/go111-logging
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -v -o main


FROM alpine

COPY --from=builder /go/src/github.com/vvakame/til/appengine/go111-logging/main /main

CMD ["/main"]
