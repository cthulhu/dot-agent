#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <tag>" >&2
  echo "Example: $0 v0.1.0" >&2
  exit 1
fi

tag="$1"
version="${tag#v}"
url="https://github.com/cthulhu/dot-agent/archive/refs/tags/${tag}.tar.gz"

echo "Fetching ${url} ..."
curl -fsSL "$url" | shasum -a 256

echo
echo "Update Formula/dot-agent.rb:"
echo "  url \"${url}\""
echo "  version \"${version}\""
