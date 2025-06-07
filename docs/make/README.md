# Gutils Build & Distro Creation Guide

This guide explains how to build the Gutils utilities, kernel, root filesystem, and ISO image.

---

## Prerequisites
- Go toolchain (1.22+)
- GNU Make
- Kernel build dependencies (gcc, make, etc.)
- syslinux (isolinux) and xorriso for ISO creation

---

## Fixing ISO Creation in Fedora.

When building in Fedora, you may find that there's not a "isolinux" package. Here's what to do.

```bash
sudo dnf install syslinux
sudo ln -s /usr/share/syslinux /usr/lib/ISOLINUX
```

This should fix ISO creation as of Fedora 42.

## To Install Dependencies using APT (as of Ubuntu 24.04)

```bash
sudo apt update && sudo apt install -y \
  build-essential \
  bc \
  bison \
  flex \
  libssl-dev \
  libelf-dev \
  libncurses-dev \
  grub-efi-arm64-bin \
  grub-mkrescue \
  xorriso \
  mtools \
  qemu-system-aarch64 \
  qemu-efi-aarch64 \
  busybox-static \
  golang \
  make \
  cmake
```

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
