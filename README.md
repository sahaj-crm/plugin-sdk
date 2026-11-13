# Sahaj CRM Plugin SDK

SDK and protocol for Sahaj CRM extensions. Plugins are executable child
processes that speak newline-delimited JSON-RPC over stdin/stdout.

## Protocol

One JSON object per line:

```json
{"jsonrpc":"2.0","id":1,"method":"info","params":{}}
```

Methods:

- `info` — id, name, version, API version, targets, events
- `contribute` — declarative fields for a modal
- `handle` — before/after/failed lifecycle events
- `shutdown` — optional graceful exit

The full contract is in [spec/protocol.md](spec/protocol.md).

## Languages

| Language | Package |
|---|---|
| Go | `github.com/sahaj-crm/plugin-sdk` |
| JavaScript | `@sahaj-crm/plugin-sdk` |
| Python | `sahaj-crm-plugin-sdk` |

## Go

```go
package main

import (
    "context"
    extensionapi "github.com/sahaj-crm/plugin-sdk"
)

type Extension struct{}

func (Extension) Info(context.Context) (*extensionapi.Info, error) {
    return &extensionapi.Info{
        ID: "com.example.optin",
        Name: "Example Opt-in",
        Version: "0.1.0",
        APIVersion: extensionapi.APIVersion,
        Targets: []string{"course.new", "course.edit"},
        Events: []string{"course.create", "course.update"},
    }, nil
}

func (Extension) Contribute(context.Context, *extensionapi.ContributionRequest) (*extensionapi.ContributionResponse, error) {
    return &extensionapi.ContributionResponse{Sections: []extensionapi.Section{...}}, nil
}

func (Extension) Handle(context.Context, *extensionapi.Event) (*extensionapi.EventResponse, error) {
    return &extensionapi.EventResponse{}, nil
}

func main() {
    extensionapi.Serve(Extension{})
}
```

Build a plugin into a manifest directory:

```sh
mkdir -p plugins/example-optin
go build -o plugins/example-optin/plugin ./examples/go-optin
```

```json
{
  "command": "./plugin",
  "args": [],
  "cwd": "."
}
```

## JavaScript

```js
import { serve } from "@sahaj-crm/plugin-sdk";

serve({
  info: () => ({ ... }),
  contribute: (req) => ({ sections: [...] }),
  handle: (event) => ({}),
});
```

## Python

```python
from sahaj_crm_plugin_sdk import serve

serve(
    info=lambda: {...},
    contribute=lambda req: {"sections": [...]},
    handle=lambda event: {},
)
```

## Examples

- `examples/go-optin`
- `examples/js-optin`
- `examples/python-optin`

Install all three into a CRM plugins directory:

```sh
./install-examples.sh /path/to/sahaj-crm/plugins
```

The JS and Python examples vendor a copy of their SDK so they can be copied
directly into the plugins directory. The Go example is compiled by the install
script.

## Streaming and concurrency

The current protocol is request/response and the host serializes calls. One
plugin process is started at host boot and reused. stdout is reserved for
protocol messages; logs go to stderr.

## License

MIT
