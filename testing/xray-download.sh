#!/usr/bin/env bash

DAT_PATH=${DAT_PATH:-.}
INSTALL_VERSION="v1.8.4"
NAME=xray

curl() {
  $(type -P curl) -L -q --retry 5 --retry-delay 10 --retry-max-time 60 "$@"
}


identify_the_operating_system_and_architecture() {
  if [[ "$(uname)" == 'Linux' ]]; then
    case "$(uname -m)" in
      'i386' | 'i686')
        MACHINE='32'
        ;;
      'amd64' | 'x86_64')
        MACHINE='64'
        ;;
      'armv5tel')
        MACHINE='arm32-v5'
        ;;
      'armv6l')
        MACHINE='arm32-v6'
        grep Features /proc/cpuinfo | grep -qw 'vfp' || MACHINE='arm32-v5'
        ;;
      'armv7' | 'armv7l')
        MACHINE='arm32-v7a'
        grep Features /proc/cpuinfo | grep -qw 'vfp' || MACHINE='arm32-v5'
        ;;
      'armv8' | 'aarch64')
        MACHINE='arm64-v8a'
        ;;
      'mips')
        MACHINE='mips32'
        ;;
      'mipsle')
        MACHINE='mips32le'
        ;;
      'mips64')
        MACHINE='mips64'
        lscpu | grep -q "Little Endian" && MACHINE='mips64le'
        ;;
      'mips64le')
        MACHINE='mips64le'
        ;;
      'ppc64')
        MACHINE='ppc64'
        ;;
      'ppc64le')
        MACHINE='ppc64le'
        ;;
      'riscv64')
        MACHINE='riscv64'
        ;;
      's390x')
        MACHINE='s390x'
        ;;
      *)
        echo "error: The architecture is not supported."
        exit 1
        ;;
    esac
  else
    echo "error: This operating system is not supported."
    exit 1
  fi
}

download_xray() {
  DOWNLOAD_LINK="https://github.com/XTLS/Xray-core/releases/download/$INSTALL_VERSION/Xray-linux-$MACHINE.zip"
  echo "Downloading Xray archive: $DOWNLOAD_LINK"
  if ! curl -x "${PROXY}" -R -H 'Cache-Control: no-cache' -o "$ZIP_FILE" "$DOWNLOAD_LINK"; then
    echo 'error: Download failed! Please check your network or try again.'
    return 1
  fi
  return 0
  echo "Downloading verification file for Xray archive: $DOWNLOAD_LINK.dgst"
  if ! curl -x "${PROXY}" -sSR -H 'Cache-Control: no-cache' -o "$ZIP_FILE.dgst" "$DOWNLOAD_LINK.dgst"; then
    echo 'error: Download failed! Please check your network or try again.'
    return 1
  fi
  if [[ "$(cat "$ZIP_FILE".dgst)" == 'Not Found' ]]; then
    echo 'error: This version does not support verification. Please replace with another version.'
    return 1
  fi

  # Verification of Xray archive
  CHECKSUM=$(cat "$ZIP_FILE".dgst | awk -F '= ' '/256=/ {print $2}')
  LOCALSUM=$(sha256sum "$ZIP_FILE" | awk '{printf $1}')
  if [[ "$CHECKSUM" != "$LOCALSUM" ]]; then
    echo 'error: SHA256 check failed! Please check your network or try again.'
    return 1
  fi
}

decompression() {
  if ! unzip -q "$1" -d "$TMP_DIRECTORY"; then
    echo 'error: Xray decompression failed.'
    "rm" -r "$TMP_DIRECTORY"
    echo "removed: $TMP_DIRECTORY"
    exit 1
  fi
  echo "info: Extract the Xray package to $TMP_DIRECTORY and prepare it for installation."
}

install_xray() {
  # Install Xray binary to /usr/local/bin/ and $DAT_PATH
  install -m 755 "${TMP_DIRECTORY}/$NAME" "$DAT_PATH/$NAME"
}


main() {
  if [ -f "$DAT_PATH/$NAME" ]; then
    echo "xray already exists"
    exit 0
  fi

  identify_the_operating_system_and_architecture

  # Two very important variables
  TMP_DIRECTORY="$(mktemp -d)"
  ZIP_FILE="${TMP_DIRECTORY}/Xray-linux-$MACHINE.zip"

  if ! download_xray; then
    "rm" -r "$TMP_DIRECTORY"
    echo "removed: $TMP_DIRECTORY"
    exit 1
  fi
  decompression "$ZIP_FILE"

  install_xray
  "rm" -r "$TMP_DIRECTORY"
  echo "removed: $TMP_DIRECTORY"
}

main "$@"
