BIN := mackerel-remora
VERSION := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)

.PHONY: all
all: clean build

BUILD_LDFLAGS := "\
	-X main.version=$(VERSION) \
	-X main.revision=$(REVISION)"

.PHONY: deps
deps:
	go get -d ./...

.PHONY: build
build: deps
	go build -ldflags=$(BUILD_LDFLAGS) -o build/$(BIN) ./cmd/$(BIN)/...

.PHONY: test-deps
test-deps:
	go get -d -t ./...
	go get -u golang.org/x/lint/golint

.PHONY: test
test: test-deps
	go test -v ./...

.PHONY: lint
lint: test-deps
	go vet ./...
	golint -set_exit_status ./...

.PHONY: clean
clean:
	rm -fr build
	go clean

.PHONY: linux
linux:
	GOOS=linux go build -ldflags=$(BUILD_LDFLAGS) -o build/$(BIN) ./cmd/$(BIN)/...

.PHONY: docker
docker: linux
	docker build -t aknow/$(BIN):$(VERSION) -t $(BIN):$(VERSION) .

.PHONY: dockerhub
dockerhub:
	docker push aknow/mackerel-remora:$(VERSION)
	docker push aknow/mackerel-remora:latest

.PHONY: version
version:
	echo $(VERSION)
