# AppEngine/SE go111 でstdoutでログ出してみる

```
$ gcloud --project $GOOGLE_CLOUD_PROJECT app deploy
```

```
$ gcloud --project $GOOGLE_CLOUD_PROJECT builds submit --tag gcr.io/$GOOGLE_CLOUD_PROJECT/go111-logging
$ gcloud --project $GOOGLE_CLOUD_PROJECT beta run deploy go111-logging --image gcr.io/$GOOGLE_CLOUD_PROJECT/go111-logging --region us-central1 --set-env-vars GOOGLE_CLOUD_PROJECT=$GOOGLE_CLOUD_PROJECT 
```

## Reference

https://twitter.com/furuyamayuki/statuses/1118382226205487106
https://github.com/yfuruyama/stackdriver-request-context-log
https://medium.com/google-cloud-jp/gae-%E3%81%AE%E3%83%AD%E3%82%B0%E3%81%AB%E6%86%A7%E3%82%8C%E3%81%A6-895ebe927c4

https://gcpug.slack.com/archives/C0D60LCAE/p1555837886044100

https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
https://cloud.google.com/logging/docs/agent/configuration#special-fields
https://cloud.google.com/appengine/docs/standard/go111/writing-application-logs

https://cloud.google.com/appengine/docs/standard/go112/go-differences

## わかったこと

* Request Log と App Log という異なる概念がある
* Request Log はAppEngineだと自動的に出る
* App Log で適切にtraceとか設定すればStackdriver Loggingでグルーピングはされる
* appengineのLog APIと異なるのは、app logでのseverityがrequest logにpropagateされない
    * なので、よくやってたWarningのapp logを含むrequest logの一覧を出す みたいなのができない
* Stackdriverが各ログエントリを名寄せしているのでフレッシュじゃないログだと完全性が保証されない可能性がある？

