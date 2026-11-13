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

`values` contains only this extension's stored data, keyed by the simple field
`param`. `record` is a snapshot of the entity model.

Result:

```json
{
  "fields": [{
    "name": "Enable the example feature",
    "param": "opt_in",
    "type": "checkbox",
    "value": true,
    "mode": "",
    "placeholder": "",
    "note": "Stored under the extension's own data.",
    "options": []
  }]
}
```

`param` is the extension-local field key. The host namespaces it to
`ext[<extension>][<param>]` before rendering and stores the value under the
simple key in the extension's own data row. The host overlays stored `values`
by `param`. `note` renders below the input. Supported field types:
`checkbox`, `checkboxes`, `text`, `textarea`, `select`, `number`.

## `handle`

Params:

```json
{
  "name": "course.update",
  "phase": "before",
  "entity": "course",
  "action": "update",
  "id": "4",
  "input": {
    "Record": {
      "ID": 4,
      "Name": "Introduction to Sahaj",
      "StartDate": "2026-01-01"
    },
    "Extras": {
      "com.example.optin": {
        "opt_in": "true",
        "options": "one"
      }
    }
  },
  "output": null,
  "error": ""
}
```

Result:

```json
{"error": ""}
```

### Phases

- `before` runs before the database operation. Return a non-empty `error` to
  abort. No writes have happened yet.
- `after` runs after the operation commits. `output` contains the saved result.
- `failed` runs after an operation error. `error` contains the reason.

Errors from `after` and `failed` handlers are logged and do not change the
operation result.

### Input

`input` is the JSON form of the operation's input.

For entity **create** and **update**, it is:

```json
{
  "Record": { "...": "entity fields..." },
  "Extras": {
    "<extension id>": {
      "<field param>": "value"
    }
  }
}
```

`Extras` contains every extension field submitted with this operation, grouped
by extension ID, then by the field's simple `param`.

- Web forms submit strings. A checkbox is `"true"` or `"false"`; multi-value
  fields are arrays of strings.
- API clients may submit JSON types directly: booleans, numbers, strings, and
  arrays.
- A field that was not submitted is absent.
- Plugins should read only `Extras["<their own id>"]`.

`Record` keys use the host model's JSON names. Today those are Go field names
such as `Name`, `StartDate`, and `Status`.

For archive and delete, `input` is just the entity ID. There are no `Extras`.

### Output

`output` is only present on `after`:

- create/update/archive: the resulting entity
- delete: `{"ID": <id>}`

The host does not consume `output`. It is provided so `after` handlers can use
the committed result for their own side effects. For a course create it is the
course itself, not the generated sessions/list/form/campaign/page.

### EventResponse

```json
{"error": "", "values": {}}
```

- `error` aborts a `before` event when non-empty.
- `values` is reserved for future value patching and is currently ignored by
  the host.

## `shutdown`

Params: `{}`

Result:

```json
{"ok": true}
```

The plugin should exit after replying. The host may also close stdin.
