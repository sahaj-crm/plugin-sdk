# Sahaj CRM Plugin SDK

SDK and protocol for Sahaj CRM extensions.

A plugin is **one directory containing one native executable and one
`plugin.json` manifest**:

```text
plugins/example-optin/
  plugin.json
  plugin
```

The host runs the executable directly. It does not care whether it was written
in Go, JavaScript, Python, Rust, or anything else.

## Manifest

The manifest is authoritative for plugin metadata. There is no runtime `info`
handshake.

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

`executable` is relative to the plugin directory. On Windows the host also
tries `<executable>.exe` when the exact file is missing.

## Protocol

One JSON object per line over stdin/stdout:

```json
{"jsonrpc":"2.0","id":1,"method":"contribute","params":{}}
```

Methods:

- `contribute` — declarative fields for a modal
- `handle` — before/after/failed lifecycle events
- `shutdown` — optional graceful exit

The full contract is in [spec/protocol.md](spec/protocol.md).

## Fields

`contribute` returns a flat list of fields:

```json
{
  "fields": [{
    "name": "Enable the example feature",
    "param": "opt_in",
    "type": "checkbox",
    "note": "Stored under the extension's own data."
  }]
}
```

- `param` is the extension-local field key. The host namespaces it to
  `ext[<extension>][<param>]` for the form and stores the value under the
  simple key in `extension_data`.
- `note` renders as help text below the field.
- Field types: `checkbox`, `checkboxes`, `text`, `textarea`, `select`,
  `number`.
- `mode: "required"` shows the required marker.
- `options` supplies select/checkbox choices.

## Using field values in `handle`

For entity create/update events, `event.input` is the serialized operation
input:

```json
{
  "Record": {
    "ID": 4,
    "Name": "Introduction to Sahaj"
  },
  "Extras": {
    "com.example.optin": {
      "opt_in": "true",
      "options": "one"
    }
  }
}
```

`Extras` is namespaced by extension ID and then by the **simple field param**.
Read only your own extension's entry.

Types:

- web forms submit strings:
  - checkbox: `"true"` / `"false"`
  - multi-value field: `["one", "two"]`
- API clients may send JSON booleans, numbers, strings, or arrays.

### Go

```go
func (Extension) Handle(_ context.Context, ev *extensionapi.Event) (*extensionapi.EventResponse, error) {
    input, err := extensionapi.DecodeEntityInput[map[string]any](ev)
    if err != nil {
        return nil, err
    }

    values := input.Extras.Values("com.example.optin")
    if ev.Phase == extensionapi.Before && !extensionapi.AsBool(values["opt_in"]) {
        return &extensionapi.EventResponse{Error: "opt-in is required"}, nil
    }

    if ev.Phase == extensionapi.After {
        output, err := extensionapi.DecodeOutput[map[string]any](ev)
        if err != nil {
            return nil, err
        }
        log.Printf("saved id: %v", output["ID"])
    }

    return &extensionapi.EventResponse{}, nil
}
```

`extensionapi.AsString` and `extensionapi.AsStrings` normalize single and
multi-value fields.

### JavaScript

```js
import { asBool, extensionValues } from "@sahaj-crm/plugin-sdk";

handle: (event) => {
  const values = extensionValues(event.input?.Extras, "com.example.optin");
  if (event.phase === "before" && !asBool(values.opt_in)) {
    return { error: "opt-in is required" };
  }
  return {};
},
```

### Python

```python
from sahaj_crm_plugin_sdk import as_bool, extension_values

def handle(event):
    values = extension_values(
        event.get("input", {}).get("Extras"),
        "com.example.optin",
    )
    if event["phase"] == "before" and not as_bool(values.get("opt_in")):
        return {"error": "opt-in is required"}
    return {}
```

### Patching entity fields

A `before` handler can populate entity fields by returning `patch`:

```json
{
  "patch": {
    "MeetingLink": "https://zoom.us/j/123456789",
    "MeetingID": "123456789"
  }
}
```

Go:

```go
return &extensionapi.EventResponse{
    Patch: map[string]any{
        "MeetingLink": "https://zoom.us/j/123456789",
        "MeetingID":   "123456789",
    },
}, nil
```

Keys are Go field names, matching the `Record` payload. Values are decoded into
the field's Go type. Patches are only supported for entity create/update and
are applied before the database write. After/failed handlers receive the
patched input. Plugins are trusted, so any matching exported field can be
patched.

### Phases

- `before`: validate; return `error` to abort.
- `after`: the operation committed. `event.output` contains the saved result.
- `failed`: the operation failed. `event.error` contains the reason.

`output` is not used by the host. It exists for `after` handlers that need the
committed result for their own side effects, for example the new record ID or
the final field values. For a course create it is the course itself, not the
internally generated sessions, list, form, campaign, or page.

For archive/delete, `input` is just the entity ID and there is no `Extras`.
Delete `output` is `{"ID": <id>}`.

## Languages

The SDK is an authoring tool; the installed artifact is always a native
executable.

| Language | Package | Build |
|---|---|---|
| Go | `github.com/sahaj-crm/plugin-sdk` | `go build` |
| JavaScript | `@sahaj-crm/plugin-sdk` | `bun build --compile` |
| Python | `sahaj-crm-plugin-sdk` | PyInstaller |

## Go

```go
package main

import (
    "context"
    extensionapi "github.com/sahaj-crm/plugin-sdk"
)

type Extension struct{}

func (Extension) Contribute(context.Context, *extensionapi.ContributionRequest) (*extensionapi.ContributionResponse, error) {
    return &extensionapi.ContributionResponse{Fields: []extensionapi.Field{...}}, nil
}

func (Extension) Handle(context.Context, *extensionapi.Event) (*extensionapi.EventResponse, error) {
    return &extensionapi.EventResponse{}, nil
}

func main() {
    extensionapi.Serve(Extension{})
}
```

```sh
go build -o plugins/go-full/plugin ./examples/go-full
```

## JavaScript

```js
import { serve } from "@sahaj-crm/plugin-sdk";

serve({
  contribute: (req) => ({ fields: [...] }),
  handle: (event) => ({}),
});
```

Bundle to a native executable with Bun:

```sh
bun build --compile plugin.mjs --outfile plugins/example-optin/plugin
```

## Python

```python
from sahaj_crm_plugin_sdk import serve

serve(
    contribute=lambda req: {"fields": [...]},
    handle=lambda event: {},
)
```

Bundle to a native executable with PyInstaller:

```sh
python3 -m PyInstaller --onefile --name plugin \
  --paths python examples/python-full/plugin.py
```

## Examples

Simple examples have one checkbox. Full examples add a select with options,
notes, and output logging.

- `examples/go-simple`
- `examples/go-full`
- `examples/js-simple`
- `examples/js-full`
- `examples/python-simple`
- `examples/python-full`

Install all six into a CRM plugins directory:

```sh
./install-examples.sh /path/to/sahaj-crm/plugins
```

Requirements:

- Go
- Bun
- Python 3 with `pyinstaller`

The script builds each example into a native executable and copies the
manifest. No runtime or SDK files are installed into the plugins directory.

## Streaming and concurrency

The current protocol is request/response and the host serializes calls. One
plugin process is started at host boot and reused. stdout is reserved for
protocol messages; logs go to stderr.

## License

MIT
