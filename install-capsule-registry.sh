#!/bin/bash
LAST_CAPSULE_REGISTRY_VERSION="v0.0.0"
echo "System: ${OSTYPE} $(uname -m)"

if [ -z "$CAPSULE_REGISTRY_PATH" ]
then
    CAPSULE_REGISTRY_PATH="$HOME/.local/bin"
fi

if [[ $1 = "help" ]]
then
    echo "usage: $0"
    echo "The script will detect the OS & ARCH and use the last version of capsule (${LAST_CAPSULE_REGISTRY_VERSION})"
    echo "You can force the values by setting these environment variables:"
    echo "- CAPSULE_REGISTRY_OS (linux, darwin)"
    echo "- CAPSULE_REGISTRY_ARCH (amd64, arm64)"
    echo "- CAPSULE_REGISTRY_VERSION (default: ${LAST_CAPSULE_REGISTRY_VERSION})"
    echo "- CAPSULE_REGISTRY_PATH (default: ${CAPSULE_REGISTRY_PATH})"
    exit 0
fi

if [ -z "$CAPSULE_REGISTRY_VERSION" ]
then
    CAPSULE_REGISTRY_VERSION=$LAST_CAPSULE_REGISTRY_VERSION
fi

if [ -z "$CAPSULE_REGISTRY_OS" ]
then
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
      CAPSULE_REGISTRY_OS="linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
      CAPSULE_REGISTRY_OS="darwin"
    else
      CAPSULE_REGISTRY_OS="linux"
    fi
fi

if [ -z "$CAPSULE_REGISTRY_ARCH" ]
then
    if [[ "$(uname -m)" == "x86_64" ]]; then
      CAPSULE_REGISTRY_ARCH="amd64"
    elif [[ "$(uname -m)" == "arm64" ]]; then
      CAPSULE_REGISTRY_ARCH="arm64"
    elif [[ $(uname -m) == "aarch64" ]]; then
      CAPSULE_REGISTRY_ARCH="arm64"
    else
      CAPSULE_REGISTRY_ARCH="amd64"
    fi
fi

MODULE="capsule-registry"


echo "Installing ${MODULE}..."
wget https://github.com/bots-garden/capsule-registry/releases/download/${CAPSULE_REGISTRY_VERSION}/${MODULE}-${CAPSULE_REGISTRY_VERSION}-${CAPSULE_REGISTRY_OS}-${CAPSULE_REGISTRY_ARCH}.tar.gz
tar -zxf ${MODULE}-${CAPSULE_REGISTRY_VERSION}-${CAPSULE_REGISTRY_OS}-${CAPSULE_REGISTRY_ARCH}.tar.gz --directory ${CAPSULE_REGISTRY_PATH}
rm ${MODULE}-${CAPSULE_REGISTRY_VERSION}-${CAPSULE_REGISTRY_OS}-${CAPSULE_REGISTRY_ARCH}.tar.gz

echo "Capsule[module: ${MODULE}] $(capsule version) is installed 🎉"
