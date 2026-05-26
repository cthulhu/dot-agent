#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <tag>" >&2
  echo "Example: $0 v0.7.0" >&2
  exit 1
fi

tag="$1"
version="${tag#v}"
url="https://github.com/cthulhu/dot-agent/releases/download/${tag}/dot-agent_${version}_windows_amd64.zip"

echo "Fetching ${url} ..."
sha=$(curl -fsSL "$url" | shasum -a 256 | awk '{print $1}')

echo
echo "SHA-256 Checksum: ${sha}"
echo
echo "Update chocolatey/dot-agent.nuspec:"
echo "  <version>${version}</version>"
echo
echo "Update chocolatey/tools/chocolateyInstall.ps1:"
echo "  checksum      = '${sha}'"
