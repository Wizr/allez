GO = go
GLIDE = glide

.PHONY: all
all: darwin linux

glide.lock: glide.yaml
	$(GLIDE) update

.PHONY: darwin
darwin: main.go core libs toolez glide.lock static
	GOOS=darwin GOARCH=amd64 $(GO) build -o bin/allez.darwin main.go

.PHONY: linux
linux: main.go core libs toolez glide.lock static
	GOOS=linux GOARCH=amd64 $(GO) build -o bin/allez.linux main.go

# set npm/yarn to use taobao mirror
# npm config set registry https://registry.npm.taobao.org
static: client client/src client/src/**
	cd client && yarn install && yarn build
