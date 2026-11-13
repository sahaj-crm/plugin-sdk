// Package extensionapi is the contract shared by the CRM and extension plugins.
// It deliberately depends only on the standard library so external modules can
// import it without pulling in internal packages.
package extensionapi

import (
	"context"
	"encoding/json"
)

// Phase is the point in an operation lifecycle at which an event is emitted.
type Phase string

const (
	Before Phase = "before"
	After  Phase = "after"
	Failed Phase = "failed"
)

// Event is a serialized lifecycle event sent to interested extensions.
type Event struct {
	Name   string          `json:"name"`   // "course.create"
	Phase  Phase           `json:"phase"`  // before | after | failed
	Entity string          `json:"entity"` // "course"
	Action string          `json:"action"` // "create"
	ID     string          `json:"id,omitempty"`
	Input  json.RawMessage `json:"input,omitempty"`
	Output json.RawMessage `json:"output,omitempty"`
	Error  string          `json:"error,omitempty"`
}

// EventResponse lets a before-handler reject or modify an operation.
// Patch values are applied to the entity record before saving, keyed by Go
// field name (for example "MeetingLink").
type EventResponse struct {
	Error  string         `json:"error,omitempty"`
	Patch  map[string]any `json:"patch,omitempty"`
	Values map[string]any `json:"values,omitempty"` // reserved for extension data
}

// Info describes an extension and what it wants to receive.
type Info struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	APIVersion int      `json:"api_version"`
	Events     []string `json:"events,omitempty"`  // "course.create", ...
	Targets    []string `json:"targets,omitempty"` // "course.new", "course.edit", ...
}

// Supported form field types. These match the host's StructField.Type values.
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

// Field mirrors the host's StructField so it can be rendered by the same form
// template. Param is the field key; the host namespaces it as
// "ext[<extension>][<param>]" before rendering.
type Field struct {
	Name        string   `json:"name,omitempty"`
	Param       string   `json:"param"`
	Value       any      `json:"value,omitempty"`
	Type        string   `json:"type"`
	Mode        string   `json:"mode,omitempty"`
	Placeholder string   `json:"placeholder,omitempty"`
	Note        string   `json:"note,omitempty"`
	Options     []Option `json:"options,omitempty"`
}

// ContributionRequest asks an extension for fields for a modal target.
type ContributionRequest struct {
	Entity string          `json:"entity"`
	Action string          `json:"action"` // new | edit | view
	ID     string          `json:"id,omitempty"`
	Record json.RawMessage `json:"record,omitempty"`
	Values map[string]any  `json:"values,omitempty"` // this extension's stored values
}

// ContributionResponse is the declarative UI contribution of an extension.
type ContributionResponse struct {
	Fields []Field `json:"fields"`
}

// Extension is the behavior implemented by a plugin. Its metadata lives in
// plugin.json, so no runtime handshake is required.
type Extension interface {
	Contribute(context.Context, *ContributionRequest) (*ContributionResponse, error)
	Handle(context.Context, *Event) (*EventResponse, error)
}

// APIVersion is the current extension contract version.
const APIVersion = 1
