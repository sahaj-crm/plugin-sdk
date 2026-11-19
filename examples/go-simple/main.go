// Minimal Go plugin example: one checkbox and a Handle logger.
//
//	mkdir -p plugins/go-simple
//	go build -o plugins/go-simple/plugin ./examples/go-simple
package main

import (
	"context"
	"log"

	pluginapi "github.com/sahaj-crm/plugin-sdk"
)

const pluginID = "com.example.gosimple"

type Plugin struct{}

func (Plugin) Contribute(
	_ context.Context,
	_ *pluginapi.ContributionRequest,
) (*pluginapi.ContributionResponse, error) {
	return &pluginapi.ContributionResponse{
		Fields: []pluginapi.Setting{{
			Name: "Enable the simple Go feature",
			Key:  "enabled",
			Type: pluginapi.FieldCheckbox,
			Note: "Minimal example with one checkbox.",
		}},
	}, nil
}

func (Plugin) Handle(
	_ context.Context,
	ev *pluginapi.Event,
) (*pluginapi.EventResponse, error) {
	input, err := pluginapi.DecodeEntityInput[map[string]any](ev)
	if err != nil {
		return nil, err
	}
	values := input.Extras.Values(pluginID)
	log.Printf(
		"simple go plugin: %s %s enabled=%v",
		ev.Phase,
		ev.Name,
		pluginapi.AsBool(values["enabled"]),
	)
	return &pluginapi.EventResponse{}, nil
}

func main() {
	pluginapi.Serve(Plugin{})
}
