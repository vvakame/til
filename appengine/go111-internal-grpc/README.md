# start gRPC server and client access in localhost on go111 runtime

```
$ ./setup.sh
$ ./generate.sh
```

```
$ go run .

# call gRPC echo.Say via hand-written HTTP handler
$ curl -X POST http://localhost:8080/echo -H "Content-Type: application/json" --data '{"message_id":"foobar", "message_body": "hello"}'
{"message_id":"foobar","message_body":"hello","received":{"seconds":1544842978,"nanos":140999000}}

# call gRPC echo.Say via grpc-gateway
$ curl -X POST http://localhost:8080/v1/echo/say -H "Content-Type: application/json" --data '{"message_id":"foobar", "message_body": "hello"}'
{"message_id":"foobar","message_body":"hello","received":"2018-12-15T02:56:00.416423Z"}
```

```
$ gcloud --project $PROJECT_ID app deploy

$ curl -X POST https://go111-internal-grpc-dot-vvakame-playground.appspot.com/echo -H "Content-Type: application/json" --data '{"message_id":"foobar", "message_body": "hello"}'
{"message_id":"foobar","message_body":"hello","received":{"seconds":1544843639,"nanos":164338033}}

$ curl -X POST https://go111-internal-grpc-dot-vvakame-playground.appspot.com/v1/echo/say -H "Content-Type: application/json" --data '{"message_id":"foobar", "message_body": "hello"}'
{"message_id":"foobar","message_body":"hello","received":"2018-12-15T03:14:29.548454932Z"}
```
