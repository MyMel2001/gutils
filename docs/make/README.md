# Gutils Build & Distro Creation Guide

This guide explains how to build the Gutils utilities, kernel, root filesystem, and ISO image.

---

## Prerequisites
- Go toolchain (1.22+)
- GNU Make
- Kernel build dependencies (gcc, make, etc.)
- syslinux (isolinux) and xorriso for ISO creation

---

## To Install Dependencies using APT (as of Ubuntu 24.04/Fedora 42)

```bash
sudo apt update && sudo apt install -y \
  build-essential \
  bc \
  bison \
  flex \
  libssl-dev \
  libelf-dev \
  libncurses-dev \
  grub-efi-amd64-bin \
  xorriso \
  mtools \
  busybox-static \
  golang \
  make \
  cmake gcc-12
```
then
```
sudo update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-12 60
```

In Fedora:

```
sudo dnf install -y podman distrobox
distrobox create --name gutils-devel --image ubuntu:24.04
```
then
```
distrobox enter gutils-devel
```
and follow Ubuntu instructions.

## Building Utilities
To build all utilities:
```
make
```
Binaries will be placed in the `bin/` directory.

---

## Building the Kernel and Root Filesystem
To build the kernel, rootfs, and initramfs:
```
make -f Makefile.kernel all DOSU_PASS=yourpassword
```
- Replace `yourpassword` with the password you want for `dosu`.
- The rootfs will use `/bin/highway` as the default shell.

---

## Creating a Bootable ISO
To create a bootable ISO image:
```
make -f Makefile.kernel iso DOSU_PASS=yourpassword
```
The ISO will be named `gutils-linux.iso`.

---

## One-Step Distro Build
To build everything (utilities, kernel, rootfs, ISO) in one step:
```
make distro DOSU_PASS=yourpassword
```

---

## Cleaning Build Artifacts
To clean all build outputs:
```
make clean
make -f Makefile.kernel clean
```

---

## Notes
- The kernel version and build options can be changed in `Makefile.kernel`.
- All utilities in `bin/` are included in the rootfs and ISO automatically. 
