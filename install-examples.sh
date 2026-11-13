#!/bin/sh
# Install the example plugins into a Sahaj CRM plugins directory.
#
#   ./install-examples.sh /path/to/sahaj-crm/plugins
#
# The default target is ../sahaj-crm/plugins, which matches the workspace
# layout used by the two repositories.
set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
TARGET=${1:-"$ROOT/../sahaj-crm/plugins"}

mkdir -p "$TARGET/go-optin" "$TARGET/js-optin" "$TARGET/python-optin"

cd "$ROOT"
go build -o "$TARGET/go-optin/plugin" ./examples/go-optin

cp examples/go-optin/plugin.json "$TARGET/go-optin/"
cp examples/js-optin/plugin.json examples/js-optin/plugin.mjs examples/js-optin/sdk.js "$TARGET/js-optin/"
cp examples/python-optin/plugin.json examples/python-optin/plugin.py "$TARGET/python-optin/"
cp -R examples/python-optin/sahaj_crm_plugin_sdk "$TARGET/python-optin/"
rm -rf "$TARGET/python-optin/__pycache__"

echo "Installed example plugins into $TARGET"
echo "Restart the CRM process to load them."
