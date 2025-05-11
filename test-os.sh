#!/bin/bash

# Script to create and run a QEMU VM with:
# - 2GB RAM
# - 16GB disk storage
# - 2 CPU threads
# - Using gutils-linux.iso

# Configuration variables
VM_NAME="gutils-vm"
RAM_SIZE="2G"
DISK_SIZE="16G"
CPU_THREADS=2
ISO_FILE="gutils-linux.iso"
DISK_IMG="${VM_NAME}.qcow2"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if the ISO file exists
if [ ! -f "$ISO_FILE" ]; then
    echo -e "${YELLOW}Warning: ISO file '${ISO_FILE}' not found in the current directory.${NC}"
    read -p "Do you want to continue anyway? (y/n): " CONTINUE
    if [[ ! "$CONTINUE" =~ ^[Yy]$ ]]; then
        echo "Exiting."
        exit 1
    fi
fi

# Check if QEMU is installed
if ! command -v qemu-system-x86_64 &> /dev/null; then
    echo "QEMU is not installed. Installing QEMU..."
    
    # Check the package manager and install QEMU
    if command -v dnf &> /dev/null; then
        sudo dnf install -y qemu-kvm
    elif command -v apt-get &> /dev/null; then
        sudo apt-get update
        sudo apt-get install -y qemu-kvm qemu-system-x86
    elif command -v pacman &> /dev/null; then
        sudo pacman -S qemu
    else
        echo "Could not detect package manager. Please install QEMU manually."
        exit 1
    fi
fi

# Check if disk image already exists
if [ -f "$DISK_IMG" ]; then
    echo -e "${YELLOW}Disk image '${DISK_IMG}' already exists.${NC}"
    read -p "Do you want to create a new one? This will delete the existing one. (y/n): " CREATE_NEW
    if [[ "$CREATE_NEW" =~ ^[Yy]$ ]]; then
        rm "$DISK_IMG"
    else
        echo "Using existing disk image."
    fi
fi

# Create disk image if it doesn't exist
if [ ! -f "$DISK_IMG" ]; then
    echo "Creating disk image '${DISK_IMG}' with size ${DISK_SIZE}..."
    qemu-img create -f qcow2 "$DISK_IMG" "$DISK_SIZE"
    if [ $? -ne 0 ]; then
        echo "Failed to create disk image. Exiting."
        exit 1
    fi
    echo -e "${GREEN}Disk image created successfully.${NC}"
fi

# Run the VM
echo -e "${GREEN}Starting QEMU VM with:${NC}"
echo "  - RAM: $RAM_SIZE"
echo "  - Disk: $DISK_IMG ($DISK_SIZE)"
echo "  - CPU Threads: $CPU_THREADS"
echo "  - ISO: $ISO_FILE"
echo

qemu-system-x86_64 \
    -name "$VM_NAME" \
    -m "$RAM_SIZE" \
    -smp "$CPU_THREADS" \
    -boot d \
    -drive file="$DISK_IMG",format=qcow2 \
    -cdrom "$ISO_FILE" \
    -enable-kvm \
    -machine type=q35,accel=kvm \
    -device virtio-net-pci,netdev=net0 \
    -netdev user,id=net0 \
    -display gtk

# Check if QEMU exited successfully
if [ $? -eq 0 ]; then
    echo -e "${GREEN}VM session ended normally.${NC}"
else
    echo -e "${YELLOW}VM session ended with an error code: $?${NC}"
fi

exit 0
