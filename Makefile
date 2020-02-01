REGISTRY_NAME?=docker.io/sbezverk
IMAGE_VERSION?=0.0.0

.PHONY: all jalapeno-gateway jalapeno-client compile-gateway compile-client container-gateway container-client push clean test

ifdef V
TESTARGS = -v -args -alsologtostderr -v 5
else
TESTARGS =
endif
BIN=./

all: jalapeno-gateway jalapeno-client

jalapeno-gateway:
	mkdir -p bin
	$(MAKE) -C ./cmd compile-gateway

jalapeno-client:
	mkdir -p bin
	$(MAKE) -C ./cmd/gateway-client compile-client

container-gateway: jalapeno-gateway
	docker build -t $(REGISTRY_NAME)/jalapeno-gateway-debug:$(IMAGE_VERSION) -f ./build/Dockerfile.gateway .

container-client: jalapeno-client
	docker build -t $(REGISTRY_NAME)/jalapeno-client:$(IMAGE_VERSION) -f ./build/Dockerfile.client .

push: container-gateway container-client
	docker push $(REGISTRY_NAME)/jalapeno-gateway-debug:$(IMAGE_VERSION)
	docker push $(REGISTRY_NAME)/jalapeno-client:$(IMAGE_VERSION)

clean:
	rm -rf bin

test:
	GO111MODULE=on go test `go list ./... | grep -v 'vendor'` $(TESTARGS)
	GO111MODULE=on go vet `go list ./... | grep -v vendor`
