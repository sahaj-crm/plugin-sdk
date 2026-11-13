// Example Go extension. Build it into a manifest directory:
//
//	mkdir -p plugins/example-optin
//	go build -o plugins/example-optin/plugin ./examples/go-optin
package main

import (
	"context"
	"log"

	extAPI "github.com/sahaj-crm/plugin-sdk"
)

var name = "com.example.optin"

type Extension struct{}

func (Extension) Contribute(
	_ context.Context,
	req *extAPI.ContributionRequest,
) (*extAPI.ContributionResponse, error) {
	if req.Action != "new" && req.Action != "edit" {
		return &extAPI.ContributionResponse{}, nil
	}
	return &extAPI.ContributionResponse{
		Fields: []extAPI.Field{{
			Name:  "Enable the example feature",
			Param: "feature_enabled",
			Type:  extAPI.FieldCheckbox,
			Note:  "Note about this checkbox",
		}, {
			Name:  "Choose wisely",
			Param: "feature_mode",
			Type:  extAPI.FieldSelect,
			Note:  "Go select",
			Options: []extAPI.Option{
				{Value: "one", Label: "1"},
				{Value: "two", Label: "2"},
			},
		}},
	}, nil
}

func (Extension) Handle(
	_ context.Context,
	ev *extAPI.Event,
) (*extAPI.EventResponse, error) {
	log.Printf("example extension: %s %s", ev.Phase, ev.Name)

	input, err := extAPI.DecodeEntityInput[map[string]any](ev)
	if err != nil {
		return nil, err
	}

	// Extras is namespaced by extension id. Read only your own fields.
	values := input.Extras.Values(name)
	enabled := extAPI.AsBool(values["feature_enabled"])
	mode := extAPI.AsString(values["feature_mode"])

	log.Printf("plugin fields: feature_enabled=%v feature_mode=%s", enabled, mode)

	if ev.Phase == extAPI.Before && !enabled {
		log.Printf("example extension: feature is disabled")
	}

	if ev.Phase == extAPI.After {
		output, err := extAPI.DecodeOutput[map[string]any](ev)
		if err != nil {
			return nil, err
		}
		log.Printf("saved record: %v", output["ID"])
	}

	return &extAPI.EventResponse{}, nil
}

func main() {
	extAPI.Serve(Extension{})
}
