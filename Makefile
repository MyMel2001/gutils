# Makefile for Gutils

GOARCH ?= arm64
GOOS ?= linux

# List of utilities
UTILS = \
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
	uname \
	unexpand \
	uniq \
	unlink \
	wc \
	whereis \
	who \
	whoami \
	xargs \
	yes \
	mount \
	ping \
	ip \
	wifi-connect \
	ethernet-connect \
	dhcp-get \
	porridge \
	interwebz \
	lnkr \
	g \
	qmachine \
	swirl \
	hserve \
	zip \
	unzip \
	install-distro \
	expand-fs \
	susie \
	mkusr \
	shutdown \
	rmdir \
	sed \
	stty \
	mknod \
	od \
	pr \
	expr \
	chgrp \
	cmp \
	bruv

# Default target: build all utilities
all: $(UTILS)

# Build each utility
$(UTILS):
	@mkdir -p bin
	go build -o bin/$@ ./cmd/$@

# Clean build outputs
clean:
	rm -rf bin

# Build a ready-to-use Linux image with all utilities and custom rootfs
# Usage: make distro DOSU_PASS=yourpassword
.PHONY: distro

distro:
	@if [ -z "$(DOSU_PASS)" ]; then \
		echo "ERROR: DOSU_PASS must be set (e.g. make distro DOSU_PASS=yourpassword)"; \
		exit 1; \
	fi
	$(MAKE) all
	$(MAKE) -f Makefile.kernel all DOSU_PASS="$(DOSU_PASS)"

# Build a ready-to-use Linux image with all utilities and custom rootfs
# Usage: make quickreimage DOSU_PASS=yourpassword	
.PHONY: quickreimage

quickreimage:
	@if [ -z "$(DOSU_PASS)" ]; then \
		echo "ERROR: DOSU_PASS must be set (e.g. make quickreimage DOSU_PASS=yourpassword)"; \
		exit 1; \
	fi
	$(MAKE) -f Makefile.kernel quickreimage DOSU_PASS="$(DOSU_PASS)"

# Alias for compatibility with Makefile.kernel
utils: all

.PHONY: all clean $(UTILS)

bin/%: cmd/%/main.go
	@mkdir -p bin
	go build -o $@ ./cmd/$*/ 
