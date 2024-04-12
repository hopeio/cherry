#!/usr/bin/env bash

PROTOC_VERSION=${1:-"26.1"}
DOWNLOAD_PATH=${2:-"."}

ARCH=$(uname -m)
case $ARCH in
  "arm64")
    ARCH="aarch_64"
    ;;
  "arm32")
    ARCH="aarch_32"
    ;;
  "aarch64")
    ARCH="aarch_64"
    ;;
  "aarch32")
    ARCH="aarch_32"
    ;;
esac

OS=$(uname)
case $OS in
  "Darwin")
    OS="osx"
    ;;
  "Linux")
    OS="linux"
    ;;
esac

ARCHIVE_NAME="protoc-$PROTOC_VERSION-$OS-$ARCH"
echo $DOWNLOAD_PATH/$ARCHIVE_NAME
cd $DOWNLOAD_PATH
curl -LO "https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/$ARCHIVE_NAME.zip"

unzip -o $ARCHIVE_NAME -d $ARCHIVE_NAME
export PATH=$DOWNLOAD_PATH/$ARCHIVE_NAME/bin:$PATH

