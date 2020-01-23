GCLOUD_RUN_NAME=funny-failing-webhook
GCLOUD_IMAGE?=gcr.io/drnic-257704/funny-failing-webhook
GCLOUD_REGION?=us-central1

.PHONY: all docker cloudrun

all: docker cloudrun

docker:
	docker build -t $(GCLOUD_IMAGE) .
	docker push $(GCLOUD_IMAGE)

cloudrun:
	gcloud run deploy $(GCLOUD_RUN_NAME) \
		--image $(GCLOUD_IMAGE) \
		--platform managed \
		--region $(GCLOUD_REGION) \
		--allow-unauthenticated
