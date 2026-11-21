# Sahaj CRM Plugin SDK

SDK and protocol for Sahaj CRM plugins.

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

The manifest is authoritative for plugin metadata.

```json
{
  "id": "com.example.optin",
  "name": "Example Opt-in",
  "description": "Adds a sample opt-in checkbox to course forms.",
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
{ "jsonrpc": "2.0", "id": 1, "method": "contribute", "params": {} }
```

Methods:

- `handshake` — protocol/readiness check handled by the SDK
- `contribute` — declarative fields for a modal
- `handle` — before/after/failed lifecycle events
- `shutdown` — optional graceful exit

The full contract is in [spec/protocol.md](spec/protocol.md).

## Fields

`contribute` returns a flat list of `Setting` objects. The same shape is used
for deployment settings:

```json
{
  "fields": [
    {
      "name": "Enable the example feature",
      "key": "opt_in",
      "type": "checkbox",
      "note": "Stored under the plugin's own data."
    }
  ]
}
```

- `key` is the plugin-local field key. The host namespaces it to
  `plugin[<plugin>][<key>]` for the form and stores the value under the
  simple key in `plugin_data`.
- `note` renders as help text below the field.
- Special types: `textarea`, `select`, `checkbox`, `checkboxes`, `date`,
  `datetime`, `duration`, `route`.
- Any other valid HTML input type works: `text`, `email`, `url`, `tel`,
  `number`, `password`, `search`, `hidden`, `color`, `range`, ...
- `mode`: `""`, `required`, `disabled`, `readonly`, `hidden`, `controlled`.
- `options` supplies `select` / `checkboxes` choices as
  `{"value": "...", "label": "..."}`.

See [spec/protocol.md](spec/protocol.md) for the full field, target and event
reference.

## Settings

Plugins can declare deployment-wide settings in `plugin.json`:

```json
{
  "settings": [
    {
      "key": "api_key",
      "name": "API Key",
      "type": "password",
      "mode": "required",
      "secret": true
    },
    {
      "key": "default_duration",
      "name": "Default Duration",
      "type": "number",
      "value": "60"
    }
  ]
}
```

Settings are stored by the host, configured by an admin, and passed to the
plugin on each `contribute` and `handle` call:

```go
apiKey := req.Settings["api_key"]
```

```go
apiKey := ev.Settings["api_key"]
```

Only the calling plugin's settings are included. `secret` settings are
never returned by the admin API after saving; submitting an empty secret keeps
the existing value.

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

`Extras` is namespaced by plugin ID and then by the **simple field key**.
Read only your own plugin's entry.

Types:

- web forms submit strings:
  - checkbox: `"true"` / `"false"`
  - multi-value field: `["one", "two"]`
- API clients may send JSON booleans, numbers, strings, or arrays.

### Go

```go
func (Plugin) Handle(_ context.Context, ev *pluginapi.Event) (*pluginapi.EventResponse, error) {
    input, err := pluginapi.DecodeEntityInput[map[string]any](ev)
    if err != nil {
        return nil, err
    }

    values := input.Extras.Values("com.example.optin")
    if ev.Phase == pluginapi.Before && !pluginapi.AsBool(values["opt_in"]) {
        return &pluginapi.EventResponse{Error: "opt-in is required"}, nil
    }

    if ev.Phase == pluginapi.After {
        output, err := pluginapi.DecodeOutput[map[string]any](ev)
        if err != nil {
            return nil, err
        }
        log.Printf("saved id: %v", output["ID"])
    }

    return &pluginapi.EventResponse{}, nil
}
```

`pluginapi.AsString` and `pluginapi.AsStrings` normalize single and
multi-value fields.

### JavaScript

```js
import { asBool, pluginValues } from "@sahaj-crm/plugin-sdk";

handle: (event) => {
  const values = pluginValues(event.input?.Extras, "com.example.optin");
  if (event.phase === "before" && !asBool(values.opt_in)) {
    return { error: "opt-in is required" };
  }
  return {};
},
```

### Python

```python
from sahaj_crm_plugin_sdk import as_bool, plugin_values

def handle(event):
    values = plugin_values(
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
return &pluginapi.EventResponse{
    Patch: map[string]any{
        "MeetingLink": "https://zoom.us/j/123456789",
        "MeetingID":   "123456789",
    },
}, nil
```

Keys are Go field names, matching the `Record` payload. Values are decoded into
the field's Go type. Patches are only supported for entity create/update and
are applied before the database write. After/failed handlers receive the
patched input.

Protected fields are omitted from `Record`. Hidden and readonly fields are
visible to plugins but cannot be patched; protected fields cannot be patched
either.

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

| Language   | Package                           | Build                 |
| ---------- | --------------------------------- | --------------------- |
| Go         | `github.com/sahaj-crm/plugin-sdk` | `go build`            |
| JavaScript | `@sahaj-crm/plugin-sdk`           | `bun build --compile` |
| Python     | `sahaj-crm-plugin-sdk`            | PyInstaller           |

## Go

```go
package main

import (
    "context"
    pluginapi "github.com/sahaj-crm/plugin-sdk"
)

type Plugin struct{}

func (Plugin) Contribute(context.Context, *pluginapi.ContributionRequest) (*pluginapi.ContributionResponse, error) {
    return &pluginapi.ContributionResponse{Fields: []pluginapi.Setting{...}}, nil
}

func (Plugin) Handle(context.Context, *pluginapi.Event) (*pluginapi.EventResponse, error) {
    return &pluginapi.EventResponse{}, nil
}

func main() {
    pluginapi.Serve(Plugin{})
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
notes, output logging, and two deployment settings (`default_mode` and a secret
`api_key`) that show how settings are read in `contribute` and `handle`.

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

### Prebuilt examples

CI builds the same six archives on every push to `main` and on every `v*` tag,
so they can be installed straight from a URL without a local toolchain.

Rolling builds from `main` (always the newest code):

```text
https://github.com/sahaj-crm/plugin-sdk/releases/download/examples-latest/<example>_<platform>.zip
```

The release description lists every URL, grouped into one block per platform,
so the block for your machine can be copied straight into the CRM plugin
installer to install all six at once.

Pinned releases use the usual latest-release URL:

```text
https://github.com/sahaj-crm/plugin-sdk/releases/latest/download/<example>_<platform>.zip
```

Examples:

```text
go-simple_darwin_arm64.zip
js-full_linux_amd64.zip
python-simple_darwin_arm64.zip
```

Platforms are `darwin_arm64` and `linux_amd64`. Each archive contains
`plugin.json` and the executable at its root, which is the layout the CRM
installer expects. PyInstaller cannot cross-compile, so every platform is built
on its own runner; add a matrix entry to cover more.

## Streaming and concurrency

The current protocol is request/response and the host serializes calls. One
plugin process is started at host boot and reused. stdout is reserved for
protocol messages; logs go to stderr.

## License

MIT
