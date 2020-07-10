GOLINT ?= golangci-lint

all: build lint tidy
.PHONY: all

SUBMODULES=

FFI_PATH:=./extern/filecoin-ffi/
FFI_DEPS:=.install-filcrypto
FFI_DEPS:=$(addprefix $(FFI_PATH),$(FFI_DEPS))

$(FFI_DEPS): .filecoin-build ;

.filecoin-build: $(FFI_PATH)
	$(MAKE) -C $(FFI_PATH) $(FFI_DEPS:$(FFI_PATH)%=%)
	@touch $@

.update-modules:
	git submodule update --init --recursive
	@touch $@

chainval: .update-modules .filecoin-build
	go build ./...
.PHONY: chainval
SUBMODULES+=chainval

build: $(SUBMODULES)

clean:
	rm -f .filecoin-build
	rm -f .update-modules

resources:
	go generate ./box/...

tidy:
	go mod tidy
.PHONY: tidy

lint:
	$(GOLINT) run ./...
.PHONY: lint
