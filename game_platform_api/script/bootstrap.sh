#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=game_platform_api
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}