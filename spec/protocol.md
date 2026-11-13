# Sahaj Plugin Protocol v1

A plugin is one directory containing a native executable and `plugin.json`.
The host runs the executable directly; there is no runtime handshake.

## Manifest

```json
{
  "id": "com.example.optin",
  "name": "Example Opt-in",
  "version": "0.1.0",
  "api_version": 1,
  "executable": "plugin",
  "targets": ["course.new", "course.edit"],
  "events": ["course.create", "course.update"]
}
```

`targets` and `events` support exact names and a trailing `*` wildcard.
`api_version` is validated by the host before the plugin is loaded.

## Transport rules

- UTF-8 JSON, one object per line.
- stdout contains protocol messages only.
- calls are serialized by the host; a plugin may process them sequentially.
- the plugin process is started once and reused.
- closing stdin or receiving `shutdown` ends the process.
- requests use JSON-RPC 2.0-style envelopes.

## Request

```json
{"jsonrpc":"2.0","id":1,"method":"contribute","params":{}}
```

## Response

```json
{"jsonrpc":"2.0","id":1,"result":{}}
```

```json
{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"reason"}}
```

## `contribute`

Params:

```json
{
  "entity": "course",
  "action": "new",
  "id": "new",
  "record": {"Name": "New Course"},
  "values": {"opt_in": true}
}
```

`values` contains only this extension's stored data. `record` is a snapshot of
the entity model.

Result:

```json
{
  "fields": [{
    "name": "Enable the example feature",
    "param": "ext[com.example.optin][opt_in]",
    "type": "checkbox",
    "value": true,
    "mode": "",
    "placeholder": "",
    "options": []
  }]
}
```

`param` is the full HTML input name and its value is stored under that key.
The host overlays stored `values` by `param`. Supported field types:
`checkbox`, `checkboxes`, `text`, `textarea`, `select`, `number`.

## `handle`

Params:

```json
{
  "name": "course.create",
  "phase": "before",
  "entity": "course",
  "action": "create",
  "id": "",
  "input": {"Record": {}, "Extras": {}},
  "output": null,
  "error": ""
}
```

Result:

```json
{"error": ""}
```

A non-empty `error` from a `before` handler aborts the operation. Errors from
`after` and `failed` handlers are logged by the host and do not change the
operation result.

## `shutdown`

Params: `{}`

Result:

```json
{"ok": true}
```

The plugin should exit after replying. The host may also close stdin.
