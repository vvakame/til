# Try go111 on GAE/Go

```
$ gcloud --project $PROJECT_ID app deploy
$ gcloud --project $PROJECT_ID beta tasks queues create-app-engine-queue go111-sample-queue

$ gcloud --project $PROJECT_ID beta tasks queues describe go111-sample-queue
```
