#!/bin/bash

set -e
set -x

NAME=MyBasicPlugin
# http://stackoverflow.com/a/1371283/358804
BIN=${PWD##*/}

go build

# reinstall
if cf plugins | grep -q "$NAME"; then
  cf uninstall-plugin "$NAME"
fi
cf install-plugin -f "$BIN"

cf basic-plugin-command
