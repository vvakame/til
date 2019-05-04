# ory/fosite + Cloud Datastore integration

based on https://github.com/ory/fosite-example

```
$ gcloud builds submit --tag gcr.io/${PROJECT_ID}/oauth2idp-example
$ gcloud beta run deploy --image gcr.io/${PROJECT_ID}/oauth2idp-example --set-env-vars=DATASTORE_PROJECT_ID=${PROJECT_ID},BASE_URL=https://foo-bar.a.run.app
```

```
$ ./serve.sh
$ go run main.go
```
