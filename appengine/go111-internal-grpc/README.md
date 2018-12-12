# start gRPC server and client access in localhost on go111 runtime

```
$ go get -u github.com/kazegusuri/grpcurl
$ go run .

$ echo '{"message_id":"foobar", "message_body": "hello"}' | grpcurl -k call localhost:5000 vvakame.echo.Echo.Say
{"message_id":"foobar","message_body":"hello","received":"2018-12-12T09:12:10.157027Z"}

$ curl -X POST http://localhost:8080/echo -H "Content-Type: application/json" --data '{"message_id":"foobar", "message_body": "hello"}'
{"message_id":"foobar","message_body":"hello","received":{"seconds":1544606049,"nanos":877374000}}

$ curl -X POST https://go111-internal-grpc-dot-vvakame-playground.appspot.com/echo -H "Content-Type: application/json" --data '{"message_id":"foobar", "message_body": "hello"}'
{"message_id":"foobar","message_body":"hello","received":{"seconds":1544606049,"nanos":877374000}}
```

```
$ gcloud --project $PROJECT_ID app deploy
```
