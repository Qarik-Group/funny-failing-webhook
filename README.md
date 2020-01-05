# Funny Failing Webhook for Kubernetes

This demonstration project shows a very simple Kubernetes webhook server, that is deployed publicly on Google Cloud Run, and sample configuration for your Kubernetes cluster.

Any requests to CREATE pods into a specific namespace will fail, with a random funny message.

This example webhook server assumes it has HTTPS provided for it by Cloud Run, and does not need to talk to your Kubernetes cluster, so its much simpler to write than common webooks that might want to interact with your Kubernetes cluster, or might need to manage their own TLS certificates for HTTPS.
