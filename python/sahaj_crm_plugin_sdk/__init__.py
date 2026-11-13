"""Sahaj CRM plugin SDK.

Plugins are executable processes speaking newline-delimited JSON-RPC over
stdin/stdout. `serve` implements the loop; extension authors provide `info`,
`contribute` and `handle` callables (or an object with those methods).
"""

import json
import sys
from typing import Any, Callable, Optional, TypedDict

JSONRPC_VERSION = "2.0"

Extras = dict[str, dict[str, Any]]


class EntityInput(TypedDict):
    """Input shape for entity create/update events."""

    Record: dict[str, Any]
    Extras: Extras


def extension_values(extras: Optional[Extras], extension_id: str) -> dict[str, Any]:
    """Return one extension's values from an event's Extras map."""
    if not extras:
        return {}
    return extras.get(extension_id) or {}


def as_bool(value: Any) -> bool:
    """Normalize form/API boolean values."""
    return value is True or value == 1 or (
        isinstance(value, str) and value.lower() in ("true", "1")
    )


def as_string(value: Any) -> str:
    """Normalize a single form/API value."""
    return "" if value is None else str(value)


def as_strings(value: Any) -> list[str]:
    """Normalize single or multi-value form/API values."""
    if value is None or value == "":
        return []
    if isinstance(value, list):
        return [str(item) for item in value]
    return [str(value)]


def serve(
    contribute: Optional[Callable[[dict], Any]] = None,
    handle: Optional[Callable[[dict], Any]] = None,
    extension: Any = None,
) -> None:
    """Run an extension until stdin closes or a shutdown request arrives."""
    if extension is not None:
        contribute = contribute or extension.contribute
        handle = handle or extension.handle
    if contribute is None or handle is None:
        raise ValueError("serve requires contribute and handle callables")

    for line in sys.stdin:
        if not line.strip():
            continue
        try:
            request = json.loads(line)
        except json.JSONDecodeError as exc:
            _error(0, -32700, "invalid JSON: " + str(exc))
            continue

        request_id = request.get("id", 0)
        method = request.get("method")
        params = request.get("params") or {}

        try:
            if method == "contribute":
                result = contribute(params)
            elif method == "handle":
                result = handle(params)
            elif method == "shutdown":
                _reply(request_id, {"ok": True})
                break
            else:
                _error(request_id, -32601, "unknown method: " + str(method))
                continue
            _reply(request_id, result)
        except Exception as exc:  # noqa: BLE001 - protocol boundary
            _error(request_id, -32000, str(exc))


def _reply(request_id: int, result: Any) -> None:
    print(
        json.dumps({"jsonrpc": JSONRPC_VERSION, "id": request_id, "result": result}),
        flush=True,
    )


def _error(request_id: int, code: int, message: str) -> None:
    print(
        json.dumps({
            "jsonrpc": JSONRPC_VERSION,
            "id": request_id,
            "error": {"code": code, "message": message},
        }),
        flush=True,
    )
