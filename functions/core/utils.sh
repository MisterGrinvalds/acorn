#!/bin/sh
# Core shell utilities

# Run a bash shell as another user
bash-as() {
    sudo -u "$1" /bin/bash
}

# Remove all environment variables matching pattern
rmenv() {
    unset $(env | grep -i "${1:-prox}" | grep -oE '^[^=]+')
}
