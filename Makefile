# Makefile for Sneed Coreutils

# List of utilities
UTILS = highway ls cat echo cp mv rm printf mkdir grep head alias tail wc dosu

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