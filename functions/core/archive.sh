#!/bin/sh
# Archive creation utilities

# Create a .tar.gz archive of a file or folder
mktar() {
    tar cvzf "${1%%/}.tar.gz" "${1%%/}/"
}

# Create a .zip archive of a file or folder
mkzip() {
    zip -r "${1%%/}.zip" "${1%%/}/"
}
