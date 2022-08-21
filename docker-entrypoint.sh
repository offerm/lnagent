#!/bin/sh
set -e

# this script starts lnagent

PROGNAME=$(basename $0)

STARTUP="/lnagent $@"

echo "$PROGNAME: Starting $STARTUP"
exec ${STARTUP}
