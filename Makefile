GIT_VERSION := $(shell git describe --abbrev=6 --dirty --always --tags)
GIT_VERSION_SHA1 := $(shell git rev-parse HEAD)

ifeq ($V, 1)
	VERBOSE =
else
	VERBOSE = @
endif

include config.mk

VMM_LIB_SRC := \
	src/lib/vmm/intel/vmcs.c \
	src/lib/vmm/intel/vmx.c \
	src/lib/vmm/intel/vmx_msr.c \
	\
	src/lib/vmm/io/vatpic.c \
	src/lib/vmm/io/vatpit.c \
	src/lib/vmm/io/vhpet.c \
	src/lib/vmm/io/vioapic.c \
	src/lib/vmm/io/vlapic.c \
	src/lib/vmm/io/vpmtmr.c \
	src/lib/vmm/io/vrtc.c \
	\
	src/lib/vmm/vmm.c \
	src/lib/vmm/vmm_api.c \
	src/lib/vmm/vmm_callout.c \
	src/lib/vmm/vmm_host.c \
	src/lib/vmm/vmm_instruction_emul.c \
	src/lib/vmm/vmm_ioport.c \
	src/lib/vmm/vmm_lapic.c \
	src/lib/vmm/vmm_mem.c \
	src/lib/vmm/vmm_stat.c \
	src/lib/vmm/vmm_util.c \
	src/lib/vmm/x86.c

HYPERKIT_LIB_SRC := \
	src/lib/acpitbl.c \
	src/lib/atkbdc.c \
	src/lib/block_if.c \
	src/lib/consport.c \
	src/lib/dbgport.c \
	src/lib/fwctl.c \
	src/lib/inout.c \
	src/lib/ioapic.c \
	src/lib/log.c \
	src/lib/md5c.c \
	src/lib/mem.c \
	src/lib/mevent.c \
	src/lib/mptbl.c \
	src/lib/pci_ahci.c \
	src/lib/pci_emul.c \
	src/lib/pci_hostbridge.c \
	src/lib/pci_irq.c \
	src/lib/pci_lpc.c \
	src/lib/pci_uart.c \
	src/lib/pci_virtio_9p.c \
	src/lib/pci_virtio_block.c \
	src/lib/pci_virtio_net_tap.c \
	src/lib/pci_virtio_net_vmnet.c \
	src/lib/pci_virtio_net_vpnkit.c \
	src/lib/pci_virtio_rnd.c \
	src/lib/pci_virtio_sock.c \
	src/lib/pm.c \
	src/lib/post.c \
	src/lib/rtc.c \
	src/lib/smbiostbl.c \
	src/lib/task_switch.c \
	src/lib/uart_emul.c \
	src/lib/virtio.c \
	src/lib/xmsr.c

FIRMWARE_LIB_SRC := \
	src/lib/firmware/bootrom.c \
	src/lib/firmware/kexec.c \
	src/lib/firmware/fbsd.c \
	src/lib/firmware/multiboot.c

HYPERKIT_SRC := src/hyperkit.c

HAVE_OCAML_QCOW := $(shell if ocamlfind query qcow prometheus-app uri logs logs.fmt mirage-unix >/dev/null 2>/dev/null ; then echo YES ; else echo NO; fi)

ifeq ($(HAVE_OCAML_QCOW),YES)
CFLAGS += -DHAVE_OCAML=1 -DHAVE_OCAML_QCOW=1 -DHAVE_OCAML=1

LIBEV_FILE=/usr/local/lib/libev.a
LIBEV=$(shell if test -e $(LIBEV_FILE) ; then echo $(LIBEV_FILE) ; fi )

# prefix vsock file names if PRI_ADDR_PREFIX
# is defined. (not applied to aliases)
ifneq ($(PRI_ADDR_PREFIX),)
CFLAGS += -DPRI_ADDR_PREFIX=\"$(PRI_ADDR_PREFIX)\"
endif

# override default connect socket name if 
# CONNECT_SOCKET_NAME is defined 
ifneq ($(CONNECT_SOCKET_NAME),)
CFLAGS += -DCONNECT_SOCKET_NAME=\"$(CONNECT_SOCKET_NAME)\"
endif

OCAML_SRC := \
	src/lib/mirage_block_ocaml.ml

OCAML_C_SRC := \
	src/lib/mirage_block_c.c

OCAML_WHERE := $(shell ocamlc -where)
OCAML_PACKS := cstruct cstruct.lwt io-page io-page.unix uri mirage-block \
	mirage-block-unix qcow unix threads lwt lwt.unix logs logs.fmt   \
	mirage-unix prometheus-app conduit-lwt cohttp.lwt
OCAML_LDLIBS := -L $(OCAML_WHERE) \
	$(shell ocamlfind query cstruct)/cstruct.a \
	$(shell ocamlfind query cstruct)/libcstruct_stubs.a \
	$(shell ocamlfind query io-page)/io_page.a \
	$(shell ocamlfind query io-page-unix)/io_page_unix.a \
	$(shell ocamlfind query io-page-unix)/libio_page_unix_stubs.a \
	$(shell ocamlfind query lwt.unix)/liblwt_unix_stubs.a \
	$(shell ocamlfind query lwt.unix)/lwt_unix.a \
	$(shell ocamlfind query lwt.unix)/lwt.a \
	$(shell ocamlfind query threads)/libthreadsnat.a \
	$(shell ocamlfind query mirage-block-unix)/libmirage_block_unix_stubs.a \
	$(shell ocamlfind query base)/libbase_stubs.a \
        $(LIBEV) \
	-lasmrun -lbigarray -lunix

build/hyperkit.o: CFLAGS += -I$(OCAML_WHERE)
endif

SRC := \
	$(VMM_LIB_SRC) \
	$(HYPERKIT_LIB_SRC) \
	$(FIRMWARE_LIB_SRC) \
	$(OCAML_C_SRC) \
	$(HYPERKIT_SRC)

OBJ := $(SRC:src/%.c=build/%.o) $(OCAML_SRC:src/%.ml=build/%.o)
DEP := $(OBJ:%.o=%.d)
INC := -Isrc/include

CFLAGS += -DVERSION=\"$(GIT_VERSION)\" -DVERSION_SHA1=\"$(GIT_VERSION_SHA1)\"

TARGET = build/hyperkit

all: $(TARGET) | build

.PHONY: clean all test test-qcow
.SUFFIXES:

-include $(DEP)

build:
	@mkdir -p build

src/include/xhyve/dtrace.h: src/lib/dtrace.d
	@echo gen $<
	$(VERBOSE) $(DTRACE) -h -s $< -o $@

$(SRC): src/include/xhyve/dtrace.h

build/%.o: src/%.c
	@echo cc $<
	@mkdir -p $(dir $@)
	$(VERBOSE) $(ENV) $(CC) $(CFLAGS) $(INC) $(DEF) -MMD -MT $@ -MF build/$*.d -o $@ -c $<

$(OCAML_C_SRC:src/%.c=build/%.o): CFLAGS += -I$(OCAML_WHERE)
build/%.o: src/%.ml
	@echo ml $<
	@mkdir -p $(dir $@)
	$(VERBOSE) $(ENV) ocamlfind ocamlopt -thread -package "$(OCAML_PACKS)" -c $< -o build/$*.cmx
	$(VERBOSE) $(ENV) ocamlfind ocamlopt -thread -linkpkg -package "$(OCAML_PACKS)" -output-obj -o $@ build/$*.cmx

$(TARGET).sym: $(OBJ)
	@echo ld $(notdir $@)
	$(VERBOSE) $(ENV) $(LD) $(LDFLAGS) -Xlinker $(TARGET).lto.o -o $@ $(OBJ) $(LDLIBS) $(OCAML_LDLIBS)
	@echo dsym $(notdir $(TARGET).dSYM)
	$(VERBOSE) $(ENV) $(DSYM) $@ -o $(TARGET).dSYM

$(TARGET): $(TARGET).sym
	@echo strip $(notdir $@)
	$(VERBOSE) $(ENV) $(STRIP) $(TARGET).sym -o $@

clean:
	@rm -rf build
	@rm -f src/include/xhyve/dtrace.h
	@rm -f test/vmlinuz test/initrd.gz
	@rm -f test/disk.qcow2

test/vmlinuz test/initrd.gz:
	@cd test; ./tinycore.sh

test: $(TARGET) test/vmlinuz test/initrd.gz
	@(cd test && ./test_linux.exp)

test-qcow: $(TARGET) test/vmlinuz test/initrd.gz
	@(cd test && ./test_linux_qcow.exp)


## ----------- ##
## Artifacts.  ##
## ----------- ##

.PHONY: artifacts
artifacts: build/LICENSE build/COMMIT

.PHONY: build/LICENSE
build/LICENSE:
	@echo "  GEN     " $@
	@find src -type f | xargs awk '/^\/\*-/{p=1;print FILENAME ":";print;next} p&&/^.*\*\//{print;print "";p=0};p' > $@.tmp
	@opam config exec -- make -C repo list-licenses
	@cat repo/OCAML-LICENSES >> $@.tmp
	@mv $@.tmp $@

.PHONY: build/COMMIT
build/COMMIT:
	@echo "  GEN     " $@
	@git rev-parse HEAD > $@.tmp
	@mv $@.tmp $@
