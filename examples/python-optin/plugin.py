#!/usr/bin/env python3
"""Example Python extension. The SDK is vendored in the plugin directory, so
this can be copied directly into the host's plugins folder."""

import sys

from sahaj_crm_plugin_sdk import serve


def info():
    return {
        "id": "com.example.pyoptin",
        "name": "Python Example",
        "version": "0.1.0",
        "api_version": 1,
        "targets": ["course.new", "course.edit"],
        "events": ["course.create", "course.update"],
    }


def contribute(_request):
    return {
        "sections": [{
            "title": "Python Options",
            "fields": [{
                "key": "py_flag",
                "label": "Enable the Python feature",
                "type": "checkbox",
            }],
        }],
    }


def handle(event):
    print(f"python extension: {event['phase']} {event['name']}", file=sys.stderr)
    return {}


serve(info=info, contribute=contribute, handle=handle)
