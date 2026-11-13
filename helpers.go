// Helpers for decoding handle event payloads.
package extensionapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Extras is the namespaced extension field data submitted with an operation:
// extension id -> field param -> value.
type Extras map[string]map[string]any

// EntityInput is the input shape for entity create/update events.
type EntityInput[T any] struct {
	Record T      `json:"Record"`
	Extras Extras `json:"Extras"`
}

// DecodeInput unmarshals a handle event's input into T.
func DecodeInput[T any](ev *Event) (T, error) {
	var out T
	if ev == nil || len(ev.Input) == 0 {
		return out, nil
	}
	err := json.Unmarshal(ev.Input, &out)
	return out, err
}

// DecodeEntityInput unmarshals a create/update event input.
func DecodeEntityInput[T any](ev *Event) (EntityInput[T], error) {
	return DecodeInput[EntityInput[T]](ev)
}

// DecodeOutput unmarshals a handle event's output into T. Output is only
// present on after events.
func DecodeOutput[T any](ev *Event) (T, error) {
	var out T
	if ev == nil || len(ev.Output) == 0 {
		return out, nil
	}
	err := json.Unmarshal(ev.Output, &out)
	return out, err
}

// Values returns this extension's field values, or nil when absent.
func (e Extras) Values(extension string) map[string]any {
	if e == nil {
		return nil
	}
	return e[extension]
}

// AsBool normalizes form/API boolean values.
func AsBool(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(v, "true") || v == "1"
	case float64:
		return v == 1
	case int:
		return v == 1
	default:
		return false
	}
}

// AsString normalizes a single form/API value.
func AsString(value any) string {
	if value == nil {
		return ""
	}
	if s, ok := value.(string); ok {
		return s
	}
	return fmt.Sprint(value)
}

// AsStrings normalizes single or multi-value form/API values.
func AsStrings(value any) []string {
	switch v := value.(type) {
	case nil:
		return nil
	case string:
		if v == "" {
			return nil
		}
		return []string{v}
	case []string:
		return v
	case []any:
		out := make([]string, len(v))
		for i, item := range v {
			out[i] = AsString(item)
		}
		return out
	default:
		return []string{AsString(v)}
	}
}
