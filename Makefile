include build.mk

export GO111MODULE=on
# export GOSUMDB=off

BUILD_ENV_PARAMS:=CGO_ENABLED=0

space := $(subst ,, )
CURDIR_ESCAPE:=$(subst $(space),\ ,$(CURDIR))

LOCAL_BIN:=$(CURDIR_ESCAPE)/bin
LINT_BIN:=$(LOCAL_BIN)/golangci-lint
LINT_VERSION:=2.12.2
INSTALLED_LINT_VERSION:=$(shell if [ -x "$(LINT_BIN)" ]; then "$(LINT_BIN)" version 2>/dev/null | sed -E 's/.* version ([^ ]+) .*/\1/'; fi)

###### TEST ######
.PHONY: test
test:
	go test ./... -count=1 -timeout=60s -v -short
###### TEST ######

###### LINT ######
.PHONY: install-lint
install-lint:
ifneq ("$(INSTALLED_LINT_VERSION)","$(LINT_VERSION)")
	$(info Installing golangci-lint v$(LINT_VERSION))
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$(LINT_VERSION)
# Устанавливаем текущий путь для исполняемого файла линтера.
else
	$(info Golangci-lint v$(LINT_VERSION) is already installed)
endif

.PHONY: lint
lint: install-lint
	$(info Running lint against changed files...)
	$(LINT_BIN) run \
		--new-from-rev=origin/master \
		--config=.golangci.yml \
		./...

.PHONY: lint-full
lint-full: install-lint
	$(info Running lint against all project files...)
	$(LINT_BIN) run \
		--config=.golangci.yml \
		./...
###### LINT ######
