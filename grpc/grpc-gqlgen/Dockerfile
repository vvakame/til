FROM golang:1.12 as builder

WORKDIR /work
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o exec-bin

FROM alpine

RUN apk --no-cache --update add ca-certificates
COPY --from=builder /work/exec-bin /exec-bin

CMD ["/exec-bin"]
