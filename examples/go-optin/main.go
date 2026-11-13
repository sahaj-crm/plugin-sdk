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

type Input struct {
	Record map[string]any            `json:"Record"`
	Extras map[string]map[string]any `json:"Extras"`
}

func (Extension) Handle(
	_ context.Context,
	ev *extensionapi.Event,
) (*extensionapi.EventResponse, error) {
	log.Printf("example extension: %s %s", ev.Phase, ev.Name)

	var input Input
	if err := json.Unmarshal(ev.Input, &input); err != nil {
		return nil, err
	}

	// Extras is namespaced by extension id. Read only your own fields.
	values := input.Extras[name]
	optIn := values["opt_in"]
	choice := values["options"]

	log.Printf("plugin fields: opt_in=%v options=%v", optIn, choice)

	// Web forms submit strings ("true"/"false"); API clients may send booleans.
	if ev.Phase == extensionapi.Before && optIn == "false" {
		log.Printf("example extension: opt-in is disabled")
	}

	return &extensionapi.EventResponse{}, nil
}

func main() {
	extensionapi.Serve(Extension{})
}
