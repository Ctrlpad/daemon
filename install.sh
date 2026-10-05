#!/bin/sh
set -e

OWNER="ctrlpad"
REPO="daemon"
BINARY="ctrlpad-daemon"
PROJECT_NAME="ctrlpad-daemon"

usage() {
  this=$1
  cat <<EOF
$this: download ${PROJECT_NAME} binary from GitHub releases and use it as a systemd service

Usage: $this [-b bindir] [-d] [tag]
  -b sets bindir or installation directory, Defaults to /usr/local/bin
  -d turns on debug logging
  [tag] is a tag from https://github.com/${OWNER}/${REPO}/releases
        If tag is missing, then the latest will be used.

Examples:
  # Default installation (system-wide, requires sudo)
  curl -f https://raw.githubusercontent.com/${OWNER}/${REPO}/main/install.sh | sudo sh

  # Install specific version
  curl -f https://raw.githubusercontent.com/${OWNER}/${REPO}/main/install.sh | sudo sh -s -- v0.0.1

EOF
  exit 2
}

parse_args() {
  BINDIR=${BINDIR:-"/usr/local/bin"}
  while getopts "b:dh?x" arg; do
    case "$arg" in
      b) BINDIR="$OPTARG" ;;
      d) log_set_priority 10 ;;
      h | \?) usage "$0" ;;
      x) set -x ;;
    esac
  done
  shift $((OPTIND - 1))
  TAG=$1
}

execute() {
  tmpdir=$(mktemp -d)
  log_debug "downloading files into ${tmpdir}"
  http_download "${tmpdir}/${TARBALL}" "${TARBALL_URL}"
  
  srcdir="${tmpdir}"
  (cd "${tmpdir}" && untar "${TARBALL}")
  
  test ! -d "${BINDIR}" && install -d "${BINDIR}"
  
  FOUND_BIN=$(find "${tmpdir}" -type f -executable ! -name "*.tar.gz" ! -name "*.tgz" ! -name "*.zip" -print -quit)

  if [ -n "$FOUND_BIN" ]; then
    log_info "Found binary: $(basename "$FOUND_BIN") -> installing as ${BINDIR}/${BINARY}"
    install "$FOUND_BIN" "${BINDIR}/${BINARY}"
  else
    log_crit "No executable binary found in archive"
    rm -rf "${tmpdir}"
    exit 1
  fi
  
  rm -rf "${tmpdir}"
}

is_command() {
  command -v "$1" >/dev/null
}

echoerr() {
  echo "$@" 1>&2
}

_logp=6
log_set_priority() {
  _logp="$1"
}

log_priority() {
  if test -z "$1"; then
    echo "$_logp"
    return
  fi
  [ "$1" -le "$_logp" ]
}

log_debug() {
  log_priority 7 || return 0
  echoerr "$@"
}

log_info() {
  log_priority 6 || return 0
  echoerr "$@"
}

log_crit() {
  log_priority 2 || return 0
  echoerr "$@"
}

uname_os() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  if [ "$os" != "linux" ]; then
    log_crit "This script only supports Linux. Detected OS: $os"
    exit 1
  fi
  echo "$os"
}

uname_arch() {
  arch=$(uname -m)
  case $arch in
    x86_64) arch="amd64" ;;
    x86) arch="386" ;;
    i686) arch="386" ;;
    i386) arch="386" ;;
    aarch64) arch="arm64" ;;
    armv5*) arch="armv5" ;;
    armv6*) arch="armv6" ;;
    armv7*) arch="armv7" ;;
  esac
  echo "${arch}"
}

untar() {
  tarball=$1
  case "${tarball}" in
    *.tar.gz | *.tgz) tar --no-same-owner -xzf "${tarball}" ;;
    *.tar) tar --no-same-owner -xf "${tarball}" ;;
    *)
      echoerr "untar: unknown archive format for ${tarball}"
      return 1
      ;;
  esac
}

http_download_curl() {
  local_file=$1
  source_url=$2
  header=$3
  if [ -z "$header" ]; then
    code=$(curl -w '%{http_code}' -sL -o "$local_file" "$source_url")
  else
    code=$(curl -w '%{http_code}' -sL -H "$header" -o "$local_file" "$source_url")
  fi
  if [ "$code" != "200" ]; then
    log_crit "http_download_curl received HTTP status $code"
    return 1
  fi
  return 0
}

http_download() {
  log_debug "http_download $2"
  if is_command curl; then
    http_download_curl "$@"
    return
  fi
  log_crit "http_download unable to find curl"
  return 1
}

http_copy() {
  tmp=$(mktemp)
  http_download "${tmp}" "$1" "$2" || return 1
  body=$(cat "$tmp")
  rm -f "${tmp}"
  echo "$body"
}

github_release() {
  owner_repo=$1
  version=$2
  
  if [ -z "$version" ] || [ "$version" = "latest" ]; then
    giturl="https://api.github.com/repos/${owner_repo}/releases/latest"
  else
    giturl="https://api.github.com/repos/${owner_repo}/releases/tags/${version}"
  fi

  json=$(http_copy "$giturl")
  test -z "$json" && return 1
  
  version=$(echo "$json" | tr -s '\n' ' ' | sed 's/.*"tag_name":[ ]*"//' | sed 's/".*//')
  test -z "$version" && return 1
  echo "$version"
}

tag_to_version() {
  if [ -z "${TAG}" ]; then
    log_info "checking GitHub for latest tag"
  else
    log_info "checking GitHub for tag '${TAG}'"
  fi
  REALTAG=$(github_release "$OWNER/$REPO" "${TAG}") && true
  if test -z "$REALTAG"; then
    log_crit "unable to find '${TAG}' - use 'latest' or see https://github.com/${OWNER}/${REPO}/releases for details"
    exit 1
  fi
  TAG="$REALTAG"
  VERSION=${TAG#v}
}

adjust_os() {
  case ${OS} in
    linux) ADJUSTED_OS="Linux" ;;
    *) ADJUSTED_OS=$(echo "${OS}" | sed 's/./\u&/') ;;
  esac
}

adjust_arch() {
  case ${ARCH} in
    amd64) ADJUSTED_ARCH="x86_64" ;;
    386) ADJUSTED_ARCH="i386" ;;
    arm64) ADJUSTED_ARCH="arm64" ;;
    *) ADJUSTED_ARCH="${ARCH}" ;;
  esac
}

FORMAT=tar.gz
OS=$(uname_os)
ARCH=$(uname_arch)

parse_args "$@"

tag_to_version
adjust_os
adjust_arch

log_info "found version: ${VERSION} for ${TAG}/${OS}/${ARCH}"

NAME=${PROJECT_NAME}_${ADJUSTED_OS}_${ADJUSTED_ARCH}
TARBALL=${NAME}.${FORMAT}
TARBALL_URL=https://github.com/${OWNER}/${REPO}/releases/download/${TAG}/${TARBALL}

execute
