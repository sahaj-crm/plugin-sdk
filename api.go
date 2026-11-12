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

// EventResponse lets a before-handler reject an operation with an error.
// Values may patch the extension's own stored data; core entity fields are
// intentionally not mutable by extensions.
type EventResponse struct {
	Error  string         `json:"error,omitempty"`
	Values map[string]any `json:"values,omitempty"`
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

// FieldType is a declarative input type rendered by the host UI.
type FieldType string

const (
	FieldCheckbox FieldType = "checkbox"
	FieldText     FieldType = "text"
	FieldTextarea FieldType = "textarea"
	FieldSelect   FieldType = "select"
	FieldNumber   FieldType = "number"
)

// Option is one choice in a select field.
type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// Field is one extension-contributed form field. Extension and Key are filled
// by the host if the extension leaves them empty.
type Field struct {
	Extension string    `json:"extension"`
	Key       string    `json:"key"`
	Label     string    `json:"label"`
	Help      string    `json:"help,omitempty"`
	Type      FieldType `json:"type"`
	Value     any       `json:"value,omitempty"`
	Options   []Option  `json:"options,omitempty"`
	Required  bool      `json:"required,omitempty"`
}

// Section groups fields contributed by one extension.
type Section struct {
	Extension string  `json:"extension"`
	Title     string  `json:"title"`
	Fields    []Field `json:"fields"`
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
	Sections []Section `json:"sections"`
}

// Extension is implemented in-process or by a stdio child-process client.
type Extension interface {
	Info(context.Context) (*Info, error)
	Contribute(context.Context, *ContributionRequest) (*ContributionResponse, error)
	Handle(context.Context, *Event) (*EventResponse, error)
}

// APIVersion is the current extension contract version.
const APIVersion = 1
