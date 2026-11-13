// Minimal Go extension example: one checkbox and a Handle logger.
//
//	mkdir -p plugins/go-simple
//	go build -o plugins/go-simple/plugin ./examples/go-simple
package main

import (
	"context"
	"log"

	extensionapi "github.com/sahaj-crm/plugin-sdk"
)

const extensionID = "com.example.gosimple"

type Extension struct{}

func (Extension) Contribute(
	_ context.Context,
	_ *extensionapi.ContributionRequest,
) (*extensionapi.ContributionResponse, error) {
	return &extensionapi.ContributionResponse{
		Fields: []extensionapi.Field{{
			Name:  "Enable the simple Go feature",
			Param: "enabled",
			Type:  extensionapi.FieldCheckbox,
			Note:  "Minimal example with one checkbox.",
		}},
	}, nil
}

func (Extension) Handle(
	_ context.Context,
	ev *extensionapi.Event,
) (*extensionapi.EventResponse, error) {
	input, err := extensionapi.DecodeEntityInput[map[string]any](ev)
	if err != nil {
		return nil, err
	}
	values := input.Extras.Values(extensionID)
	log.Printf(
		"simple go extension: %s %s enabled=%v",
		ev.Phase,
		ev.Name,
		extensionapi.AsBool(values["enabled"]),
	)
	return &extensionapi.EventResponse{}, nil
}

func main() {
	extensionapi.Serve(Extension{})
}
