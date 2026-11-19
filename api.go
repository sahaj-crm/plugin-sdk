// Package pluginapi is the contract shared by the CRM and plugin plugins.
// It deliberately depends only on the standard library so external modules can
// import it without pulling in internal packages.
package pluginapi

import (
	"context"
	"encoding/json"
)

// ProtocolVersion is the wire protocol spoken by the host and plugins.
const ProtocolVersion = 1

// HandshakeRequest verifies that a started executable speaks the plugin
// protocol before the host marks it running.
type HandshakeRequest struct {
	Protocol int `json:"protocol"`
}

// HandshakeResponse is returned by the SDK's built-in handshake handler.
type HandshakeResponse struct {
	Protocol   int    `json:"protocol"`
	SDKVersion string `json:"sdk_version,omitempty"`
}

// Phase is the point in an operation lifecycle at which an event is emitted.
type Phase string

const (
	Before Phase = "before"
	After  Phase = "after"
	Failed Phase = "failed"
)

// Event is a serialized lifecycle event sent to interested plugins.
type Event struct {
	Name     string          `json:"name"`   // "course.create"
	Phase    Phase           `json:"phase"`  // before | after | failed
	Entity   string          `json:"entity"` // "course"
	Action   string          `json:"action"` // "create"
	ID       string          `json:"id,omitempty"`
	Input    json.RawMessage `json:"input,omitempty"`
	Output   json.RawMessage `json:"output,omitempty"`
	Error    string          `json:"error,omitempty"`
	Settings map[string]any  `json:"settings,omitempty"` // this plugin's settings
}

// EventResponse lets a before-handler reject or modify an operation.
// Patch values are applied to the entity record before saving, keyed by Go
// field name (for example "MeetingLink").
type EventResponse struct {
	Error  string         `json:"error,omitempty"`
	Patch  map[string]any `json:"patch,omitempty"`
	Values map[string]any `json:"values,omitempty"` // reserved for plugin data
}

// Manifest describes an plugin and what it wants to receive.
type Manifest struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Version     string    `json:"version"`
	APIVersion  int       `json:"api_version"`
	Executable  string    `json:"executable"`
	Events      []string  `json:"events,omitempty"`  // "course.create", ...
	Targets     []string  `json:"targets,omitempty"` // "course.new", "course.edit", ...
	Settings    []Setting `json:"settings,omitempty"`
}

// Setting describes one field. The same shape is used for deployment-wide
// settings and for entity-contributed fields. Secret is only meaningful for
// manifest settings.
type Setting struct {
	Name        string   `json:"name"`
	Key         string   `json:"key"`
	Value       any      `json:"value,omitempty"` // default or initial value
	Type        string   `json:"type"`
	Mode        string   `json:"mode,omitempty"`
	Placeholder string   `json:"placeholder,omitempty"`
	Note        string   `json:"note,omitempty"`
	Options     []Option `json:"options,omitempty"`
	Secret      bool     `json:"secret,omitempty"`
}

// Common form field types. Type accepts any valid HTML input type; these
// constants cover the types with dedicated host rendering and common values.
const (
	FieldCheckbox   = "checkbox"
	FieldCheckboxes = "checkboxes"
	FieldText       = "text"
	FieldTextarea   = "textarea"
	FieldSelect     = "select"
	FieldNumber     = "number"
)

// Option is one choice in a select field.
type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// ContributionRequest asks an plugin for fields for a modal target.
type ContributionRequest struct {
	Entity   string          `json:"entity"`
	Action   string          `json:"action"` // new | edit | view
	ID       string          `json:"id,omitempty"`
	Record   json.RawMessage `json:"record,omitempty"`
	Values   map[string]any  `json:"values,omitempty"`   // this plugin's stored values
	Settings map[string]any  `json:"settings,omitempty"` // this plugin's settings
}

// ContributionResponse is the declarative UI contribution of an plugin.
type ContributionResponse struct {
	Fields []Setting `json:"fields"`
}

// Plugin is the behavior implemented by a plugin. Its metadata lives in
// plugin.json, so no runtime handshake is required.
type Plugin interface {
	Contribute(context.Context, *ContributionRequest) (*ContributionResponse, error)
	Handle(context.Context, *Event) (*EventResponse, error)
}

// APIVersion is the current plugin contract version.
const APIVersion = 1
