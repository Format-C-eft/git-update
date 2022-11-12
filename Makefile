export GO111MODULE=on
export GOSUMDB=off

LOCAL_BIN	= $(CURDIR)/bin
GOBIN		= $(GOPATH)/bin
GOARCH		= amd64

.PHONY: bin-deps
bin-deps:
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

.PHONY: .lint
.lint:
	$(LOCAL_BIN)/golangci-lint run --config=.golangci.yaml ./...

.PHONY: lint
lint: bin-deps .lint

.PHONY: build
build: build-linux build-darwin

.PHONY: build-linux
build-linux:
	go mod download && CGO_ENABLED=0 \
	GOOS=linux GOARCH=${GOARCH} go build -o ${LOCAL_BIN}/git-update-linux-${GOARCH} ./cmd/main.go;

.PHONY: build-darwin
build-darwin:
	go mod download && CGO_ENABLED=0 \
	GOOS=darwin GOARCH=${GOARCH} go build -o ${LOCAL_BIN}/git-update-darwin-${GOARCH} ./cmd/main.go;

.PHONY: test
test:
	go test ./... -count=1 -timeout=60s -v -short

.PHONY: install
install:
	go build -o $(GOBIN)/protogen -ldflags "$(LDFLAGS)" $(CURDIR)/cmd