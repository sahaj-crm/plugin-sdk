#!/usr/bin/env python3
"""Example Python extension. Install the SDK with `pip install -e python` or
run this from the repository so the relative import below finds it."""

import pathlib
import sys

sys.path.insert(0, str(pathlib.Path(__file__).resolve().parents[2] / "python"))

from sahaj_crm_plugin_sdk import serve  # noqa: E402


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
