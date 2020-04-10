#!/bin/sh

#########################################################################
# A simple script which can be run on Cloud VMs to download and install #
# the latest und and undcli binaries. This script is BETA               #
#########################################################################

#####################################################################################################################
# Copyright 2020 Codegnosis paul@codegnosis.co.uk
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
# documentation files (the "Software"), to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and
# to permit persons to whom the Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all copies or substantial portions
# of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
# THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
# TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
# IN THE SOFTWARE.
#####################################################################################################################

LATEST_RELEASE="https://api.github.com/repos/unification-com/mainchain/releases/latest"
GH_DL_PREFIX="https://github.com/unification-com/mainchain/releases/download"
LATEST_VERSION="$(curl --silent ${LATEST_RELEASE} | grep -Po '"tag_name": "\K.*?(?=")')"
DOWNLOAD_DEST="/tmp/mainchain_tmp"
LOCAL_BIN="/usr/local/bin"
BACKUP_DIR="${HOME}/UND_BIN_BAK"

BINARIES_TO_INSTALL="und undcli"

INSTALLATION_SUMMARY=""

SKIP_INSTALL=false

bold=$(tput bold)
normal=$(tput sgr0)

echo "Latest release is v${LATEST_VERSION}"

# check binary
# check und
check () {
  CHECK_LOC="${LOCAL_BIN}/${1}"
  GO_BIN=""
  VERS=""
  echo "Check ${1}"

  echo "Checking for manual installation in GOPATH"

  if [ -n "$GOPATH" ]; then
    echo "Found GOPATH: $GOPATH"
    GO_BIN="${GOPATH}/bin/${1}"
    echo "Check for ${GO_BIN}"
    if [ -f "$GO_BIN" ]; then
      VERS=$(${GO_BIN} version)
      MSG="Found ${1} v${VERS} in GOPATH: ${GOPATH} (${GO_BIN}). Please manually uninstall this first."
      INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${MSG}"
      echo "WARNING: ${MSG}"
      SKIP_INSTALL=true
    else
      echo "${1} not found in GOPATH. Continue installation."
    fi
  fi

  echo "Checking for ${CHECK_LOC}"
  if [ ! -f "$CHECK_LOC" ]; then
    echo "${CHECK_LOC} is not currently installed."
  else
    echo "Found ${CHECK_LOC}. Checking installed version"
    VERS=$(${CHECK_LOC} version)
    echo "${1} version ${VERS} currently installed."
    if test "${LATEST_VERSION}" = "${VERS}"; then
      MSG="latest version ${1} v${LATEST_VERSION} already installed in ${CHECK_LOC}"
      INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${MSG}. Nothing to do."
      echo "${bold}NOTICE:${normal} ${MSG}. Skipping."
      SKIP_INSTALL=true
    else
      BACKUP_LOC="${BACKUP_DIR}"
      if [ ! -d "$BACKUP_LOC" ]; then
        mkdir -p "${BACKUP_LOC}"
      fi
      BACKUP_BIN="${BACKUP_LOC}/${1}-v${VERS}.tar.gz"
      echo "Backing up ${1} v${VERS} to ${BACKUP_BIN}"
      tar -czf "${BACKUP_BIN}" -C "${LOCAL_BIN}" "${1}"

      if [ -f "$BACKUP_BIN" ]; then
        MSG="${1} v${VERS} successfully backed up to: ${BACKUP_BIN}"
        echo "${MSG}"
        INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${MSG}\n"
      else
        echo "\n${bold}Failed to backup ${1} v${VERS}. Continue installation? (y to continue)${normal}"
        read -r CONTINUE_INSTALLATION
        if test "${CONTINUE_INSTALLATION}" = "y"; then
          echo "Continuing installation without backup..."
          INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${bold}WARNING:${normal} ${1} v${VERS} was NOT backed up.\n"
        else
          echo "Abort ${1} installation."
          INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\nInstallation of ${1} v${VERS} was aborted by user.\n"
          SKIP_INSTALL=true
        fi
      fi
    fi
  fi
}

# download binary_name version os
# download "undcli" "1.1.1" "linux"
download () {
  DOWLOAD_FILE="${GH_DL_PREFIX}/${2}/${1}_v${2}_${3}_x86_64.tar.gz"
  DEST="${DOWNLOAD_DEST}/${1}.tar.gz"
  echo "Downloading ${DOWLOAD_FILE} to ${DEST}"
  curl -L "${DOWLOAD_FILE}" -o "${DEST}"
  if [ -f "$DEST" ]; then
    echo "download ${DEST} successful"
  else
    echo "${bold}ERROR: failed to download ${DEST}. Exiting${normal}"
    exit
  fi
}

# extract binary_name dir
# extract "undcli"
extract () {
  ARCHIVE="${DOWNLOAD_DEST}/${1}.tar.gz"
  BIN_RES="${DOWNLOAD_DEST}/${1}"
  echo "Extracting ${ARCHIVE}"
  tar -C "${DOWNLOAD_DEST}/" -xzf "${ARCHIVE}"
  if [ -f "$BIN_RES" ]; then
    echo "${BIN_RES} extracted successfully"
  else
    echo "${bold}ERROR: extract failed.${normal}"
    exit
  fi
}

clean() {
  echo "cleaning downloads"
  echo ""
  if [ -d "$DOWNLOAD_DEST" ]; then
    echo "clean ${DOWNLOAD_DEST}"
    rm -rf "${DOWNLOAD_DEST}"
  fi
}

# install binary
# install und
install () {
  if [ -f "${DOWNLOAD_DEST}/${1}" ]; then
    echo "Found ${DOWNLOAD_DEST}/${1}. installing into ${LOCAL_BIN}/${1}"
    sudo cp "${DOWNLOAD_DEST}/${1}" "${LOCAL_BIN}/${1}"
  else
    echo "${bold}ERROR: could not find ${DOWNLOAD_DEST}/${1}. Exiting${normal}"
    exit
  fi

  if [ -f "${LOCAL_BIN}/${1}" ]; then
    echo "Found ${LOCAL_BIN}/${1} make it executable"
    sudo chmod 755 "${LOCAL_BIN}/${1}"
    if [ -x "${LOCAL_BIN}/${1}" ]; then
      echo "Done."
    else
      echo "${bold}ERROR: failed to make ${LOCAL_BIN}/${1} executable. Exiting${normal}"
      exit
    fi
  else
    echo "${bold}ERROR: could not find ${LOCAL_BIN}/${1}. Exiting${normal}"
    exit
  fi
}

# delete any previous downloads
clean

## create a tmp download dir
echo "creating ${DOWNLOAD_DEST}"
mkdir -p "${DOWNLOAD_DEST}"
echo ""

for i in ${BINARIES_TO_INSTALL}; do
  SKIP_INSTALL=false
  INSTALL_LOCATION=""
  INSTALLED_VERSION=""

  echo ""
  echo "${bold}--Attempting to install/update ${i}--${normal}"

  check "${i}"
  echo ""
  if [ "$SKIP_INSTALL" = true ] ; then
    continue
  fi

  if [ -f "${LOCAL_BIN}/${i}" ]; then
    echo "Removing old version of ${LOCAL_BIN}/${i}"
    sudo rm "${LOCAL_BIN}/${i}"
  fi

  echo "Installing ${i}"
  echo ""
  # Download the latest archive
  download "${i}" "${LATEST_VERSION}" "linux"
  echo ""
  # Extract the archive
  extract "${i}"
  echo ""
  # install the binary
  install "${i}"
  echo ""

  echo "${i} has been installed to:"
  INSTALL_LOCATION=$(which "${i}")
  echo "${INSTALL_LOCATION}"
  INSTALLED_VERSION=$("${INSTALL_LOCATION}" version)

  echo "Checking "${INSTALL_LOCATION}" version."
  echo "expected version: ${LATEST_VERSION}"
  echo "installed version: ${INSTALLED_VERSION}"

  if test "${LATEST_VERSION}" = "${INSTALLED_VERSION}"; then
    INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${bold}Latest ${i} v${INSTALLED_VERSION} successfully installed!${normal}\nLocation: ${INSTALL_LOCATION}"
    INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\nRun:\n\n${INSTALL_LOCATION} version --long\n\nto verify. Version should be ${LATEST_VERSION}.\n"
  else
    INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${bold}Something went wrong installing ${i} v${LATEST_VERSION}. See previous output.${normal}"
  fi
done

if [ -z "$HOME" ]; then
  echo "${bold}WARNING: NO HOME ENVIRONMENT VARIABLE DETECTED! ABORTING.${normal}"
  exit
fi
# remove the tmp files
clean

echo ""
echo "----------------------"
echo ""

# output installation summary
if [ -n "$INSTALLATION_SUMMARY" ]; then
  echo "${bold}Installation/Upgrade summary${normal}"
  echo "----------------------------"
  echo "${INSTALLATION_SUMMARY}"
fi

INSTALLED_UND_VERSION=$("/usr/local/bin/und" version)
INSTALLED_UNDCLI_VERSION=$("/usr/local/bin/undcli" version)

if test "${INSTALLED_UND_VERSION}" != "${INSTALLED_UNDCLI_VERSION}"; then
  echo "${bold}WARNING:${normal} und and undcli version mismatch."
  echo "und: v${INSTALLED_UND_VERSION}"
  echo "undcli: v${INSTALLED_UNDCLI_VERSION}"
  echo "Versions should match. Check for errors."
fi
