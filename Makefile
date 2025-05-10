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
	hserve

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

bin/%: cmd/%/main.go
	@mkdir -p bin
	go build -o $@ ./cmd/$*/ 