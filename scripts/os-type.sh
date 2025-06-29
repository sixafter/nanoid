#!/bin/bash
# Copyright (c) 2024 Six After, Inc.
#
# This source code is licensed under the Apache 2.0 License found in the
# LICENSE file in the root directory of this source tree.

# Detect the platform (similar to $OSTYPE)
# OS=${OSTYPE//[0-9.]/}

UNAME=$( command -v uname)

function detect_os() {
  OS=$( "${UNAME}" | tr '[:upper:]' '[:lower:]')
  case $OS in
    linux*)
      OS='Linux'
      ;;
    msys*|cygwin*|mingw*)
      # or possible 'bash on windows'
      OS='Windows'
      ;;
    nt|win*)
      OS='Windows'
      ;;
    darwin*)
      OS='macOS'
      ;;
    *) ;;
  esac

  echo $OS
}

function is_linux() {
  local OS=$(detect_os)
  if [[ $OS == 'Linux' ]]; then
    return $(true)
  fi

  return $(false)
}

function is_linux_arm() {
  if ! $(is_linux); then
    return $(false)
  fi

  local ARCH
  ARCH=$(uname -m)

  case "$ARCH" in
    arm*|aarch64)
      return $(true)
      ;;
    *)
      return $(false)
      ;;
  esac
}

function is_linux_amd() {
  if ! $(is_linux); then
    return $(false)
  fi

  local ARCH
  ARCH=$(uname -m)

  if [[ "$ARCH" == "x86_64" ]]; then
    return $(true)
  fi

  return $(false)
}

function is_linux_x86() {
  if ! $(is_linux); then
    return $(false)
  fi

  local ARCH
  ARCH=$(uname -m)

  if [[ "$ARCH" == "i386" || "$ARCH" == "i686" ]]; then
    return $(true)
  fi

  return $(false)
}

function is_macos() {
  local OS=$(detect_os)
  if [[ $OS == 'macOS' ]]; then
    return $(true)
  fi

  return $(false)
}

function is_windows() {
  local OS=$(detect_os)
  if [[ $OS == 'Windows' ]]; then
    return $(true)
  fi

  return $(false)
}

function is_macos_arm() {
  local OS=$(detect_os)
  if [[ $OS == 'macOS' && $(uname -p) == 'arm' ]]; then
    return $(true)
  fi

  return $(false)
}

function is_macos_amd() {
  local OS=$(detect_os)
  if [[ $OS == 'macOS' && $(uname -p) == 'i386' ]]; then
    return $(true)
  fi

  return $(false)
}

function is_linux_ubuntu() {
  if ! $(is_linux); then
    return $(false)
  fi

  if [[ -f /etc/os-release ]]; then
    . /etc/os-release
    if [[ $ID == 'ubuntu' ]]; then
      return $(true)
    fi
  fi

  return $(false)
}

function is_wsl() {
  if ! $(is_linux); then
    return $(false)
  fi

  # WSL detection based on /proc/version or osrelease
  if grep -qiE 'microsoft|wsl' /proc/version 2>/dev/null || \
     grep -qiE 'microsoft|wsl' /proc/sys/kernel/osrelease 2>/dev/null; then
    return $(true)
  fi

  return $(false)
}
