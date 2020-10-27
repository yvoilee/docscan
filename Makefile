BINARY_NAME=docscan
TAG:=$(shell git describe --tags)
LDFLAGS=-ldflags "-X main._VERSION_=$(TAG)"
build:
	@echo version tag: $(TAG)
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS}  -o $(BINARY_NAME) -v
build-arm64:
	@echo version tag: $(TAG)
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS}  -o $(BINARY_NAME) -v
build-ios:
	@echo version tag: $(TAG)
	go build ${LDFLAGS}  -o $(BINARY_NAME) -v
build-windows:
	@echo version tag: $(TAG)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS}  -o $(BINARY_NAME) -v