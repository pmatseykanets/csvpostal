BINARY=csvpostal
PLATFORM=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

default: build

build:
	go fmt ./...
	go vet ./...	
	@mkdir -p ./release
	@rm -rf ./release/*
	GOOS=$(PLATFORM) GOARCH=$(ARCH) CGO_ENABLED=1 go build -o ./release/$(BINARY)-$(PLATFORM)-$(ARCH)
	cp ./release/$(BINARY)-$(PLATFORM)-$(ARCH) ./$(BINARY)

install: build
	cp ./$(BINARY) $(GOBIN)

.PHONY: build install
