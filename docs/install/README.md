# Gutils System Installation Guide

This guide describes how to install the Gutils-based Linux system onto a new drive.

---

## Prerequisites
- A built Gutils ISO or image
- A target drive (e.g., /dev/sdb)
- A running system with Gutils utilities available

---

## Installation Steps

### 1. Boot into the Gutils Live System
Boot from the Gutils ISO or image using your preferred method (USB, QEMU, etc).

### 2. List Available Drives
Run:
```
lsblk
```
Identify the target drive for installation (e.g., /dev/sdb).

### 3. Run the Installer
Run:
```
install-distro
```
Follow the prompts to select the target device. The installer will copy the running system and expand the root filesystem.

### 4. First Boot
After installation, reboot the system and boot from the target drive. The default shell will be `highway`.

---

## Troubleshooting
- Ensure the target drive is not mounted or in use.
- Data on the target drive will be overwritten.
- For advanced partitioning, use `fdisk` or `parted` before running `install-distro`. 