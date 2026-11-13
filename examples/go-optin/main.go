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
			Param: "opt_in",
			Type:  extAPI.FieldCheckbox,
			Note:  "Note about this checkbox",
		}, {
			Name:  "Choose wisely",
			Param: "options",
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
	optIn := extAPI.AsBool(values["opt_in"])
	choice := extAPI.AsString(values["options"])

	log.Printf("plugin fields: opt_in=%v options=%s", optIn, choice)

	if ev.Phase == extAPI.Before && !optIn {
		log.Printf("example extension: opt-in is disabled")
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
