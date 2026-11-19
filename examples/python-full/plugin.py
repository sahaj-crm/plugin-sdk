#!/usr/bin/env python3
"""Example Python plugin. Build it as a native executable with PyInstaller;
the host only runs the resulting binary."""

import sys

from sahaj_crm_plugin_sdk import (
    as_bool,
    as_string,
    plugin_values,
    serve,
)

PLUGIN_ID = "com.example.pyfull"


def contribute(request):
    settings = request.get("settings") or {}
    default_mode = as_string(settings.get("default_mode")) or "one"
    return {
        "fields": [{
            "name": "Enable the full Python feature",
            "key": "feature_enabled",
            "type": "checkbox",
            "note": "Stored under the plugin's own data.",
        }, {
            "name": "Full Python mode",
            "key": "feature_mode",
            "type": "select",
            "value": default_mode,
            "note": "Defaults from the plugin setting.",
            "options": [
                {"value": "one", "label": "Mode 1"},
                {"value": "two", "label": "Mode 2"},
            ],
        }],
    }


def handle(event):
    input_data = event.get("input") or {}
    values = plugin_values(input_data.get("Extras"), PLUGIN_ID)
    settings = event.get("settings") or {}
    print(
        f"full python plugin: {event['phase']} {event['name']} "
        f"feature_enabled={as_bool(values.get('feature_enabled'))} "
        f"feature_mode={as_string(values.get('feature_mode'))} "
        f"settings.default_mode={as_string(settings.get('default_mode'))} "
        f"settings.api_key_configured={bool(as_string(settings.get('api_key')))}",
        file=sys.stderr,
    )

    if event["phase"] == "after":
        output = event.get("output") or {}
        print(f"full python plugin: saved id={output.get('ID')}", file=sys.stderr)

    return {}


serve(contribute=contribute, handle=handle)
