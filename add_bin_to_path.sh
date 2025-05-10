#!/usr/bin/env bash
# Usage: source ./add_bin_to_path.sh
# Temporarily adds ./bin to the beginning of your $PATH for this shell session.
export PATH="$(pwd)/bin:$PATH"
echo "Added $(pwd)/bin to PATH for this session." 