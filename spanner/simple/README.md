# Spanner素振り

* https://cloud.google.com/spanner/docs/getting-started/go/?hl=ja
* https://github.com/GoogleCloudPlatform/golang-samples/tree/7e2f8f5bf5a4825dd62305d69a31f7f30dc6b7f7/spanner/spanner_snippets

社内に立ってる共用のSpannerインスタンスに繋いであれこれしてみる。
プロジェクト構成がちょっとアレなのでメモしておく。

* Project A (Spannerが立ってるプロジェクト)
* Project B (僕の個人プロジェクト)

BでServiceAccount作ってAのプロジェクトに登録&権限貰って、SA経由でSpannerにアクセスする。

利用するのはAのSpannerだが、使うにはBのプロジェクトでもCloud SpannerのAPIが有効になっている必要がある。
次のようなエラーが出たら素直に表示されたURLにアクセスしてAPIを有効にしてリトライする。

```
2018/04/12 13:11:18 rpc error: code = PermissionDenied desc = Cloud Spanner API has not been used in project xxxx before or it is disabled. Enable it by visiting https://console.developers.google.com/apis/api/spanner.googleapis.com/overview?project=xxxx then retry. If you enabled this API recently, wait a few minutes for the action to propagate to our systems and retry.
```
