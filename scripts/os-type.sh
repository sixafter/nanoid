#!/bin/bash
# Copyright (c) 2024 Six After, Inc.
#
# This source code is licensed under the Apache 2.0 License found in the
# LICENSE file in the root directory of this source tree.

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
