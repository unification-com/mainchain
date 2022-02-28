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

#LATEST_RELEASE="https://api.github.com/repos/unification-com/mainchain/releases/latest"
GH_DL_PREFIX="https://github.com/unification-com/mainchain/releases/download"
#LATEST_VERSION="$(curl --silent ${LATEST_RELEASE} | grep -Po '"tag_name": "\K.*?(?=")')"
DOWNLOAD_DEST="/tmp/mainchain_tmp"
LOCAL_BIN="/usr/local/bin"
BACKUP_DIR="${HOME}/UND_BIN_BAK"
#SPECIFIED_VERSION=${1}
#VERSION_TO_INSTALL=${LATEST_VERSION}

BINARIES_TO_INSTALL="und undcli"

INSTALLATION_SUMMARY=""

SKIP_INSTALL=false
REQUIRE_BACKUP=false

bold=$(tput bold)
normal=$(tput sgr0)

#echo "Latest release is v${LATEST_VERSION}"

#if [ -n "$SPECIFIED_VERSION" ]; then
#  echo "user requested v${SPECIFIED_VERSION}"
#  VERSION_TO_INSTALL=${SPECIFIED_VERSION}
#fi

VERSION_TO_INSTALL="1.4.8"

echo ""
echo "------------------"
echo "DEPRECATION NOTICE"
echo "------------------"
echo ""
echo "This script will be deprecated once und v1.5.0 has been released"
echo "on MainNet. Once the network has been upgraded to Cosmos SDK v0.42.x,"
echo "cosmovisor will be the preferred method for und binary upgrades."
echo ""
echo "The last version supported by this script is:"
echo ""
echo "und v${VERSION_TO_INSTALL}"
echo ""
echo "No version higher than this should be used to run a node on either:"
echo ""
echo "FUND-Mainchain-MainNet-v1 or FUND-Mainchain-TestNet-v12"
echo ""

read -rp "Press [enter] to continue installing und v${VERSION_TO_INSTALL} or ^C to exit" key

# check binary
# check und
check() {
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
    if test "${VERSION_TO_INSTALL}" = "${VERS}"; then
      MSG="version ${1} v${VERSION_TO_INSTALL} already installed in ${CHECK_LOC}"
      INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${MSG}. Nothing to do."
      echo "${bold}NOTICE:${normal} ${MSG}. Skipping."
      SKIP_INSTALL=true
    else
      if [ "$(echo "${VERS} ${VERSION_TO_INSTALL}" | tr " " "\n" | sort --version-sort | head -n 1)" = "${VERSION_TO_INSTALL}" ]; then
        echo "${bold}WARNING: installed version ${VERS} is newer than requested version ${VERSION_TO_INSTALL}.\nContinue with installation? (y to continue)${normal}"
        read -r INSTALL_OLD
        if test "${INSTALL_OLD}" != "y"; then
          MSG="requested version ${1} v${VERSION_TO_INSTALL} older than installed ${1} v${VERS}. Installation canceled by user"
          echo "${MSG}"
          INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${MSG}"
          SKIP_INSTALL=true
        else
          echo "Continue installing ${VERSION_TO_INSTALL}"
          echo "flag ${CHECK_LOC} ${VERS} for backup"
          REQUIRE_BACKUP=true
        fi
      else
        echo "Requested version newer than installed version. Continue"
        echo "flag ${CHECK_LOC} ${VERS} for backup"
        REQUIRE_BACKUP=true
      fi
    fi
  fi
}

backup() {
  BACKUP_VERS=$("${LOCAL_BIN}/${1}" version)
  echo "Atempting to backup ${LOCAL_BIN}/${1} ${BACKUP_VERS}"
  BACKUP_LOC="${BACKUP_DIR}"
  if [ ! -d "$BACKUP_LOC" ]; then
    mkdir -p "${BACKUP_LOC}"
  fi
  BACKUP_BIN="${BACKUP_LOC}/${1}-v${BACKUP_VERS}.tar.gz"
  if [ -f "$BACKUP_BIN" ]; then
    echo "${bold}NOTE: Backup ${BACKUP_BIN} already exists. Overwrite? (y to overwrite)${normal}"
    read -r OVERWRITE_BACKUP
    if test "${OVERWRITE_BACKUP}" != "y"; then
      echo "User requested skip backup for ${LOCAL_BIN}/${1} ${BACKUP_VERS}"
      MSG="${LOCAL_BIN}/${1} ${BACKUP_VERS} not backed up - backup already exists from previous install."
      INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${MSG}\n"
      return
    fi
  fi
  echo "Backing up ${1} v${BACKUP_VERS} to ${BACKUP_BIN}"
  tar -czf "${BACKUP_BIN}" -C "${LOCAL_BIN}" "${1}"

  if [ -f "$BACKUP_BIN" ]; then
    MSG="${LOCAL_BIN}/${1} v${VERS} successfully backed up to: ${BACKUP_BIN}"
    echo "${MSG}"
    INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${MSG}\n"
  else
    echo "\n${bold}Failed to backup ${1} v${BACKUP_VERS}. Continue installation? (y to continue)${normal}"
    read -r CONTINUE_INSTALLATION
    if test "${CONTINUE_INSTALLATION}" = "y"; then
      echo "Continuing installation without backup..."
      INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${bold}WARNING:${normal} ${LOCAL_BIN}/${1} v${BACKUP_VERS} was NOT backed up.\n"
    else
      echo "Abort ${1} installation."
      INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\nInstallation of ${1} v${BACKUP_VERS} was aborted by user.\n"
      SKIP_INSTALL=true
    fi
  fi
}

# check if requested version exists
exists() {
  DOWLOAD_ARCHIVE="${2}/${1}_v${2}_${3}_x86_64.tar.gz"
  DOWLOAD_FILE="${GH_DL_PREFIX}/${DOWLOAD_ARCHIVE}"
  echo "Checking if ${1} v${2} has been released"
  if curl --output /dev/null --silent --head --fail "${DOWLOAD_FILE}"; then
    echo "Requested version ${1} v${2} is a release version."
  else
    echo "Requested version ${1} v${2} is not a release version."
    echo "Skipping ${1} installation."
    INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${1} v${2} does not exist as a release version. Installation skipped.\n"
    SKIP_INSTALL=true
  fi
}

# download binary_name version os
# download "undcli" "1.1.1" "linux"
download() {
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
extract() {
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
install() {
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

# Wrap it all in a run function
run() {

  if [ -z "$HOME" ]; then
    echo "${bold}WARNING: NO HOME ENVIRONMENT VARIABLE DETECTED! ABORTING.${normal}"
    exit
  fi

  ## create a tmp download dir
  echo "creating ${DOWNLOAD_DEST}"
  mkdir -p "${DOWNLOAD_DEST}"
  echo ""

  for i in ${BINARIES_TO_INSTALL}; do
    SKIP_INSTALL=false
    REQUIRE_BACKUP=false
    INSTALL_LOCATION=""
    INSTALLED_VERSION=""

    echo ""
    echo "${bold}--Attempting to install/update ${i}--${normal}"

    # if user requested a particular version, check it exists first
    if test "${VERSION_TO_INSTALL}" != "${LATEST_VERSION}"; then
      exists "${i}" "${VERSION_TO_INSTALL}" "linux"
      if [ "$SKIP_INSTALL" = true ]; then
        continue
      fi
    fi

    # installation pre-check
    check "${i}"

    if [ "$REQUIRE_BACKUP" = true ]; then
      backup "${i}"
    fi

    echo ""
    if [ "$SKIP_INSTALL" = true ]; then
      continue
    fi

    echo "Installing ${i}"
    echo ""
    # Download the latest archive
    download "${i}" "${VERSION_TO_INSTALL}" "linux"
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
    echo "expected version: ${VERSION_TO_INSTALL}"
    echo "installed version: ${INSTALLED_VERSION}"

    if test "${VERSION_TO_INSTALL}" = "${INSTALLED_VERSION}"; then
      INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${bold}${i} v${INSTALLED_VERSION} successfully installed!${normal}\nLocation: ${INSTALL_LOCATION}"
      INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\nRun:\n\n${INSTALL_LOCATION} version --long\n\nto verify. Version should be ${VERSION_TO_INSTALL}.\n"
    else
      INSTALLATION_SUMMARY="${INSTALLATION_SUMMARY}\n${bold}Something went wrong installing ${i} v${VERSION_TO_INSTALL}. See previous output.${normal}"
    fi
  done

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
}

run
