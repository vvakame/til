# gqlgen+gRPC integration

## やった&やりたい

1. 手でgRPC向けの.proto書く
1. 手で.protoにgqlgen向けのoptionsを付与する
1. .protoを元に *.graphql を自動生成する(今は手)
1. 通常のgqlgenに色々生成させる
1. .protoを元に *glue.go を自動生成する(今は手)
1. ResolverRootを組み立てる
    * めんどいのでwireを使った

### 課題メモ

* 複数のgRPC Serverの結果の組み合わせ
* 手で書いたほうがよい部分を自動生成部分から切り離し上書き可能にする
* gRPCから帰ってきたエラーのハンドリング
* optionsの充実
    * 特定のフィールドを隠すとか
* directiveの対応
* 共通定義の扱い
* gqlgen.yml の更新など？
* message in message とか .protoの仕様がわりと沼
* Relay Connection Spec周りの自動生成がかなり厳しそう
* Relay Global Object Identification Spec周りのかみ合わせ大変そう
    * .proto の package + ServiceName + ID をデフォにする…？
* 対応すべき.proto syntax etc
    * 各種数値型
    * oneof
    * map
    * reserved の活用(directive化とか？)

### 所感

* gqlgenのPlugin使おうかと思ったけど使う必要があまりなさそう
* gqlgenのAST組んだらschemaにdumpしてくれる機能ほしい…
    * introspection schema経由するとdirectiveの情報が落ちてしまう…

## 作業メモ


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
