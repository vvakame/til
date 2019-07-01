# gqlgen+gRPC integration

## やった&やりたい

1. 手でgRPC向けの.proto書く
1. 手で.protoにgqlgen向けのoptionsを付与する
1. .protoを元に *.graphql を自動生成する
1. .protoを元に *.gql.go を自動生成する
1. 通常のgqlgenに色々生成させる
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
* Relay Connection Spec周りの自動生成がかなり厳しそう
* Relay Global Object Identification Spec周りのかみ合わせ大変そう
    * .proto の package + ServiceName + ID をデフォにする…？
* 対応すべき.proto syntax etc
    * 各種数値型
    * oneof
    * map
    * reserved の活用(directive化とか？)
* Ruleのdestの構文は似非regexp replace的なやつじゃなくてtext/templateにしたほうが圧倒的によさそう
* 出てきたSchemaをlintに突っ込む
    * GetHoge とかを変換ルールで Hoge にすると hoge にならず Hoge のままになる… とかを倒したい
* 名前周り
    * やっぱりpackage名とService名を組み入れないとフィールド名被りが激しすぎる
    * デフォルトでは何も出力されない(ignoreされる)側に倒して開発者のお気持ちで公開したほうがいいのでは…
* Service毎skipするためにServiceOptionsの拡張も必要だわ…

### 所感

* gqlgenのPlugin使おうかと思ったけど使う必要があまりなさそう

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
