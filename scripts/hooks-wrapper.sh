#!/bin/bash
if [ -x $0.local ]; then
    $0.local "$@" || exit $?
fi
if [ -x tracked_hooks/$(basename $0) ]; then
    tracked_hooks/$(basename $0) "$@" || exit $?
fi
