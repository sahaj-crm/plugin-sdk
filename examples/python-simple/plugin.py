#!/usr/bin/env python3
"""Minimal Python plugin: one checkbox and a Handle logger."""

import sys

from sahaj_crm_plugin_sdk import as_bool, plugin_values, serve

PLUGIN_ID = "com.example.pysimple"


def contribute(_request):
    return {
        "fields": [{
            "name": "Enable the simple Python feature",
            "key": "enabled",
            "type": "checkbox",
            "note": "Minimal example with one checkbox.",
        }],
    }


def handle(event):
    input_data = event.get("input") or {}
    values = plugin_values(input_data.get("Extras"), PLUGIN_ID)
    print(
        f"simple python plugin: {event['phase']} {event['name']} "
        f"enabled={as_bool(values.get('enabled'))}",
        file=sys.stderr,
    )
    return {}


serve(contribute=contribute, handle=handle)
