// Example Go plugin. Build it into a manifest directory:
//
//	mkdir -p plugins/go-full
//	go build -o plugins/go-full/plugin ./examples/go-full
package main

import (
	"context"
	"log"

	pluginAPI "github.com/sahaj-crm/plugin-sdk"
)

var name = "com.example.gofull"

type Plugin struct{}

func (Plugin) Contribute(
	_ context.Context,
	req *pluginAPI.ContributionRequest,
) (*pluginAPI.ContributionResponse, error) {
	if req.Action != "new" && req.Action != "edit" {
		return &pluginAPI.ContributionResponse{}, nil
	}
	defaultMode := pluginAPI.AsString(req.Settings["default_mode"])
	if defaultMode == "" {
		defaultMode = "one"
	}
	return &pluginAPI.ContributionResponse{
		Fields: []pluginAPI.Setting{{
			Name: "Enable the full Go feature",
			Key:  "feature_enabled",
			Type: pluginAPI.FieldCheckbox,
			Note: "Note about this checkbox",
		}, {
			Name:  "Full Go mode",
			Key:   "feature_mode",
			Type:  pluginAPI.FieldSelect,
			Value: defaultMode,
			Note:  "Defaults from the plugin setting.",
			Options: []pluginAPI.Option{
				{Value: "one", Label: "Mode 1"},
				{Value: "two", Label: "Mode 2"},
			},
		}},
	}, nil
}

func (Plugin) Handle(
	_ context.Context,
	ev *pluginAPI.Event,
) (*pluginAPI.EventResponse, error) {
	log.Printf("full go plugin: %s %s", ev.Phase, ev.Name)
	log.Printf(
		"settings: default_mode=%s api_key_configured=%v",
		pluginAPI.AsString(ev.Settings["default_mode"]),
		pluginAPI.AsString(ev.Settings["api_key"]) != "",
	)

	input, err := pluginAPI.DecodeEntityInput[map[string]any](ev)
	if err != nil {
		return nil, err
	}

	// Extras is namespaced by plugin id. Read only your own fields.
	values := input.Extras.Values(name)
	enabled := pluginAPI.AsBool(values["feature_enabled"])
	mode := pluginAPI.AsString(values["feature_mode"])

	log.Printf("plugin fields: feature_enabled=%v feature_mode=%s", enabled, mode)

	if ev.Phase == pluginAPI.Before && !enabled {
		log.Printf("full go plugin: feature is disabled")
	}

	if ev.Phase == pluginAPI.After {
		output, err := pluginAPI.DecodeOutput[map[string]any](ev)
		if err != nil {
			return nil, err
		}
		log.Printf("full go plugin: saved id=%v", output["ID"])
	}

	return &pluginAPI.EventResponse{}, nil
}

func main() {
	pluginAPI.Serve(Plugin{})
}
