#!/bin/bash

# 1. Install the correct Golang-SDK for your distro and architecture.
# 2. Fetch deps for the project
# 3. Build


SYSTEM_ARCHITECTURE="armv6l"
GOLANG_VERSION="1.20.7"

GOLANG_SRC_URL="https://go.dev/dl/go$GOLANG_VERSION.linux-$SYSTEM_ARCHITECTURE.tar.gz"
GOLANG_SRC_FILE="go$GOLANG_VERSION.linux-$SYSTEM_ARCHITECTURE.tar.gz"

if ! -e $GOLANG_SRC_FILE ; then wget $GOLANG_SRC_URL; fi

tar -xzf $GOLANG_SRC_FILE -C ../


../go/bin/go mod tidy


../go/bin/go build

