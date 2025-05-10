# Makefile for Sneed Coreutils

# List of utilities
UTILS = highway ls cat echo cp mv rm printf mkdir grep head alias tail wc dosu touch sha256sum chmod chown shortcut sha512sum base64 du pwd false true sleep dirname lsblk whereis realpath dd sort whoami env kill ps uname date df find tar tee test uniq xargs yes no killall bse

# Default target: build all utilities
all: $(UTILS)

# Build each utility
$(UTILS):
	@mkdir -p bin
	go build -o bin/$@ ./cmd/$@

# Clean build outputs
clean:
	rm -rf bin

.PHONY: all clean $(UTILS) 