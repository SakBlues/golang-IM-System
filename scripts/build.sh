#!/usr/bin/env bash

TARGET_DIR="bin"

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "${SOURCE}" ] ; do SOURCE="$(readlink "${SOURCE}")"; done
BASE_DIR="$( cd -P "$( dirname "${SOURCE}" )/.." && pwd )"

# Delete the old dir, and make new dir.
echo "==> Removing $TARGET_DIR/ ..."
rm -rf ${BASE_DIR}/$TARGET_DIR/*
echo "==> Creating new dir $TARGET_DIR/"
mkdir -p ${BASE_DIR}/$TARGET_DIR/

echo "Building server..."
go build -o ${BASE_DIR}/$TARGET_DIR/server ${BASE_DIR}/cmd/server/main.go

echo "Building client..."
go build -o ${BASE_DIR}/$TARGET_DIR/client ${BASE_DIR}/cmd/client/main.go

 
