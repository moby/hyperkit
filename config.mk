###############################################################################
# Config                                                                      #
#                                                                             #
# [XHYVE_CONFIG_ASSERT] VMM asserts (disable for release builds?)             #
# [XHYVE_CONFIG_TRACE]  VMM event tracer                                      #
# [XHYVE_CONFIG_STATS]  VMM event profiler                                    #
###############################################################################

DEFINES := \
  -DXHYVE_CONFIG_ASSERT

###############################################################################
# Toolchain                                                                   #
###############################################################################

CC := clang
AS := clang
LD := clang
STRIP := strip
DSYM := dsymutil
DTRACE := dtrace
CODESIGN := codesign

ENV := \
  LANG=en_US.US-ASCII

###############################################################################
# CFLAGS                                                                      #
###############################################################################

CFLAGS_OPT := \
  -Os \
  -flto \
  -fstrict-aliasing

# enable everything and then selectively disable some warnings
CFLAGS_WARN := \
  -Weverything \
  -Werror \
  -pedantic \
  \
  -Wno-dollar-in-identifier-extension \
  -Wno-gnu-statement-expression \
  -Wno-packed \
  -Wno-padded \
  -Wno-reserved-id-macro \
  -Wno-unknown-warning-option \
  -Wno-unused-macros \
  -Wno-switch-enum \
  -Wno-poison-system-directories \
  -Wno-error=unused-but-set-variable

CFLAGS_DIAG := \
  -fmessage-length=152 \
  -fdiagnostics-show-note-include-stack \
  -fmacro-backtrace-limit=0 \
  -fcolor-diagnostics

CFLAGS_DBG := \
  -g

CFLAGS := \
  -arch x86_64 \
  -x c \
  -std=c11 \
  -fno-common \
  -fvisibility=hidden \
  $(DEFINES) \
  $(CFLAGS_OPT) \
  $(CFLAGS_WARN) \
  $(CFLAGS_DIAG) \
  $(CFLAGS_DBG)

###############################################################################
# LDFLAGS                                                                     #
###############################################################################

LDFLAGS_DBG := \
  -Xlinker -object_path_lto

LDFLAGS := \
  -arch x86_64 \
  -framework Hypervisor \
  -framework vmnet \
  -Wl,-no_compact_unwind \
  $(LDFLAGS_DBG)
