# stdoutでログ出力 シーズン3

```
$ gcloud --project $GOOGLE_CLOUD_PROJECT app deploy
```

```
$ gcloud --project $GOOGLE_CLOUD_PROJECT builds submit --tag gcr.io/$GOOGLE_CLOUD_PROJECT/go111-logging-s3
$ gcloud --project $GOOGLE_CLOUD_PROJECT beta run deploy go111-logging --image gcr.io/$GOOGLE_CLOUD_PROJECT/go111-logging-s3 --region us-central1 --set-env-vars GOOGLE_CLOUD_PROJECT=$GOOGLE_CLOUD_PROJECT
```

https://github.com/vvakame/til/pull/19
https://github.com/vvakame/til/pull/27
https://github.com/vvakame/til/pull/28
