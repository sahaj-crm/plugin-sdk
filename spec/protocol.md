# Sahaj Plugin Protocol v1

Plugins are executable child processes. The host writes one JSON object per
line to stdin and reads one JSON object per line from stdout. stderr is for
logs and is never parsed.

## Transport rules

- UTF-8 JSON, one object per line.
- stdout contains protocol messages only.
- calls are serialized by the host; a plugin may process them sequentially.
- the plugin process is started once and reused.
- closing stdin or receiving `shutdown` ends the process.
- requests use JSON-RPC 2.0-style envelopes.

## Request

```json
{"jsonrpc":"2.0","id":1,"method":"info","params":{}}
```

## Response

```json
{"jsonrpc":"2.0","id":1,"result":{}}
```

```json
{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"reason"}}
```

## `info`

Params: `{}`

Result:

```json
{
  "id": "com.example.optin",
  "name": "Example Opt-in",
  "version": "0.1.0",
  "api_version": 1,
  "targets": ["course.new", "course.edit"],
  "events": ["course.create", "course.update"]
}
```

`targets` and `events` support exact names and a trailing `*` wildcard.

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
  "sections": [{
    "extension": "com.example.optin",
    "title": "Example Options",
    "fields": [{
      "extension": "com.example.optin",
      "key": "opt_in",
      "label": "Enable the example feature",
      "type": "checkbox",
      "value": true,
      "options": [],
      "required": false
    }]
  }]
}
```

The host fills empty `extension` values and overlays stored `values`. Supported
field types: `checkbox`, `text`, `textarea`, `select`, `number`.

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

## Versioning

`api_version` is checked by the host during `info`. A plugin with an
unsupported version is not loaded.
