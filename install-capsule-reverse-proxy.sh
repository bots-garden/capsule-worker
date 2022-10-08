#!/bin/bash
LAST_CAPSULE_REVERSE_PROXY_VERSION="v0.0.0"
echo "System: ${OSTYPE} $(uname -m)"

if [ -z "$CAPSULE_REVERSE_PROXY_PATH" ]
then
    CAPSULE_REVERSE_PROXY_PATH="$HOME/.local/bin"
fi

if [[ $1 = "help" ]]
then
    echo "usage: $0"
    echo "The script will detect the OS & ARCH and use the last version of capsule (${LAST_CAPSULE_REVERSE_PROXY_VERSION})"
    echo "You can force the values by setting these environment variables:"
    echo "- CAPSULE_REVERSE_PROXY_OS (linux, darwin)"
    echo "- CAPSULE_REVERSE_PROXY_ARCH (amd64, arm64)"
    echo "- CAPSULE_REVERSE_PROXY_VERSION (default: ${LAST_CAPSULE_REVERSE_PROXY_VERSION})"
    echo "- CAPSULE_REVERSE_PROXY_PATH (default: ${CAPSULE_REVERSE_PROXY_PATH})"
    exit 0
fi

if [ -z "$CAPSULE_REVERSE_PROXY_VERSION" ]
then
    CAPSULE_REVERSE_PROXY_VERSION=$LAST_CAPSULE_REVERSE_PROXY_VERSION
fi

if [ -z "$CAPSULE_REVERSE_PROXY_OS" ]
then
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
      CAPSULE_REVERSE_PROXY_OS="linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
      CAPSULE_REVERSE_PROXY_OS="darwin"
    else
      CAPSULE_REVERSE_PROXY_OS="linux"
    fi
fi

if [ -z "$CAPSULE_REVERSE_PROXY_ARCH" ]
then
    if [[ "$(uname -m)" == "x86_64" ]]; then
      CAPSULE_REVERSE_PROXY_ARCH="amd64"
    elif [[ "$(uname -m)" == "arm64" ]]; then
      CAPSULE_REVERSE_PROXY_ARCH="arm64"
    elif [[ $(uname -m) == "aarch64" ]]; then
      CAPSULE_REVERSE_PROXY_ARCH="arm64"
    else
      CAPSULE_REVERSE_PROXY_ARCH="amd64"
    fi
fi

MODULE="capsule-reverse-proxy"


echo "Installing ${MODULE}..."
wget https://github.com/bots-garden/capsule-reverse-proxy/releases/download/${CAPSULE_REVERSE_PROXY_VERSION}/${MODULE}-${CAPSULE_REVERSE_PROXY_VERSION}-${CAPSULE_REVERSE_PROXY_OS}-${CAPSULE_REVERSE_PROXY_ARCH}.tar.gz
tar -zxf ${MODULE}-${CAPSULE_REVERSE_PROXY_VERSION}-${CAPSULE_REVERSE_PROXY_OS}-${CAPSULE_REVERSE_PROXY_ARCH}.tar.gz --directory ${CAPSULE_REVERSE_PROXY_PATH}
rm ${MODULE}-${CAPSULE_REVERSE_PROXY_VERSION}-${CAPSULE_REVERSE_PROXY_OS}-${CAPSULE_REVERSE_PROXY_ARCH}.tar.gz

echo "Capsule[module: ${MODULE}] $(capsule version) is installed 🎉"
