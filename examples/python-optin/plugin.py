#!/usr/bin/env python3
"""Example Python extension. Build it as a native executable with PyInstaller;
the host only runs the resulting binary."""

import sys

from sahaj_crm_plugin_sdk import serve


def contribute(_request):
    return {
        "fields": [{
            "name": "Enable the Python feature",
            "param": "py_flag",
            "type": "checkbox",
            "note": "Stored under the extension's own data.",
        }],
    }


def handle(event):
    print(f"python extension: {event['phase']} {event['name']}", file=sys.stderr)
    return {}


serve(contribute=contribute, handle=handle)
