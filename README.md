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
go build -o plugins/example-optin/plugin ./examples/go-optin
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
  --paths python examples/python-optin/plugin.py
```

## Examples

- `examples/go-optin`
- `examples/js-optin`
- `examples/python-optin`

Install all three into a CRM plugins directory:

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
