#!/bin/sh
# Enhanced cd - automatically lists directory after changing

cd() {
    builtin cd "$@" && ll
}
