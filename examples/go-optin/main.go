// Example Go extension. Build it into a manifest directory:
//
//	mkdir -p plugins/example-optin
//	go build -o plugins/example-optin/plugin ./examples/go-optin
package main

import (
	"context"
	"encoding/json"
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
			Note:  "Stored under the extension's own data.",
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

type Input struct {
	Record map[string]any
	Extras map[string]any
}

func (Extension) Handle(
	_ context.Context,
	ev *extensionapi.Event,
) (*extensionapi.EventResponse, error) {
	log.Printf("example extension: %s %s", ev.Phase, ev.Name)
	var input Input
	err := json.Unmarshal(ev.Input, &input)
	if err != nil {
		return nil, err
	}
	log.Println(input.Extras[name])
	return &extensionapi.EventResponse{}, nil
}

func main() {
	extensionapi.Serve(Extension{})
}
