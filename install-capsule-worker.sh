#!/bin/bash
LAST_CAPSULE_WORKER_VERSION="v0.0.0"
echo "System: ${OSTYPE} $(uname -m)"

if [ -z "$CAPSULE_WORKER_PATH" ]
then
    CAPSULE_WORKER_PATH="$HOME/.local/bin"
fi

if [[ $1 = "help" ]]
then
    echo "usage: $0"
    echo "The script will detect the OS & ARCH and use the last version of capsule (${LAST_CAPSULE_WORKER_VERSION})"
    echo "You can force the values by setting these environment variables:"
    echo "- CAPSULE_WORKER_OS (linux, darwin)"
    echo "- CAPSULE_WORKER_ARCH (amd64, arm64)"
    echo "- CAPSULE_WORKER_VERSION (default: ${LAST_CAPSULE_WORKER_VERSION})"
    echo "- CAPSULE_WORKER_PATH (default: ${CAPSULE_WORKER_PATH})"
    exit 0
fi

if [ -z "$CAPSULE_WORKER_VERSION" ]
then
    CAPSULE_WORKER_VERSION=$LAST_CAPSULE_WORKER_VERSION
fi

if [ -z "$CAPSULE_WORKER_OS" ]
then
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
      CAPSULE_WORKER_OS="linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
      CAPSULE_WORKER_OS="darwin"
    else
      CAPSULE_WORKER_OS="linux"
    fi
fi

if [ -z "$CAPSULE_WORKER_ARCH" ]
then
    if [[ "$(uname -m)" == "x86_64" ]]; then
      CAPSULE_WORKER_ARCH="amd64"
    elif [[ "$(uname -m)" == "arm64" ]]; then
      CAPSULE_WORKER_ARCH="arm64"
    elif [[ $(uname -m) == "aarch64" ]]; then
      CAPSULE_WORKER_ARCH="arm64"
    else
      CAPSULE_WORKER_ARCH="amd64"
    fi
fi

MODULE="capsule-worker"


echo "Installing ${MODULE}..."
wget https://github.com/bots-garden/capsule-worker/releases/download/${CAPSULE_WORKER_VERSION}/${MODULE}-${CAPSULE_WORKER_VERSION}-${CAPSULE_WORKER_OS}-${CAPSULE_WORKER_ARCH}.tar.gz
tar -zxf ${MODULE}-${CAPSULE_WORKER_VERSION}-${CAPSULE_WORKER_OS}-${CAPSULE_WORKER_ARCH}.tar.gz --directory ${CAPSULE_WORKER_PATH}
rm ${MODULE}-${CAPSULE_WORKER_VERSION}-${CAPSULE_WORKER_OS}-${CAPSULE_WORKER_ARCH}.tar.gz

echo "Capsule[module: ${MODULE}] $(capsule version) is installed 🎉"
