REGISTRY_NAME?=docker.io/sbezverk
IMAGE_VERSION?=0.0.0

.PHONY: all jalapeno-gateway container push clean test

ifdef V
TESTARGS = -v -args -alsologtostderr -v 5
else
TESTARGS =
endif

all: jalapeno-gateway

jalapeno-gateway:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -ldflags '-extldflags "-static"' -o ./bin/jalapeno-gateway ./cmd/gateway.go

container: jalapeno-gateway
	docker build -t $(REGISTRY_NAME)/jalapeno-gateway-debug:$(IMAGE_VERSION) -f ./build/Dockerfile.gateway .

push: container
	docker push $(REGISTRY_NAME)/jalapeno-gateway-debug:$(IMAGE_VERSION)

clean:
	rm -rf bin

test:
	GO111MODULE=on go test `go list ./... | grep -v 'vendor'` $(TESTARGS)
	GO111MODULE=on go vet `go list ./... | grep -v vendor`
