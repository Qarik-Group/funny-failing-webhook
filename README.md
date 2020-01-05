# Funny Failing Webhook for Kubernetes

This demonstration project shows a very simple Kubernetes webhook server, that is deployed publicly on Google Cloud Run, and sample configuration for your Kubernetes cluster.

Any requests to CREATE pods into a specific namespace will fail, with a random funny message.

This example webhook server assumes it has HTTPS provided for it by Cloud Run, and does not need to talk to your Kubernetes cluster, so its much simpler to write than common webooks that might want to interact with your Kubernetes cluster, or might need to manage their own TLS certificates for HTTPS.

## Deploy to Cloud Run

Whilst this service is already running on Cloud Run, you might want to fork and deploy it yourself. Or, more likely, I might want to do this in the future and need documentation. Luckily, I wrote myself the following documentation.

At the time of writing Google Cloud Run requires OCIs to be hosted on [Google Container Registry](https://console.cloud.google.com/gcr) (GCR):

```plain
docker build -t gcr.io/drnic-257704/funny-failing-webhook .
docker push gcr.io/drnic-257704/funny-failing-webhook
```

To deploy as an unauthenticated service to Google Cloud Run:

```plain
gcloud run deploy funny-failing-webhook \
    --image gcr.io/drnic-257704/funny-failing-webhook \
    --platform managed \
    --region us-central1 \
    --allow-unauthenticated
```

To see the available Cloud Run service and its URL:

```plain
$ gcloud run services list --platform managed
   SERVICE                REGION       URL                                                    LAST DEPLOYED BY         LAST DEPLOYED AT
✔  funny-failing-webhook  us-central1  https://funny-failing-webhook-lg2hslfa4a-uc.a.run.app  drnic@starkandwayne.com  2020-01-05T22:25:49.899Z
...
```

We use https://funny-failing-webhook-lg2hslfa4a-uc.a.run.app in our webhook configuration.

Quick confirmation that our HTTP server can receive requests by hitting the `/` endpoint:

```plain
$ curl https://funny-failing-webhook-lg2hslfa4a-uc.a.run.app
Funny Failing Webhook always rejects pod CREATE requests
Available routes:
/
/healthz
/funny-failing-webhook
```
