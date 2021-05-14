#!/bin/bash
HOOK_DIR=$(git rev-parse --show-toplevel)/scripts/hooks
if [ -x $0.local ]; then
    $0.local "$@" || exit $?
fi
if [ -x $HOOK_DIR/$(basename $0) ]; then
    $HOOK_DIR/$(basename $0) "$@" || exit $?
fi
