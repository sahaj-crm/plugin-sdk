#!/usr/bin/env python3
"""Example Python extension. Build it as a native executable with PyInstaller;
the host only runs the resulting binary."""

import sys

from sahaj_crm_plugin_sdk import (
    as_bool,
    as_string,
    extension_values,
    serve,
)

EXTENSION_ID = "com.example.pyoptin"


def contribute(_request):
    return {
        "fields": [{
            "name": "Enable the Python feature",
            "param": "feature_enabled",
            "type": "checkbox",
            "note": "Stored under the extension's own data.",
        }, {
            "name": "Choose wisely",
            "param": "feature_mode",
            "type": "select",
            "note": "Python select",
            "options": [
                {"value": "one", "label": "1"},
                {"value": "two", "label": "2"},
            ],
        }],
    }


def handle(event):
    input_data = event.get("input") or {}
    values = extension_values(input_data.get("Extras"), EXTENSION_ID)
    print(
        f"python extension: {event['phase']} {event['name']} "
        f"feature_enabled={as_bool(values.get('feature_enabled'))} "
        f"feature_mode={as_string(values.get('feature_mode'))}",
        file=sys.stderr,
    )

    if event["phase"] == "after":
        output = event.get("output") or {}
        print(f"python extension: saved id={output.get('ID')}", file=sys.stderr)

    return {}


serve(contribute=contribute, handle=handle)
