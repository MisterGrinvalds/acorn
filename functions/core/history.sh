#!/bin/sh
# History search helper

# Shorthand for history with grep
h() {
    history | grep "$1"
}
