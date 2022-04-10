# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVENDOR=$(GOCMD) mod vendor
BINARY_FOLDER=dist
BINARY_NAME=pipe

install:
	$(GOVENDOR)

update:
	$(GOGET) -u all
	$(GOVENDOR)
	$(GOCMD) mod tidy -compat=1.17


all: test build

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_FOLDER)/$(BINARY_NAME)*

# Cross compilation

build: build-linux-amd64

build-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -mod=readonly -o $(BINARY_FOLDER)/$(BINARY_NAME)

.PHONY: build
