# start gRPC server and client access in localhost on go111 runtime

```
$ ./setup.sh
$ ./generate.sh
```

```
$ go run .

# call gRPC echo.Say via hand-written HTTP handler
$ curl -X POST http://localhost:8080/echo -H "Content-Type: application/json" --data '{"message_id":"foobar", "message_body": "hello"}'
{"message_id":"foobar","message_body":"hello","received":{"seconds":1559017407,"nanos":866193000}}

# call gRPC echo.Say via grpc-gateway
$ curl -X POST http://localhost:8080/v1/echo/say -H "Content-Type: application/json" --data '{"message_id":"foobar", "message_body": "hello"}'
{"message_id":"foobar","message_body":"hello","received":"2019-05-28T04:23:48.825457Z"}

# call gRPC echo.Say via GraphQL
$ curl -X POST http://localhost:8080/api/query -H "Content-Type: application/json" -H "Accept: application/json" --data '{"query":"mutation { say(input: { clientMutationId: \"foobar\", messageBody: \"hello\" }) { clientMutationId messageBody received } }"}'
{"data":{"say":{"clientMutationId":"foobar","messageBody":"hello","received":"2019-05-28T04:29:03.92021Z"}}}
```

```
$ gcloud --project ${PROJECT_ID} builds submit --tag gcr.io/${PROJECT_ID}/grpc-gqlgen
$ gcloud --project ${PROJECT_ID} beta run deploy grpc-gqlgen --image gcr.io/${PROJECT_ID}/grpc-gqlgen --region us-central1 --allow-unauthenticated
```
