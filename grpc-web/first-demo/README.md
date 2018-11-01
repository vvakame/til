# grpc-webやってみた

結論：まだはやい gRPCみのある動くデモをdocker-compose的なものを使わずに作るのは今んとこ無理っぽい

```
$ brew install protobuf
$ git clone https://github.com/grpc/grpc-web.git && cd grpc-web && make install-plugin
$ cd client-web && npm install && npm run build
$ cd server && go run main.go
```
↑動くわけではない

https://github.com/grpc/grpc-web/commits/master 39560480263e0ab3e0cd02018a78ab88b95ad623

* メモ
    * https://github.com/grpc/grpc-web と https://github.com/improbable-eng/grpc-web/ は仕様が異なる
        * 今後互換性改善されるかもしれんけどわからん
        * https://github.com/improbable-eng/grpc-web/issues/199#issuecomment-409497863
    * Go単体でgrpc-webもサポートしているgRPCサーバを書くことは今のところできなさそう
        * Envoyないしnginxプラグインによる変換レイヤーが今のところ必要
* つらみ
    * https://github.com/grpc/grpc-web/issues/279
    * https://github.com/improbable-eng/grpc-web/issues/254
