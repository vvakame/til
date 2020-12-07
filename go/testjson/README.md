# go test の machine-readable な出力を考える

https://nihonbuson.hatenadiary.jp/entry/2018/03/10/110000
この記事読んで面白いなーと思ったので

* https://github.com/golang/go/issues/2981
* https://golang.org/cmd/test2json/

```shell
$ go test -v ./...
$ go test -json ./... | go run ./cmd/json2result
```

なんか stdout 経由の出力が同じにならないな

* https://github.com/golang/go/issues/28953

https://github.com/gotestyourself/gotestsum っていうのがあるよって教えてもらった
とりあえずこれでできる範囲で頑張るのがよいのでは…？
