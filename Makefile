PKG_NAME=confluence
BINARY_NAME=terraform-provider-confluence
INSTALL_DIR=$(HOME)/.terraform.d/plugins
VERSION=$$(cat VERSION)
TEST?=$$(go list ./...)
GOFMT_FILES?=$$(find . -name '*.go')

all: check test build

check: bin/golangci-lint
	bin/golangci-lint run

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

build: $(BINARY_NAME)

$(BINARY_NAME):
	go build -v -o $(BINARY_NAME)

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY_NAME}_${VERSION}_darwin_amd64
#	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY_NAME}_${VERSION}_freebsd_386
#	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY_NAME}_${VERSION}_freebsd_amd64
#	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY_NAME}_${VERSION}_freebsd_arm
#	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY_NAME}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY_NAME}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY_NAME}_${VERSION}_linux_arm
#	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY_NAME}_${VERSION}_openbsd_386
#	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY_NAME}_${VERSION}_openbsd_amd64
#	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY_NAME}_${VERSION}_solaris_amd64
#	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY_NAME}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY_NAME}_${VERSION}_windows_amd64

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 5m

fmt:
	gofmt -s -w $(GOFMT_FILES)

clean:
	rm -rf bin
	rm -rf site
	rm -f $(BINARY_NAME)

install: $(BINARY_NAME)
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)_v$(VERSION)

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

bin/golangci-lint:
	scripts/get-golangci.sh

.PHONY: all build check clean fmt install test testacc uninstall
