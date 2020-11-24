GO := go
GOFMT := gofmt
ARCH ?= $(shell go env GOARCH)
OS ?= $(shell go env GOOS)

.PHONY: all

IMAGE_REPOSITORY_URL = registry.wanzhixin.net/fanli/rebate
BUILD_NUMBER = $(shell git rev-parse --short HEAD)
BUILD_TAG := $(shell date +%Y%m%d%H%M%S)
COMMAND_ROOT = dev-gitlab.wanxingrowth.com/fanli/rebate/cmd

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

all: build test

build: build_service_cross_only

.PHONY: test
test:
	mkdir -p .test-result
	go test -cover -coverprofile cover.out -outputdir .test-result ./...
	go tool cover -html=.test-result/cover.out -o .test-result/coverage.html

.PHONY: build_service_local
build_service_local: build_service_protos
	mkdir -p builds/debug
	go build -o builds/debug/service -ldflags '${EXTLDFLAGS}-X dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/version.VersionDev=build.$(BUILD_NUMBER)' $(COMMAND_ROOT)

.PHONY: build_service_cross
build_service_cross: build_service_protos
	mkdir -p builds/release
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o builds/release/service -ldflags '${EXTLDFLAGS}-X dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/version.VersionDev=build.$(BUILD_NUMBER)' $(COMMAND_ROOT)

.PHONY: build_service_cross_only
build_service_cross_only:
	mkdir -p builds/release
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o builds/release/service -ldflags '${EXTLDFLAGS}-X dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/version.VersionDev=build.$(BUILD_NUMBER)' $(COMMAND_ROOT)

.PHONY: run_service_local
run_service_local: build_service_local
	# Use ./config/example.yaml copy to ./config/dev.yaml when first time debug
	builds/debug/service --config=./config/dev.yaml

build_service_image:
	docker build -t $(IMAGE_REPOSITORY_URL):$(BUILD_TAG) -f builds/docker/Dockerfile .

push_service_image: build_service_image
	docker push $(IMAGE_REPOSITORY_URL):$(BUILD_TAG)

push_service_image_clear_local: push_service_image
	docker rmi $(IMAGE_REPOSITORY_URL):$(BUILD_TAG)

push_service_image_clear_local_write_url: push_service_image_clear_local
	echo "$(IMAGE_REPOSITORY_URL):$(BUILD_TAG)" > builds/docker/latest-image

.PHONY: build_service_protos
build_service_protos:
	@sh builds/protos/protos.sh

.PHONY: pre_commit
pre_commit: test
	go fmt ./...
	go vet ./...
