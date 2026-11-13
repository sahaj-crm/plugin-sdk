#!/usr/bin/env python3
"""Example Python extension. Build it as a native executable with PyInstaller;
the host only runs the resulting binary."""

import sys

from sahaj_crm_plugin_sdk import serve


def contribute(_request):
    return {
        "fields": [{
            "name": "Enable the Python feature",
            "param": "ext[com.example.pyoptin][py_flag]",
            "type": "checkbox",
        }],
    }


def handle(event):
    print(f"python extension: {event['phase']} {event['name']}", file=sys.stderr)
    return {}


serve(contribute=contribute, handle=handle)
