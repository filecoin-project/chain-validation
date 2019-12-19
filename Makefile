all: build
.PHONY: all

SUBMODULES=

FFI_PATH:=./extern/filecoin-ffi/
FFI_DEPS:=libfilecoin.a filecoin.pc filecoin.h
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

gen:
	go run ./pkg/gen/main.go
