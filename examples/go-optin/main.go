// Example Go extension. Build it into a manifest directory:
//
//	mkdir -p plugins/example-optin
//	go build -o plugins/example-optin/plugin ./examples/go-optin
package main

import (
	"context"
	"log"

	extensionapi "github.com/sahaj-crm/plugin-sdk"
)

var name = "com.example.optin"

type Extension struct{}

func (Extension) Contribute(
	_ context.Context,
	req *extensionapi.ContributionRequest,
) (*extensionapi.ContributionResponse, error) {
	if req.Action != "new" && req.Action != "edit" {
		return &extensionapi.ContributionResponse{}, nil
	}
	return &extensionapi.ContributionResponse{
		Fields: []extensionapi.Field{{
			Name:  "Enable the example feature",
			Param: "opt_in",
			Type:  extensionapi.FieldCheckbox,
			Note:  "Note about this checkbox",
		}, {
			Name:  "Choose wisely",
			Param: "options",
			Type:  extensionapi.FieldSelect,
			Note:  "Go select",
			Options: []extensionapi.Option{
				{Value: "one", Label: "1"},
				{Value: "two", Label: "2"},
			},
		}},
	}, nil
}

func (Extension) Handle(
	_ context.Context,
	ev *extensionapi.Event,
) (*extensionapi.EventResponse, error) {
	log.Printf("example extension: %s %s", ev.Phase, ev.Name)

	input, err := extensionapi.DecodeEntityInput[map[string]any](ev)
	if err != nil {
		return nil, err
	}

	// Extras is namespaced by extension id. Read only your own fields.
	values := input.Extras.Values(name)
	optIn := extensionapi.AsBool(values["opt_in"])
	choice := extensionapi.AsString(values["options"])

	log.Printf("plugin fields: opt_in=%v options=%s", optIn, choice)

	if ev.Phase == extensionapi.Before && !optIn {
		log.Printf("example extension: opt-in is disabled")
	}

	if ev.Phase == extensionapi.After {
		output, err := extensionapi.DecodeOutput[map[string]any](ev)
		if err != nil {
			return nil, err
		}
		log.Printf("saved record: %v", output["ID"])
	}

	return &extensionapi.EventResponse{}, nil
}

func main() {
	extensionapi.Serve(Extension{})
}
