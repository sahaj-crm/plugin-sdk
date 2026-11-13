#!/bin/sh
# Build and install the example plugins as native executables.
#
#   ./install-examples.sh /path/to/sahaj-crm/plugins
#
# Requirements:
#   - Go
#   - Bun (for the JavaScript example)
#   - Python 3 with PyInstaller (for the Python example)
set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
TARGET=${1:-"$ROOT/../sahaj-crm/plugins"}
BUILD=$(mktemp -d)
trap 'rm -rf "$BUILD"' EXIT

mkdir -p "$TARGET/go-simple" "$TARGET/go-optin" "$TARGET/js-optin" "$TARGET/python-optin"

echo "Building simple Go example..."
cd "$ROOT"
go build -o "$TARGET/go-simple/plugin" ./examples/go-simple
cp examples/go-simple/plugin.json "$TARGET/go-simple/"

echo "Building Go example..."
go build -o "$TARGET/go-optin/plugin" ./examples/go-optin
cp examples/go-optin/plugin.json "$TARGET/go-optin/"

echo "Building JavaScript example with Bun..."
if ! command -v bun >/dev/null 2>&1; then
  echo "error: bun is required (https://bun.sh)" >&2
  exit 1
fi
bun build --compile examples/js-optin/plugin.mjs \
  --outfile "$TARGET/js-optin/plugin"
cp examples/js-optin/plugin.json "$TARGET/js-optin/"

echo "Building Python example with PyInstaller..."
if ! python3 -c 'import PyInstaller' >/dev/null 2>&1; then
  echo "error: PyInstaller is required (python3 -m pip install pyinstaller)" >&2
  exit 1
fi
python3 -m PyInstaller \
  --onefile --clean --noconfirm \
  --name plugin \
  --distpath "$TARGET/python-optin" \
  --workpath "$BUILD" \
  --specpath "$BUILD" \
  --paths "$ROOT/python" \
  "$ROOT/examples/python-optin/plugin.py"
cp examples/python-optin/plugin.json "$TARGET/python-optin/"

echo "Installed native example plugins into $TARGET"
echo "Restart the CRM process to load them."
