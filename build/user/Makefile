IMAGE ?= awesome-crypto-user
PROJECT ?= obarbier
DOCKER_REGISTRY ?= docker.io
# Declaring and calculating the version (tag) for the docker image in the format: <data>-<commit>
VERSION ?= $(shell date +v%Y%m%d)-$(shell git describe --tags --always)

build-image:
	docker build -t $(DOCKER_REGISTRY)/$(PROJECT)/$(IMAGE):$(VERSION) \
		-f Dockerfile ../../ --pull

build-push: build-image
	docker push $(DOCKER_REGISTRY)/$(PROJECT)/$(IMAGE):$(VERSION)
	docker tag $(DOCKER_REGISTRY)/$(PROJECT)/$(IMAGE):$(VERSION) $(DOCKER_REGISTRY)/$(PROJECT)/$(IMAGE):latest
	docker push $(DOCKER_REGISTRY)/$(PROJECT)/$(IMAGE):latest

.PHONY: build-image build-push