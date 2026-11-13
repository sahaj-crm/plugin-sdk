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

## Targets

A target is `<entity>.<action>`, where the entity is the singular name of the
owning page and the action is the modal context.

| Entity | Targets |
|---|---|
| `course` | `course.new`, `course.edit`, `course.view`, `course.timeline` |
| `contact` | `contact.new`, `contact.edit`, `contact.view`, `contact.journey` |
| `list` | `list.new`, `list.edit`, `list.view`, `list.members` |
| `campaign` | `campaign.new`, `campaign.edit`, `campaign.view` |
| `form` | `form.new`, `form.edit`, `form.view` |
| `page` | `page.new`, `page.edit`, `page.view` |
| `segment` | `segment.new`, `segment.edit`, `segment.view`, `segment.members` |
| `campaign_template` | `campaign_template.new`, `.edit`, `.view` |
| `form_template` | `form_template.new`, `.edit`, `.view` |
| `page_template` | `page_template.new`, `.edit`, `.view` |

Sub-resource modals use the owning entity's action. For example, a session view
or edit modal uses `course.view` / `course.edit`.

Wildcards are supported:

```json
"targets": ["course.*", "contact.new"]
```

## Events

Standard entities emit `create`, `update`, `archive` and `delete` events:

```text
contact.*
list.*
campaign.*
form.*
page.*
segment.*
course.*
campaign_template.*
form_template.*
page_template.*
journey_entry.*
```

Special events:

| Event | When |
|---|---|
| `course.session.update` | a generated session is updated |
| `list.member.add` | contacts are added to a list |
| `list.member.remove` | contacts are removed from a list |
| `form.submit` | a public form is submitted |

Wildcards are supported:

```json
"events": ["course.*", "list.member.*"]
```

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
by `param`.

## Fields

Every field object has this shape:

| Field | Type | Required | Meaning |
|---|---|---|---|
| `name` | string | no | label shown in the form |
| `param` | string | yes | extension-local key; host namespaces it |
| `value` | any | no | initial or stored value |
| `type` | string | yes | field/input type |
| `mode` | string | no | render/validation mode |
| `placeholder` | string | no | input placeholder |
| `note` | string | no | helper text under the input |
| `options` | array | no | choices for `select` / `checkboxes` |

### Special field types

These types have dedicated rendering:

| Type | Result |
|---|---|
| `textarea` | multi-line text input |
| `select` | dropdown; requires `options` |
| `checkbox` | single boolean; submits `"true"` or `"false"` |
| `checkboxes` | multi-select checkbox group; requires `options`; stores submitted value(s) |
| `date` | date input using the host's date formatting |
| `datetime` | `datetime-local` input |
| `duration` | two HH/MM number inputs |
| `route` | route input with host availability checking |

### Native input types

Any other valid HTML input type is supported. The type is used directly as:

```html
<input type="{{ .Type }}" ... />
```

Examples:

- `text`
- `email`
- `url`
- `tel`
- `number`
- `password`
- `search`
- `hidden`
- `color`
- `range`

`file` is not supported because the host does not process multipart file
uploads for extension fields.

### Options

`options` is an array of:

```json
{"value": "one", "label": "1"}
```

It is required for `select` and `checkboxes`.

### Modes

`mode` is a single string:

| Mode | Behaviour |
|---|---|
| `""` | normal field |
| `required` | required marker and HTML `required` attribute |
| `disabled` | disabled input |
| `readonly` | readonly input |
| `hidden` | hide the entire field row |
| `controlled` | hide the entire field row |

Use one of these values. Other values are emitted as HTML attributes and are
not recommended.

### Values

- Web forms submit strings.
- `checkbox` stores `"true"` / `"false"`.
- `checkboxes` stores a string for one selection and an array for multiple.
- `duration` stores the two input values as an array.
- API clients may send JSON booleans, numbers, strings, or arrays.
- Stored extension values overwrite the field's `value` when the modal opens.

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
{
  "error": "",
  "patch": {
    "MeetingLink": "https://zoom.us/j/123",
    "MeetingID": "123"
  },
  "values": {}
}
```

- `error` aborts a `before` event when non-empty.
- `patch` applies entity field changes before a create/update is saved. Keys
  are Go field names, matching the `Record` payload. The value is decoded into
  the field's Go type.
- `values` is reserved and currently ignored by the host.

Patches are only supported on entity create/update events. They are applied
after all interested before handlers run, so a later plugin overrides an
earlier plugin's value for the same field. After/failed handlers receive the
patched input.

## `shutdown`

Params: `{}`

Result:

```json
{"ok": true}
```

The plugin should exit after replying. The host may also close stdin.
