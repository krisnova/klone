#!/usr/bin/env bash
# Copyright © 2017 Kris Nova <kris@nivenly.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#  _  ___
# | |/ / | ___  _ __   ___
# | ' /| |/ _ \| '_ \ / _ \
# | . \| | (_) | | | |  __/
# |_|\_\_|\___/|_| |_|\___|
#
# BOOTSTRAP.sh will bootstrap and run klone on any system
# Usage: BOOTSTRAP.sh <query> <bash:command>

if [ -z "$1" ]; then
    echo "Usage: BOOTSTRAP.sh <query> <command>"
    exit 1
fi
QUERY=${1}

CMD=""
if [ -n "$2" ]; then
    CMD=${2}
fi

exists() {
  command -v "$1" >/dev/null 2>&1
}

HACK="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
VERSION=$(cat "${HACK}/../pkg/local/version.go" | grep Version | cut -d '"' -f 2)


DOWNLOAD_URL="https://github.com/kris-nova/klone/releases/download/v1.1.1/linux-amd64"
INSTALL_DIR="/usr/local/bin"
BIN_NAME="darwin-amd64"

if [ "$(uname)" == "Darwin" ]; then
    echo "Darwin"
    BIN_NAME="darwin-amd64"
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
   echo "Linux"
   BIN_NAME="linux-amd64"
elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW64_NT" ]; then
    echo "Win64"
    https://github.com/kris-nova/klone/releases/download/v${VERSION}/windows-amd64
    BIN_NAME="windows-amd64"
fi
DOWNLOAD_URL="https://github.com/kris-nova/klone/releases/download/v${VERSION}/${BIN_NAME}"

if exists klone; then
    echo "Exists"
else
    echo "Downloading klone"
    # assume wget
    wget $DOWNLOAD_URL
    chmod +x $BIN_NAME
    mv $BIN_NAME $INSTALL_DIR/klone
    PATH=$PATH:$INSTALL_DIR
fi

klone ${QUERY}
if [ -n "$CMD" ]; then
    eval ${CMD}
fi


