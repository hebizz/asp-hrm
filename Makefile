IMAGE_COMMON=harbor.jiangxingai.com/library/
#IMAGE_COMMON=registry.jiangxingai.com:5000/
APP_NAME=asp-hrm
VER=$(shell cat ./VERSION)

PLATA32=arm32v7
PLATA64=arm64v8
PLATX64=x8664

.PHONY: build test clean prepare update docker

GO = CGO_ENABLED=0 GO111MODULE=on go

MICROSERVICES=cmd/asp-hrm

#.PHONY: $(MICROSERVICES)

#DOCKERS=docker_asp-hrm
#.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 1.0.0)
BUILD_D=$(shell date "+%m/%d/%Y %R %Z")
GIT_SHA=$(shell git rev-parse --short HEAD)
GO_VERS=$(shell go version)

PROJ=github.com/greystone/asp-hrm

LDFLAG=-ldflags \"-X "$(PROJ)/pkg/config.BuildVersion=$(VERSION)" \
                 -X "$(PROJ)/pkg/config.BuildTime=$(BUILD_D)" \
                 -X "$(PROJ)/pkg/config.BuildHash=$(GIT_SHA)" \
                 -X "$(PROJ)/pkg/config.GoVersion=$(GO_VERS)" \"
.PHONY: asp-hrm

asp-hrm: main.go

	GOOS=linux GOARCH=amd64  $(GO) build  $(LDFLAG_ARM) -o asp-hrm

asp-hrm-mac: main.go

	GOOS=darwin GOARCH=amd64  $(GO) build  $(LDFLAG_ARM) -o asp-hrm

asp-hrm-arm: main.go

	GOOS=linux GOARCH=arm64  $(GO) build  $(LDFLAG_ARM) -o asp-hrm-arm


build: $(MICROSERVICES)
	$(GO) build ./...

cmd/asp-hrm:
	$(GO) build $(GOFLAGS) -o $@ ./cmd

test:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
	./bin/test-attribution-txt.sh
	./bin/test-go-mod-tidy.sh

clean:
	rm -f $(MICROSERVICES)

build-docker: asp-hrm
	docker build -t $(IMAGE_FULL) .

push:
	docker push $(IMAGE_FULL)

x86: PLAT=$(PLATX64)
x86: IMAGE_FULL=$(IMAGE_COMMON)$(APP_NAME)/$(PLAT)/others:$(VER)
x86: build-docker push