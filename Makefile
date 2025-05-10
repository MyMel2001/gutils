# Makefile for Sneed Coreutils

# List of utilities
UTILS = \
	alias \
	base64 \
	basename \
	calc \
	bse \
	cat \
	chmod \
	chown \
	cp \
	cut \
	date \
	dd \
	df \
	dirname \
	dosu \
	du \
	echo \
	env \
	expand \
	false \
	find \
	fold \
	grep \
	head \
	highway \
	id \
	kill \
	killall \
	ls \
	lsblk \
	mkdir \
	mv \
	nl \
	no \
	paste \
	printf \
	ps \
	pwd \
	realpath \
	rm \
	seq \
	sha256sum \
	sha512sum \
	shortcut \
	sleep \
	sort \
	split \
	stat \
	sum \
	tac \
	tail \
	tar \
	tee \
	test \
	touch \
	tr \
	true \
	unexpand \
	uniq \
	unlink \
	wc \
	whereis \
	who \
	whoami \
	xargs \
	yes

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

BINARIES = \
	bin/cat \
	bin/echo \
	bin/ls \
	bin/cp \
	bin/mv \
	bin/rm \
	bin/mkdir \
	bin/touch \
	bin/pwd \
	bin/true \
	bin/false \
	bin/head \
	bin/tail \
	bin/wc \
	bin/sleep \
	bin/date \
	bin/env \
	bin/uname \
	bin/df \
	bin/du \
	bin/ps \
	bin/kill \
	bin/killall \
	bin/yes \
	bin/no \
	bin/tee \
	bin/test \
	bin/uniq \
	bin/xargs \
	bin/find \
	bin/tar \
	bin/sha256sum \
	bin/sha512sum \
	bin/whereis \
	bin/realpath \
	bin/shortcut \
	bin/chmod \
	bin/chown \
	bin/dd \
	bin/base64 \
	bin/bse \
	bin/highway \
	bin/calc \
	bin/basename \
	bin/cut \
	bin/expand \
	bin/fold \
	bin/id \
	bin/nl \
	bin/paste \
	bin/seq \
	bin/split \
	bin/stat \
	bin/sum \
	bin/tac \
	bin/tr \
	bin/unexpand \
	bin/unlink \
	bin/who

bin/%: cmd/%/main.go
	@mkdir -p bin
	go build -o $@ ./cmd/$*/ 