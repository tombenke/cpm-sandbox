APP_NAME := cpm
BINARY := bin/$(APP_NAME)
GO := $(GOROOT)/bin/go
GO_BUILD_FLAGS := -ldflags="-s -w"
GO_TEST_FLAGS := -v
GO_LINT_FLAGS :=
GO_FMT_FLAGS :=
GO_VET_FLAGS :=

.SUFFIXES: .yml .png .uml
EXAMPLES_DIR := examples
EXAMPLES_YML_FILES := $(wildcard $(EXAMPLES_DIR)/*.yml)
UML_FILES := $(EXAMPLES_YML_FILES:.yml=.uml)
PNG_FILES := $(EXAMPLES_YML_FILES:.yml=.png)

all: deps test fmt vet build $(UML_FILES) $(PNG_FILES)

build:
	$(GO) build $(GO_BUILD_FLAGS) -o $(BINARY) cmd/main.go

.yml.uml:
	$(BINARY) $<

.uml.png:
	java -jar ./scripts/plantuml-1.2024.7.jar -Tpng $<

clean:
	rm -rf $(BINARY)
	rm -f $(EXAMPLES_DIR)/*.uml
	rm -f $(EXAMPLES_DIR)/*.png
	rm -f $(EXAMPLES_DIR)/*.txt

deps:
	$(GO) mod tidy
	$(GO) mod vendor

fmt:
	$(GO) fmt $(GO_FMT_FLAGS) ./...

lint:
	golangci-lint run $(GO_LINT_FLAGS)

vet:
	$(GO) vet $(GO_VET_FLAGS) ./...

test:
	$(GO) test $(GO_TEST_FLAGS) ./...
