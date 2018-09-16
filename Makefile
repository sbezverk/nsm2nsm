REGISTRY_NAME = docker.io/sbezverk
IMAGE_VERSION = latest

.PHONY: clean test

ifdef V
TESTARGS = -v -args -alsologtostderr -v 5
else
TESTARGS =
endif

nsm2nsm-server:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ./bin/nsm2nsm-server ./cmd/nsm2nsm-server/nsm2nsm-server.go

nsm2nsm-client:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ./bin/nsm2nsm-client ./cmd/nsm2nsm-client/nsm2nsm-client.go

mac-nsm2nsm-server:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin go build -a -ldflags '-extldflags "-static"' -o ./bin/nsm2nsm-server.mac ./cmd/nsm2nsm-server/nsm2nsm-server.go

mac-nsm2nsm-client:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin go build -a -ldflags '-extldflags "-static"' -o ./bin/nsm2nsm-client.mac ./cmd/nsm2nsm-client/nsm2nsm-client.go

container-server: nsm2nsm-server
	docker build -t $(REGISTRY_NAME)/nsm2nsm-server:$(IMAGE_VERSION) -f ./Dockerfile.server .

container-client: nsm2nsm-client
	docker build -t $(REGISTRY_NAME)/nsm2nsm-client:$(IMAGE_VERSION) -f ./Dockerfile.client .

push: container-server container-client
	docker push $(REGISTRY_NAME)/nsm2nsm-server:$(IMAGE_VERSION)
	docker push $(REGISTRY_NAME)/nsm2nsm-client:$(IMAGE_VERSION)

clean:
	rm -rf bin

test:
	go test `go list ./... | grep -v 'vendor'` $(TESTARGS)
	go vet `go list ./... | grep -v vendor`
