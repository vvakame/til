# Google Cloud Functions を Go でやってみる

```bash
$ gcloud --project $PROJECT_ID functions deploy test-hello-world --runtime go111 --trigger-http --entry-point Func
```

Go 1.12 で作った go.mod をデプロイしようとすると怒られた…。

```bash
Deploying function (may take a while - up to 2 minutes)...failed.
ERROR: (gcloud.functions.deploy) OperationError: code=3, message=Build failed: go: finding github.com/vvakame/til/cloud-functions/go-hello v0.0.0
go build github.com/vvakame/til/cloud-functions/go-hello: module requires Go 1.12
```
