BINARY=cert-manager
PROJECT=arian-241419

# example make VERSION=v0.0.1
all: build tag push

build:
	docker build -t ${BINARY} .

tag:
	docker tag ${BINARY} gcr.io/${PROJECT}/${BINARY}:latest

push:
	gcloud docker -- push gcr.io/${PROJECT}/${BINARY}:latest

.PHONY: build tag push
