#!/bin/sh
# Build and install the example plugins as native executables.
#
#   ./install-examples.sh /path/to/sahaj-crm/plugins
#
# Requirements:
#   - Go
#   - Bun (for the JavaScript examples)
#   - Python 3 with PyInstaller (for the Python examples)
set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
TARGET=${1:-"$ROOT/../sahaj-crm/plugins"}
BUILD=$(mktemp -d)
trap 'rm -rf "$BUILD"' EXIT

mkdir -p \
  "$TARGET/go-simple" "$TARGET/go-full" \
  "$TARGET/js-simple" "$TARGET/js-full" \
  "$TARGET/python-simple" "$TARGET/python-full"

cd "$ROOT"

echo "Building Go examples..."
go build -o "$TARGET/go-simple/plugin" ./examples/go-simple
cp examples/go-simple/plugin.json "$TARGET/go-simple/"
go build -o "$TARGET/go-full/plugin" ./examples/go-full
cp examples/go-full/plugin.json "$TARGET/go-full/"

echo "Building JavaScript examples with Bun..."
if ! command -v bun >/dev/null 2>&1; then
  echo "error: bun is required (https://bun.sh)" >&2
  exit 1
fi
bun build --compile examples/js-simple/plugin.mjs \
  --outfile "$TARGET/js-simple/plugin"
cp examples/js-simple/plugin.json "$TARGET/js-simple/"
bun build --compile examples/js-full/plugin.mjs \
  --outfile "$TARGET/js-full/plugin"
cp examples/js-full/plugin.json "$TARGET/js-full/"

echo "Building Python examples with PyInstaller..."
if ! python3 -c 'import PyInstaller' >/dev/null 2>&1; then
  echo "error: PyInstaller is required (python3 -m pip install pyinstaller)" >&2
  exit 1
fi
for example in python-simple python-full; do
  mkdir -p "$BUILD/$example"
  python3 -m PyInstaller \
    --onefile --clean --noconfirm \
    --name plugin \
    --distpath "$TARGET/$example" \
    --workpath "$BUILD/$example" \
    --specpath "$BUILD/$example" \
    --paths "$ROOT/python" \
    "$ROOT/examples/$example/plugin.py"
  cp "examples/$example/plugin.json" "$TARGET/$example/"
done

echo "Installed native example plugins into $TARGET"
echo "Restart the CRM process to load them."
