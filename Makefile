GO_VERSION=1.17

GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_VENDOR=$(GO_CMD) mod vendor

DOCKER_IMAGE_NAME=cenk1cekn2/nginx-ingress

BINARY_FOLDER=dist
BINARY_NAME=pipe

install:
	$(GO_VENDOR)

update:
	$(GO_GET) -u all
	$(GO_VENDOR)
	$(GO_CMD) mod tidy -compat=$(GO_VERSION)


all: test build

test:
	$(GO_TEST) -v ./...

clean:
	$(GO_CLEAN)
	rm -f $(BINARY_FOLDER)/$(BINARY_NAME)*

# Cross compilation

build: build-linux-amd64

build-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -mod=readonly -o $(BINARY_FOLDER)/$(BINARY_NAME)

.PHONY: build

build-docker: build build-docker-image

build-docker-image:
	docker build -t $(DOCKER_IMAGE_NAME):test .

run:
	./$(BINARY_FOLDER)/$(BINARY_NAME) $(ARGS)
